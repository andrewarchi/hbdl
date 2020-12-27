package hb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// Library-related errors:
var (
	ErrNoHomeJSONData = errors.New("home json data not found")
)

// GetGamekeys fetches the gamekeys (IDs) for each purchase.
func (c *Client) GetGamekeys() ([]string, error) {
	resp, err := c.c.Get("https://www.humblebundle.com/home/purchases")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("gamekeys: status %s", resp.Status)
	}
	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	data := findNode(node, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == "user-home-json-data" {
					return true
				}
			}
		}
		return false
	})
	if data == nil {
		return nil, ErrNoHomeJSONData
	}
	var d struct {
		Gamekeys []string `json:"gamekeys"`
		// Other fields present
	}
	r := strings.NewReader(data.FirstChild.Data)
	if err := json.NewDecoder(r).Decode(&d); err != nil {
		return nil, err
	}
	return d.Gamekeys, nil
}

func findNode(node *html.Node, matches func(*html.Node) bool) *html.Node {
	if node == nil {
		return nil
	}
	if matches(node) {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if match := findNode(child, matches); match != nil {
			return match
		}
	}
	return nil
}
