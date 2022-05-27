package site

import (
	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/sirupsen/logrus"
)

// Website - a site that we will be checking
type Website struct {
	ID                   int
	URL                  string
	IsUp                 bool
	DB                   *storage.Connection
	Logger               *logrus.Logger
	standardWarningSent  bool
	emergencyWarningSent bool
}
