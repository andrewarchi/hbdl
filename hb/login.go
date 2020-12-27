package hb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

// Login-related errors:
var (
	ErrNoCSRFCookie  = errors.New("csrf_cookie not found")
	ErrGuardRequired = errors.New("guard required: enter the code sent to your email address to verify your account")
	ErrGuardInvalid  = errors.New("guard invalid: the provided code is invalid")
)

// Login signs into an account with username (email), password, and
// guard. The guard 2FA code is sent to your email address when guard is
// empty.
func (c *Client) Login(username, password, guard string) error {
	csrf, err := c.getCSRF()
	if err != nil {
		return err
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
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("CSRF-Prevention-Token", csrf)

	resp, err := c.c.Do(req)
	if err != nil {
		return err
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
				return ErrGuardInvalid
			case data.GuardRequired:
				return ErrGuardRequired
			case data.Errors != nil:
				return LoginError(data.Errors)
			}
		}
		return fmt.Errorf("login: status %s", resp.Status)
	}
	return nil
}

var wwwURL = &url.URL{Scheme: "https", Host: "www.humblebundle.com"}

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

	for _, cookie := range c.c.Jar.Cookies(wwwURL) {
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

// LoadCookies reads cookies from a file and adds them to the client.
func (c *Client) LoadCookies(cookieFile string) error {
	f, err := os.Open(cookieFile)
	if err != nil {
		return err
	}
	var cookies []*http.Cookie
	if err := json.NewDecoder(f).Decode(&cookies); err != nil {
		return err
	}
	c.c.Jar.SetCookies(wwwURL, cookies)
	return nil
}

// SaveCookies writes the client's cookies to a file.
func (c *Client) SaveCookies(cookieFile string) error {
	f, err := os.Create(cookieFile)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(c.c.Jar.Cookies(wwwURL))
}
