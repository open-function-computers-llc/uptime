package notifications

import (
	"net/http"
	"net/url"
	"os"
)

func SendHTTPRequest(domain string) {
	sendGotifyRequest := os.Getenv("GOTIFY_NOTIFICATION")
	if sendGotifyRequest != "true" {
		return
	}

	data := url.Values{
		"message":  {"OFCO.911!!! " + domain + " is Down!"},
		"priority": {"8"},
	}

	http.PostForm(os.Getenv("GOTIFY_HOST")+"message?token="+os.Getenv("GOTIFY_TOKEN"), data)
}
