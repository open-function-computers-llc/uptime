package server

import "net/http"

// here are all the routes and their handlers
func (s *Server) setRoutes() {
	// basic routes
	routes := map[string]http.HandlerFunc{
		"/":             s.handleIndex(),
		"/add":          s.handleSiteForm(),
		"/store":        s.handleStoreSite(),
		"/remove/{id}":  s.handleRemoveSite(),
		"/restore/{id}": s.handleRestoreSite(),
		"/purge/{id}":   s.handlePurgeSite(),
		"/details/{id}": s.handleSiteDetailsByID(),
		"/details":      s.handleSiteDetails(),
		"/deleted":      s.handleDeleted(),
	}

	// wrap the routes in basic middleware stack
	for route, handler := range routes {
		s.router.HandleFunc(route, s.logRequest(handler))
	}
}
