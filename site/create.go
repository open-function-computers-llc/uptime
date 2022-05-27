package site

import (
	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/sirupsen/logrus"
)

// Create - Make a new instance of a Website struct
func Create(address string, dbConn *storage.Connection, logger *logrus.Logger) Website {
	w := Website{
		URL:    address,
		IsUp:   true,
		DB:     dbConn,
		Logger: logger,
	}
	logger.Info("Created Website:", address)

	siteDatabaseID := w.GetSiteID(dbConn)
	if siteDatabaseID == 0 {
		err := storage.AddSite(w.URL, dbConn)
		if err != nil {
			logger.Info("Couldn't add new site to DB:", err)
		}
	}
	return w
}
