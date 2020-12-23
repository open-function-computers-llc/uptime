package server

import (
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/open-function-computers-llc/uptime/site"
)

func (s *Server) handleIndex() http.HandlerFunc {
	buttonHTML := "<a class='button' href='/add'>Add Site</a>"

	return func(w http.ResponseWriter, r *http.Request) {
		sites := site.GetSites(s.storage)
		keys := make([]int, 0, len(sites))
		for k := range sites {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		output := "<h1>Sites</h1>" + buttonHTML + "<ul>"
		for _, index := range keys {
			site := sites[index]
			status := "Online"
			class := "online"

			if !site.IsUp {
				status = "DOWN!!!"
				class = "down"
			}
			output += "<li class='" + class + "'><p>" + site.URL + ": " + status + "</p>"
			output += "<span> <a class='button' href='/details?url=" + site.URL + "'>Details</a> <a class='button' href='/remove/" + strconv.Itoa(site.ID) + "'>Delete</a></span></li>"
		}
		output += "</ul>" + buttonHTML + extraHTML
		io.WriteString(w, output)
	}
}
