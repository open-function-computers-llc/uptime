package models

import "time"

type Outage struct {
	ID        int       `json:"id"`
	WebsiteID int       `json:"siteID"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Duration  int       `json:"duration"`
}
