package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/open-function-computers-llc/uptime/site"
	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// Server - an http server
type Server struct {
	sites           []site.Website
	router          *mux.Router
	logger          *logrus.Logger
	port            string
	storage         *storage.Connection
	shutdownChannel *chan string
}

// Bootstrap - share a pointer to a SQL DB storage struct with this server
func (s *Server) Bootstrap(dbConn *storage.Connection, logger *logrus.Logger, shutdownChan *chan string) {
	s.storage = dbConn

	// set up logging
	s.logger = logger

	// set up the communication channel
	s.shutdownChannel = shutdownChan

	// set app configuration
	s.processConfiguration()

	// set up routing
	s.router = mux.NewRouter()
	s.setRoutes()
}

// Create return a new instance of a server booted up and ready to go
func Create() Server {
	s := Server{}

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
