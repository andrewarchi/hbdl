package hb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// Errors from API
var (
	ErrNoCSRFCookie  = errors.New("cookie not found: csrf_cookie")
	ErrGuardRequired = errors.New("guard required: enter the code sent to your email address to verify your account")
	ErrGuardInvalid  = errors.New("guard invalid: check the code sent to your email address")
)

// Login signs into an account with username (email), password, and
// guard. The guard 2FA code is sent to your email address.
func (c *Client) Login(username, password, guard string) (*http.Response, error) {
	csrf, err := c.getCSRF()
	if err != nil {
		return nil, err
	}

	form := url.Values{
		"username":                 {username},
		"password":                 {password},
		"guard":                    {guard},
		"access_token":             {""},
		"access_token_provider_id": {""},
		"goto":                     {"/"},
		"qs":                       {""},
	}
	fr := strings.NewReader(form.Encode())
	req, err := http.NewRequest("POST", "https://www.humblebundle.com/processlogin", fr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("CSRF-Prevention-Token", csrf)

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var data struct {
			GuardRequired bool                `json:"humble_guard_required"`
			Errors        map[string][]string `json:"errors"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
			switch {
			case data.GuardRequired && guard != "":
				return nil, ErrGuardInvalid
			case data.GuardRequired:
				return nil, ErrGuardRequired
			case data.Errors != nil:
				return nil, LoginError(data.Errors)
			}
		}
		return nil, fmt.Errorf("login: status %s", resp.Status)
	}
	return resp, nil
}

var hbURL = &url.URL{Scheme: "https", Host: "humblebundle.com"}

func (c *Client) getCSRF() (string, error) {
	if c.csrf != "" {
		return c.csrf, nil
	}

	// Visit some page so the CSRF cookie is set
	resp, err := c.c.Get("https://www.humblebundle.com/login")
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("csrf: status %s", resp.Status)
	}
	resp.Body.Close()

	for _, cookie := range c.c.Jar.Cookies(hbURL) {
		if cookie.Name == "csrf_cookie" {
			c.csrf = cookie.Value
			return cookie.Value, nil
		}
	}
	return "", ErrNoCSRFCookie
}

// LoginError is returned on username or password failure.
type LoginError map[string][]string

func (err LoginError) Error() string {
	var b strings.Builder
	categories := make([]string, 0, len(err))
	for category := range err {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	for _, category := range categories {
		errs := err[category]
		if b.Len() != 0 {
			b.WriteByte('\n')
		}
		switch len(errs) {
		case 0:
		case 1:
			fmt.Fprintf(&b, "%s: %s", category, errs[0])
		default:
			fmt.Fprintf(&b, "%s:", category)
			for _, e := range errs {
				fmt.Fprintf(&b, "\n  %s", e)
			}
		}
	}
	return b.String()
}
