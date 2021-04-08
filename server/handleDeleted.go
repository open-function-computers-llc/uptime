package server

import (
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/open-function-computers-llc/uptime/site"
)

func (s *Server) handleDeleted() http.HandlerFunc {
	buttonHTML := "<a class='button' href='/'>Home</a>"

	return func(w http.ResponseWriter, r *http.Request) {
		sites := site.GetDeletedSites(s.storage)
		keys := make([]int, 0, len(sites))
		for k := range sites {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		output := "<h1>Sites</h1>" + buttonHTML + "<ul>"
		for _, index := range keys {
			site := sites[index]
			output += "<li><p>" + site.URL + "</p>"
			output += "<span> <a class='button' href='/restore/" + strconv.Itoa(site.ID) + "'>Restore</a> <a class='button' href='/purge/" + strconv.Itoa(site.ID) + "'>Purge</a></span></li>"
		}
		if len(keys) == 0 {
			output += "<li><p>No deleted sites.</p></li>"
		}
		output += "</ul>" + buttonHTML + extraHTML
		io.WriteString(w, output)
	}
}
