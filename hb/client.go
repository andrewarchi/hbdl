package hb

import (
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

// Client manages Humble Bundle API requests.
type Client struct {
	c    http.Client
	csrf string
}

// NewClient constructs a client.
func NewClient() (*Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	return &Client{c: http.Client{Jar: jar}}, nil
}
