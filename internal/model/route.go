package model

// Route defines the structure for an API route
type Route struct {
	Path        string
	Service     string
	ServicePort int
}
