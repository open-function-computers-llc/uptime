package site

import (
	"net/http"

	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/sirupsen/logrus"
)

// Create - Make a new instance of a Website
func Create(address string, dbConn *storage.Connection, logger *logrus.Logger, client *http.Client, timeout int) Website {
	w := Website{
		URL:           address,
		IsUp:          true,
		DB:            dbConn,
		Logger:        logger,
		httpClient:    client,
		clientTimeout: timeout,
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
