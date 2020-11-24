package server

import (
	"io"
	"net/http"
	"strconv"

	"github.com/open-function-computers-llc/uptime/site"
)

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sites := site.GetSites(s.storage)

		output := `<h1>Sites</h1><ul>`
		for _, site := range sites {
			output += "<li>" + site.URL + ": "
			if site.IsUp {
				output += "Online"
			} else {
				output += "DOWN!!!"
			}
			output += "<br /> <a href='/details/" + strconv.Itoa(site.ID) + "'>Details</a> <a href='/remove/" + strconv.Itoa(site.ID) + "'>Delete</a></li>"
		}
		output += "</ul><a href='/add'>Add Site</a>"
		io.WriteString(w, output)
	}
}
