package server

import "net/http"

// here are all the routes and their handlers
func (s *server) setRoutes() {
	// basic routes
	routes := map[string]http.HandlerFunc{
		"POST /remove/{id}":        s.handleSiteRemove(),
		"POST /restore/{id}":       s.handleSiteRestore(),
		"GET /webhooks/{id}":       s.handleSiteWebhooks(),
		"POST /store-webhook/{id}": s.handleSiteStoreWebhook(),
		"DELETE /delete-webhook":   s.handleSiteDeleteWebhook(),
		"GET /edit/{id}":           s.handleSiteEdit(),
		"POST /update/{id}":        s.handleSiteUpdate(),
		"GET /purge/{id}":          s.handleSitePurgeCheck(),
		"POST /purge/{id}":         s.handleSitePurge(),
		"GET /api/load-sites":      s.handleApiSites(),
		"GET /add":                 s.handleSiteAdd(),
		"POST /store":              s.handleSiteStore(),
		"GET /details/{id}":        s.handleSiteDetails(),
		"GET /details":             s.handleSiteDetails(),
		"GET /":                    s.handleIndex(),
	}

	// wrap the routes in basic middleware stack
	for route, handler := range routes {
		s.router.Handle(route, s.inertiaManager.Middleware(s.logRequest(handler)))
	}

	s.router.Handle("GET /dist/", http.FileServer(http.FS(s.distFS)))
}
