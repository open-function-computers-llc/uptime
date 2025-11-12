package server

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/open-function-computers-llc/uptime/models"
	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/petaki/inertia-go"
	"github.com/sirupsen/logrus"
)

// Server - an http server
type server struct {
	sites          map[int]*models.Site
	router         *http.ServeMux
	storage        *storage.Connection
	curlTimeout    int
	port           string
	distFS         fs.FS
	inertiaManager *inertia.Inertia
}

// Bootstrap - share a pointer to a SQL DB storage struct with this server
func Create(fileSystem fs.FS, url string) (server, error) {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	connection, err := storage.NewConnection(logger)
	if err != nil {
		return server{}, err
	}

	s := server{
		sites:          map[int]*models.Site{},
		storage:        connection,
		port:           os.Getenv("APP_PORT"),
		distFS:         fileSystem,
		inertiaManager: inertia.NewWithFS(url, "dist/index.html", "", fileSystem),
	}

	// share the app environment with the frontend
	s.inertiaManager.Share("appEnvironment", os.Getenv("APP_ENV"))

	s.inertiaManager.ShareFunc("assetPath", s.assetPath)

	// set up routing
	s.router = &http.ServeMux{}
	s.setRoutes()

	s.log("Loading sites...")
	err = s.loadSites()
	if err != nil {
		return s, err
	}

	return s, nil
}

// Serve will start up the http server
func (s *server) Serve() error {
	s.log("Starting server on port " + s.port)

	return http.ListenAndServe(":"+s.port, s.router)
}

func (s *server) log(messages ...interface{}) {
	s.storage.Logger.Info(messages...)
}

func (s *server) loadSites() error {
	sites, err := storage.GetSites(s.storage)

	// loop over the fresh sites in the DB
	for _, site := range sites {
		//  add to the server list if it's missing
		if _, ok := s.sites[site.ID]; !ok {
			s.sites[site.ID] = site
			continue
		}

		// update the existing site with the latest goodies from the DB
		s.sites[site.ID].IsDeleted = site.IsDeleted
		s.sites[site.ID].Uptime_1day = site.Uptime_1day
		s.sites[site.ID].Uptime_7day = site.Uptime_7day
		s.sites[site.ID].Uptime_30day = site.Uptime_30day
		s.sites[site.ID].Uptime_60day = site.Uptime_60day
		s.sites[site.ID].Uptime_90day = site.Uptime_90day
		s.sites[site.ID].URL = site.URL
		s.sites[site.ID].Meta = site.Meta
	}

	// make sure all loaded sites that are not deleted are monitoring
	for _, site := range s.sites {
		if site.IsDeleted {
			storage.StopMonitoringSite(site, s.storage.Logger)
			continue
		}
		if !site.IsMonitoring {
			go storage.Monitor(site, s.storage)
		}
	}
	return err
}
