package site

import (
	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/sirupsen/logrus"
)

// Website - a site that we will be checking
type Website struct {
	IsUp                 bool
	standardWarningSent  bool
	emergencyWarningSent bool
	DB                   *storage.Connection
	Logger               *logrus.Logger
	ID                   int
	clientTimeout        int
	URL                  string
}
