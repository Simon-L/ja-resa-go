package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jkrecek/caldav-go/caldav"
	calentities "github.com/jkrecek/caldav-go/caldav/entities"
	"github.com/jkrecek/caldav-go/icalendar/components"
)

var client *caldav.Client

// EventJSON ...
type EventJSON struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Telephone string    `json:"tel"`
	Password  string    `json:"password"`
}

func uriFromPath(path string) (URI string, err error) {
	switch path {
	case "/a/music":
		return "music_test", nil
	case "/a/live-perf":
		return "live-perf_test", nil
	case "/a/ja-events":
		return "ja-events_test", nil
	case "/a/redbox":
		return "red-box_test", nil
	default:
		return "", errors.New("")
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	path, err := uriFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	fmt.Println(string(body))
	var e EventJSON
	err = json.Unmarshal(body, &e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	fmt.Println(e)
	fmt.Println(path)

	path = fmt.Sprintf("/%s/%s.ics", path, e.ID)
	err = client.DeleteEvent(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "{}")
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	path, err := uriFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	fmt.Println(string(body))
	var e EventJSON
	err = json.Unmarshal(body, &e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	fmt.Println(e)

	uuid := fmt.Sprintf("jaresa-%d", e.Start.Unix())
	e.ID = uuid
	putEvent := components.NewEventWithEnd(uuid, e.Start, e.End)
	putEvent.Summary = e.Title
	putEvent.Description = e.Telephone

	// generate an ICS filepath
	path = fmt.Sprintf("/%s/%s.ics", path, uuid)

	fmt.Println(putEvent)
	// save the event to the server, then fetch it back out
	if err = client.PutEvents(path, putEvent); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	body, err = json.Marshal(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	fmt.Println(body)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(body)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	path, err := uriFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	s, _ := strconv.Atoi(r.URL.Query().Get("start"))
	e, _ := strconv.Atoi(r.URL.Query().Get("end"))
	var eventsj []EventJSON

	start := time.Unix(int64(s), 0).Truncate(time.Hour).UTC()
	end := time.Unix(int64(e), 0).Truncate(time.Hour).UTC()
	query, err := calentities.NewEventRangeQuery(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	events, err := client.QueryEvents(path, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	for _, ev := range events {
		evj := EventJSON{
			ev.UID,
			ev.Summary,
			ev.DateStart.NativeTime(),
			ev.DateEnd.NativeTime(),
			ev.Description,
			ev.Description,
		}
		// fmt.Fprintf(w, "%s", evj)
		fmt.Println(ev)
		eventsj = append(eventsj, evj)
	}
	content, err := json.Marshal(eventsj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if len(eventsj) == 0 {
		fmt.Fprintf(w, "[]")
	} else {
		fmt.Fprintf(w, "%s", content)
	}
}

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	hc := &http.Client{Transport: tr}

	// create a reference to your CalDAV-compliant server
	var server, err = caldav.NewServer(ServerURL)
	fmt.Println(err)

	// create a CalDAV client to speak to the server
	client = caldav.NewClient(server, hc)

	// start executing requests!
	err = client.ValidateServer("/")
	fmt.Println(err)

	http.HandleFunc("/a/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlePost(w, r)
		case http.MethodGet:
			handleGet(w, r)
		case http.MethodDelete:
			handleDelete(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	})

	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
