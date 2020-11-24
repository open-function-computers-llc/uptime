package server

import (
	"net/http"
	"os"
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
func (s *Server) Bootstrap(dbConn *storage.Connection) {
	s.storage = dbConn

	// set up logging
	s.logger = logrus.New()
	s.logger.SetOutput(os.Stdout)

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

// SetChannel - set the broadcast channel to destroy site monitors
func (s *Server) SetChannel(c *chan string) {
	s.shutdownChannel = c
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
