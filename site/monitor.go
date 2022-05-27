package site

import (
	"time"
)

// Monitor - periodically make an HTTP GET request to the site's URL, and optionally log it in the database
func (s *Website) Monitor(shutdownChan *chan string) {
	go func() {
		secondsDown := 0

		for {
			select {
			case msg := <-*shutdownChan:
				if msg == s.URL {
					s.Logger.Info("Shutting down monitor for " + s.URL)
					return
				}

				s.Logger.Info("Site: " + s.URL + "passing url back to channel " + msg)
				*shutdownChan <- msg
			default:
				// nothing to do as the default, this is just here so that the
				// channel checking is non-blocking
			}

			statusCodeOrTimeoutValue := s.getStatusCode()
			// s.Logger.Info(s.URL+":", statusCode)
			if statusCodeOrTimeoutValue <= 99 {
				// fake status code returned because of slow response
				secondsDown += statusCodeOrTimeoutValue
				s.setSiteDown(s.DB, secondsDown)

				secondsDown += 5
				time.Sleep(time.Second * 5) // wait 5 seconds and try again
				continue
			}
			if statusCodeOrTimeoutValue == 200 {
				s.setSiteUp(s.DB, secondsDown)

				secondsDown = 0

				time.Sleep(time.Second * 30) // wait 30 seconds and try again
				continue
			}

			s.setSiteDown(s.DB, secondsDown)

			// wait 5 seconds and try again
			time.Sleep(time.Second * 5)
			secondsDown += 5
		}
	}()
}
