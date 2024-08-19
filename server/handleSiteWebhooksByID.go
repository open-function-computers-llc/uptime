package server

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/open-function-computers-llc/uptime/site"
)

// Depreciated. Check out the new handler s.handleSiteDetails()
func (s *Server) handleSiteWebhooksByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		siteID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
		}
		site, err := site.FindWebsiteByID(siteID, s.storage, s.logger)
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
		}

		output := "<h1>Webhooks for " + site.URL + "</h1>"

		standardHooks, err := site.GetWebHooks("standard")
		if err != nil {
			w.Write([]byte("<h1>Error getting webhooks for site " + site.URL + "</h1><br />" + err.Error()))
		}
		emergencyHooks, err := site.GetWebHooks("emergency")
		if err != nil {
			w.Write([]byte("<h1>Error getting webhooks for site " + site.URL + "</h1><br />" + err.Error()))
		}
		allHooks := append(standardHooks, emergencyHooks...)

		for _, hook := range allHooks {
			output += "<p><strong>(" + hook.Verb + ") " + hook.Name + ":</strong> " + hook.URL + "<a href='/remove-webhook/" + strconv.Itoa(hook.ID) + "' class='remove-button'>&times;</a></p>"
		}

		addForm = strings.ReplaceAll(addForm, "%%SITEID%%", r.PathValue("id"))

		io.WriteString(w, htmlWrap(output+addForm))
	}
}
