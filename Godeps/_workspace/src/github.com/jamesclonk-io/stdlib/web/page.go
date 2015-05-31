package web

import "net/http"

// PageMaster contains various default values.
type PageMaster struct {
	Title      string
	Content    interface{}
	Data       interface{}
	Template   string
	StatusCode int
	Navigation Navigation
}

// Page represents a page to be rendered on the browser.
type Page struct {
	Title      string
	Navigation Navigation
	ActiveLink string
	Headers    http.Header
	Content    interface{}
	Data       interface{}
	Template   string
	StatusCode int
	Error      error
}

// Navigation represents the navigation for the web page.
type Navigation []NavigationElement

// NavigationElement is a navigation element of the web page or nested inside another NavigationElement.
type NavigationElement struct {
	Name     string              `json:"name"`
	Link     string              `json:"link,omitempty"`
	Icon     string              `json:"icon,omitempty"`
	Dropdown []NavigationElement `json:"dropdown,omitempty"`
}
