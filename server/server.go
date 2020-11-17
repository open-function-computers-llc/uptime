package server

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/open-function-computers-llc/uptime/site"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// Server - an http server
type Server struct {
	sites  []site.Website
	router *mux.Router
	logger *logrus.Logger
	port   string
}

// Create return a new instance of a server booted up and ready to go
func Create() Server {
	s := Server{}

	// set up logging
	s.logger = logrus.New()
	s.logger.SetOutput(os.Stdout)

	// set up routing
	s.router = mux.NewRouter()
	s.setRoutes()

	// set app configuration
	s.processConfiguration()

	return s
}

// Serve will start up the http server
func (s *Server) Serve() error {
	_, err := strconv.Atoi(s.port)
	if err != nil {
		return err
	}

	s.log("Starting server on port " + s.port)

	// CORS stuff
	rules := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	})
	corsMux := rules.Handler(s.router)
	return http.ListenAndServe(":"+s.port, corsMux)
}

func (s *Server) log(messages ...interface{}) {
	s.logger.Info(messages...)
}

func (s *Server) processConfiguration() {
	s.port = "8888"
}
