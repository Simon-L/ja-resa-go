package entities

import (
	"encoding/xml"
)

// a CalDAV calendar query object
type ContactQuery struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:carddav addressbook-query"`
	Prop    *Prop    `xml:",omitempty"`
	//Filter  *Filter           `xml:",omitempty"`
}

// creates a new CalDAV query for iCalendar events from a particular time range
func NewContactQueryWithProps(props ...string) *ContactQuery {
	cProps := make([]*CProp, len(props))
	for i, prop := range props {
		cProps[i] = &CProp{Name: prop}
	}

	return &ContactQuery{
		Prop: &Prop{
			GetETag: new(GetETag),
			AddressData: &AddressData{
				Prop: cProps,
			},
		},
	}
}

func NewDefaultContactQuery() *ContactQuery {
	return &ContactQuery{
		Prop: &Prop{
			GetETag:     new(GetETag),
			AddressData: new(AddressData),
		},
	}
}
