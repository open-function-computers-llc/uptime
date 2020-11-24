package main

import (
	"log"

	"github.com/open-function-computers-llc/uptime/server"
	"github.com/open-function-computers-llc/uptime/site"
	"github.com/open-function-computers-llc/uptime/storage"
)

func main() {
	db, err := setUpStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer shutDownStorage(db)

	appStorage := storage.Connection{
		DB: db,
	}

	err = storage.BootstrapSites(&appStorage)
	if err != nil {
		log.Fatal(err)
	}

	existingSites := site.GetSites(&appStorage)

	shutDownChannel := make(chan string)

	for _, existingSite := range existingSites {
		site := site.Create(existingSite.URL, &appStorage)
		site.Monitor(shutDownChannel)
	}

	server := server.Create()
	server.Bootstrap(&appStorage)
	server.SetChannel(&shutDownChannel)
	server.Serve()
}
