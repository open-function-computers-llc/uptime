package site

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Website - a site that we will be checking
type Website struct {
	URL                string
	IsUp               bool
	UpAt               time.Time
	LastOutageDuration time.Duration
}

// Create - Make a new instance of a Website struct
func Create(address string) Website {
	w := Website{
		URL: address,
	}
	log.Println("Created Website:", address)
	return w
}

func (s *Website) Monitor() {
	go func() {
		for {
			statusCode := s.getStatusCode()
			log.Println(s.URL+":", statusCode)
			if statusCode == 200 {
				s.IsUp = true
				time.Sleep(time.Second * 15)
				continue
			}

			s.IsUp = false
			time.Sleep(time.Second * 1)
			fmt.Println("looping")
		}
	}()
}

func (s *Website) getStatusCode() int {
	if s.URL == "" {
		return 404
	}

	resp, err := http.Get(s.URL)
	if err != nil {
		log.Println(err.Error())
		return 500
	}

	return resp.StatusCode
}
