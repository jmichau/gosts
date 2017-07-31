// Package gosts provide middleware for support Strict-Transport-Security header.
// Read more here: https://tools.ietf.org/html/rfc6797
package gosts

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

var (
	// headerPieceValue handles cached last part of header value
	headerPieceValue string

	// maxAge handles Info.MaxAge
	maxAge int

	// expires handles Info.Expires
	expires time.Time
)

// Info contains the HSTS configurations
type Info struct {
	// MaxAge sets the duration (number of seconds) that the HSTS is valid for.
	// This value can't be less than zero. This property is required
	// cause of Expires that use MaxAge as a fallback if duration between now and Expires is less than zero.
	MaxAge int `yaml:"MaxAge"`

	// Expires sets the date after which the header will not be valid.
	// If this property is set the MaxAge property is not respected,
	// but if duration between Expires and now is less than zero then MaxAge will be respected automatically.
	Expires time.Time `yaml:"Expires"`

	// IncludeSubDomains specifying that this HSTS Policy also applies to any hosts whose
	// domain names are subdomains of the Known HSTS Host's domain name.
	IncludeSubDomains bool `yaml:"IncludeSubDomains"`

	// SendPreloadDirective sets whether the preload directive should be set. The directive allows browsers to
	// confirm that the site should be added to a preload list. (see https://hstspreload.appspot.com/)
	SendPreloadDirective bool `yaml:"SendPreloadDirective"`
}

// Configure HSTS middleware
func Configure(i *Info) error {
	// Validate
	if i.MaxAge < 0 {
		return errors.New("HSTS MaxAge duration can't be less than zero.")
	}

	if i.IncludeSubDomains {
		headerPieceValue += "; includeSubDomains"
	}
	if i.SendPreloadDirective {
		headerPieceValue += "; preload"
	}

	expires = i.Expires
	maxAge = i.MaxAge

	return nil
}

// Header is http middleware that adds RFC 6797 Strict-Transport-Security header.
// Note that this header is ignored by browsers on not-secure response
// e.g. when connection is over HTTP protocol or SSL certificate is self-signed.
func Header(h http.Handler) http.Handler {

	if expires.Sub(time.Now()) > 0*time.Second {
		// Expires strategy

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			if expires.Sub(now) > 0 {
				w.Header().Set("Strict-Transport-Security", "max-age="+strconv.Itoa(int(expires.Sub(now).Seconds()))+headerPieceValue)
			} else if maxAge >= 0 {
				w.Header().Set("Strict-Transport-Security", "max-age="+strconv.Itoa(maxAge)+headerPieceValue)
			}

			h.ServeHTTP(w, r)
		})
	}
	// Else MaxAge strategy
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age="+strconv.Itoa(maxAge)+headerPieceValue)

		h.ServeHTTP(w, r)
	})
}
