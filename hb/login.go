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
	"time"
)

// Login-related errors:
var (
	ErrNoCSRFCookie  = errors.New("csrf_cookie not found")
	ErrGuardRequired = errors.New("guard required: enter the code sent to your email address to verify your account")
	ErrGuardInvalid  = errors.New("guard invalid: the provided code is invalid")
	Err2FARequired   = errors.New("code required: enter the code from your two-factor authenticator to verify your account")
	Err2FAInvalid    = errors.New("code invalid: the provided two-factor code is invalid")
	ErrLoginFailed   = errors.New("login failed")
)

// Login signs into an account with username (email) and password. A
// guard code or 2FA code is required.
func (c *Client) Login(username, password string) error {
	return c.login(username, password, "", "")
}

// LoginGuard signs into an account with username (email), password, and
// a guard code sent to your email address. When the guard code is
// empty, an email will be sent.
func (c *Client) LoginGuard(username, password, guard string) error {
	return c.login(username, password, guard, "")
}

// Login2FA signs into an account with username (email), password, and
// a 2FA code from an authenticator app.
func (c *Client) Login2FA(username, password, code string) error {
	return c.login(username, password, "", code)
}

func (c *Client) login(username, password, guard, code string) error {
	csrf, err := c.getCSRF()
	if err != nil {
		return err
	}

	form := url.Values{
		"username":                 {username},
		"password":                 {password},
		"guard":                    {guard},
		"code":                     {code},
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
			GuardRequired     bool                `json:"humble_guard_required"`
			TwoFactorRequired bool                `json:"two_factor_required"`
			TwoFactorType     string              `json:"twofactor_type"` // "google"
			Errors            map[string][]string `json:"errors"`
			Success           bool                `json:"success"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
			switch {
			case data.GuardRequired:
				if guard == "" {
					return ErrGuardRequired
				}
				return ErrGuardInvalid
			case data.TwoFactorRequired:
				if code == "" {
					return Err2FARequired
				}
				return Err2FAInvalid
			case data.Errors != nil:
				return LoginError(data.Errors)
			case !data.Success:
				return ErrLoginFailed
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
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("csrf: status %s", resp.Status)
	}

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
	var jsonCookies []*jsonCookie
	if err := json.NewDecoder(f).Decode(&jsonCookies); err != nil {
		return err
	}
	cookies := make([]*http.Cookie, len(jsonCookies))
	for i, c := range jsonCookies {
		cookies[i] = (*http.Cookie)(c)
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
	cookies := c.c.Jar.Cookies(wwwURL)
	jsonCookies := make([]*jsonCookie, len(cookies))
	for i, c := range cookies {
		jsonCookies[i] = (*jsonCookie)(c)
	}
	return json.NewEncoder(f).Encode(jsonCookies)
}

// jsonCookie adds json tags to http.Cookie
type jsonCookie struct {
	Name       string        `json:"name,omitempty"`
	Value      string        `json:"value,omitempty"`
	Path       string        `json:"path,omitempty"`
	Domain     string        `json:"domain,omitempty"`
	Expires    time.Time     `json:"expires,omitempty"`
	RawExpires string        `json:"raw_expires,omitempty"`
	MaxAge     int           `json:"max_age,omitempty"`
	Secure     bool          `json:"secure,omitempty"`
	HttpOnly   bool          `json:"http_only,omitempty"`
	SameSite   http.SameSite `json:"same_site,omitempty"`
	Raw        string        `json:"raw,omitempty"`
	Unparsed   []string      `json:"unparsed,omitempty"`
}
