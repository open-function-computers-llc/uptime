package main

import (
	"log"
	"os"

	"github.com/open-function-computers-llc/uptime/server"
	"github.com/open-function-computers-llc/uptime/site"
	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	// shared logger instance
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	db, err := setUpStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer shutDownStorage(db)

	appStorage := storage.Connection{
		DB:     db,
		Logger: logger,
	}

	err = storage.BootstrapSites(&appStorage)
	if err != nil {
		log.Fatal(err)
	}

	existingSites := site.GetSites(&appStorage)

	shutDownChannel := make(chan string)

	for _, existingSite := range existingSites {
		site := site.Create(existingSite.URL, &appStorage, logger)
		site.Monitor(shutDownChannel)
	}

	server := server.Create()
	server.Bootstrap(&appStorage, logger)
	server.SetChannel(&shutDownChannel)
	server.Serve()
}
