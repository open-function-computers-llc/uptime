package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/open-function-computers-llc/uptime/server"
	"github.com/open-function-computers-llc/uptime/site"
	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	// shared logger instance
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	// first things first, read env
	err := godotenv.Load()
	if err != nil {
		logger.Error(err)
	}

	db, err := setUpStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	appStorage := storage.Connection{
		DB:     db,
		Logger: logger,
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
