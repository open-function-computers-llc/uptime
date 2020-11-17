package main

import (
	"github.com/open-function-computers-llc/uptime/server"
	"github.com/open-function-computers-llc/uptime/site"
)

var sites []site.Website

func main() {
	addresses := []string{
		"https://openfunctioncomputers.com",
		"http://kce.ofco.test",
		"http://localhost:8000",
	}

	for _, address := range addresses {
		site := site.Create(address)
		site.Monitor()
		sites = append(sites, site)
	}

	server := server.Create()
	server.Serve()
}
