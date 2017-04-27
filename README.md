# go-hsts
Package for provide middleware for support Strict-Transport-Security header.
## Install:
`go get github.com/mj420/go-hsts`

## Usage:
```golang
package main

import (
	"io"
	"net/http"
	"time"

	"github.com/mj420/go-hsts"
	"github.com/pressly/chi"
)

// -----------------------------------------------------------------------------
// Main
// -----------------------------------------------------------------------------

func main() {
	// config for hsts middleware
	hstsConf := &go_hsts.Info{
		MaxAge:               60 * 60 * 24,
		Expires:              time.Now().Add(24 * time.Hour),
		IncludeSubDomains:    true,
		SendPreloadDirective: false,
	}

	r := chi.NewRouter()

	// middleware
	go_hsts.Configure(hstsConf)
	r.Use(go_hsts.Header)

	r.Get("/", helloHandlerGET)
	// start listener
	http.ListenAndServeTLS(":443", "server.crt", "server.key", r)
}

// -----------------------------------------------------------------------------
// Handlers
// -----------------------------------------------------------------------------

func helloHandlerGET(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World!")
}

```
## Options

Option | Type | Description
--- | --- | ---
`MaxAge` | int | MaxAge sets the duration (number of seconds) that the HSTS is valid for. This value can't be less than zero. This property is required cause of Expires that use MaxAge as a fallback if duration between now and Expires is less than zero.
`Expires` | time.Time | Expires sets the date after which the header will not be valid. If this property is set the MaxAge property is not respected, but if duration between Expires and now is less than zero then MaxAge will be respected automatically.
`IncludeSubDomains` | bool | IncludeSubDomains specifying that this HSTS Policy also applies to any hosts whose domain names are subdomains of the Known HSTS Host's domain name.
`SendPreloadDirective` | bool | SendPreloadDirective sets whether the preload directive should be set. The directive allows browsers to confirm that the site should be added to a preload list. (see https://hstspreload.appspot.com/)



## TODO:
- [ ] add json config support
- [ ] add example that use yaml and/or json config
- [ ] add tests
