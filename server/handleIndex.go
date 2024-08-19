package server

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/open-function-computers-llc/uptime/site"
)

func printPercent(p float64) string {
	return strconv.FormatFloat(p*100, 'f', 2, 64) + "%"
}

func (s *Server) handleIndex() http.HandlerFunc {
	buttonHTML := "<span class='button-wrapper'><a class='button' href='/add'>Add Site</a>"
	buttonHTML += "<a class='button' href='/deleted'>Show Deleted Sites</a></span>"

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
			class := "online"

			if !site.IsUp {
				class = "down"
			}
			output += "<li class='" + class + "'>"
			output += "<p>" + site.URL + "<br /><small>"
			output += "<strong>1:</strong>" + printPercent(site.CalcUptime(1, s.storage))
			output += "<br /><strong>7:</strong>" + printPercent(site.CalcUptime(7, s.storage))
			output += "<br /><strong>30:</strong>" + printPercent(site.CalcUptime(30, s.storage))
			output += "<br /><strong>60:</strong>" + printPercent(site.CalcUptime(60, s.storage))
			output += "<br /><strong>90:</strong>" + printPercent(site.CalcUptime(90, s.storage))
			output += "</small></p>"
			output += "<span> <a class='button' href='/details?url=" + site.URL + "'>Details</a> <a class='button' href='/webhooks/" + strconv.Itoa(site.ID) + "'>Webhooks</a> <a class='button' href='/remove/" + strconv.Itoa(site.ID) + "'>Delete</a></span></li>"
		}
		output += "</ul>" + buttonHTML

		// used for page "refreshing"
		isAJAX := r.URL.Query().Get("ajax")
		if isAJAX == "true" {
			w.Write([]byte(output))
			return
		}

		w.Write([]byte(htmlWrap(output)))
	}
}
