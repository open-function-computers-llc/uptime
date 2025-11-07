package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/open-function-computers-llc/uptime/certificate"
	"github.com/open-function-computers-llc/uptime/domain"
	"github.com/open-function-computers-llc/uptime/models"
	"github.com/open-function-computers-llc/uptime/storage"
)

func (s *server) handleSiteDetails() http.HandlerFunc {
	type detailsOutput struct {
		CertInfo   certificate.CertificateInfo `json:"certInfo"`
		DomainInfo domain.WhoisInfo            `json:"domainInfo"`
		Outages    []*models.Outage            `json:"outages"`
		Up         bool                        `json:"up"`
		Uptime     map[string]float64          `json:"uptime"`
		Url        string                      `json:"url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		var site *models.Site

		// find site by URL
		siteURL := r.FormValue("url")
		if siteURL != "" {
			for _, loadedSite := range s.sites {
				if siteURL == loadedSite.URL {
					site = loadedSite
					break
				}
			}
		}

		// find site by ID
		if site == nil {
			siteID, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				http.Redirect(w, r, "/?error=invalid details `id` param", http.StatusFound)
			}
			site = s.sites[siteID]
		}

		// at this point, we better have a site
		if site == nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
			return
		}

		outages, err := storage.SiteOutages(site.ID, s.storage)
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
			return
		}

		updays1, _ := strconv.Atoi(strings.Replace(site.Uptime_1day, ".", "", 1))
		updays7, _ := strconv.Atoi(strings.Replace(site.Uptime_7day, ".", "", 1))
		updays30, _ := strconv.Atoi(strings.Replace(site.Uptime_30day, ".", "", 1))
		updays60, _ := strconv.Atoi(strings.Replace(site.Uptime_60day, ".", "", 1))
		updays90, _ := strconv.Atoi(strings.Replace(site.Uptime_90day, ".", "", 1))

		domainInfo, err := domain.GetDomainInfo(site.URL)
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
			return
		}
		certInfo, err := certificate.GetCertificateInfo(site.URL)
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
			return
		}

		output := detailsOutput{
			Up:         site.IsUp,
			Url:        site.URL,
			Outages:    outages,
			CertInfo:   certInfo,
			DomainInfo: domainInfo,
			Uptime: map[string]float64{
				"days1":  float64(updays1) / 10000,
				"days7":  float64(updays7) / 10000,
				"days30": float64(updays30) / 10000,
				"days60": float64(updays60) / 10000,
				"days90": float64(updays90) / 10000,
			},
		}
		dataJSON, _ := json.Marshal(output)
		w.Header().Set("Content-Type", "application/json")
		w.Write(dataJSON)
	}
}
