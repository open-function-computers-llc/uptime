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
	storage         *storage.Connection
	shutdownChannel *chan string
	httpClient      *http.Client
	clientTimeout   int
	port            string
}

// Bootstrap - share a pointer to a SQL DB storage struct with this server
func Create(dbConn *storage.Connection, logger *logrus.Logger, shutdownChan *chan string, client *http.Client, timeout int) Server {
	s := Server{
		storage:         dbConn,
		logger:          logger,
		shutdownChannel: shutdownChan,
		httpClient:      client,
		clientTimeout:   timeout,
	}

	// set app configuration
	s.processConfiguration()

	// set up routing
	s.router = mux.NewRouter()
	s.setRoutes()

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
