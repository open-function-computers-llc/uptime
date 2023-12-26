package main

import (
	"log"
	"net/http"
	"os"
	"time"

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

	err = site.CloseOutPreviousOutages(&appStorage)
	if err != nil {
		log.Fatal(err)
	}
	existingSites := site.GetSites(&appStorage)

	// http client for all requests
	httpClientTimeout := 10
	client := &http.Client{
		Timeout: time.Duration(httpClientTimeout) * time.Second,
	}
	shutDownChannel := make(chan string)

	for _, existingSite := range existingSites {
		site := site.Create(existingSite.URL, &appStorage, logger, client, httpClientTimeout)
		site.Monitor(&shutDownChannel)
	}

	server := server.Create(&appStorage, logger, &shutDownChannel, client, httpClientTimeout)
	server.Serve()
}
