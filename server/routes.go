package server

import "net/http"

// here are all the routes and their handlers
func (s *Server) setRoutes() {
	// basic routes
	routes := map[string]http.HandlerFunc{
		"/": s.handleIndex(),
		// "/intent": s.handleIntent(),
	}

	// wrap the routes in basic middleware stack
	for route, handler := range routes {
		s.router.HandleFunc(route, s.logRequest(handler))
	}
}
