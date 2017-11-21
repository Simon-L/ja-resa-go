# ja-resa-go

* Add server URL in `config.go`
* Get static content: `git clone https://github.com/crucialHawg/ja-resa-web static`
* Run with `go run ja-resa.go config.go`

Content is served at `0.0.0.0:8080/`

`caldav-go` is vendored because of a *bug (??)* with HTTP requests not including the Depth header.
