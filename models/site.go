package models

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Site - a site that we will be checking
type Site struct {
	IsUp                 bool
	IsDeleted            bool
	IsMonitoring         bool
	StandardWarningSent  bool
	EmergencyWarningSent bool
	ID                   int
	URL                  string
	Meta                 string
	Uptime_1day          string
	Uptime_7day          string
	Uptime_30day         string
	Uptime_60day         string
	Uptime_90day         string
}

func (s *Site) GetStatusCodeAndError() (int, error) {
	if strings.TrimSpace(s.URL) == "" {
		return 404, errors.New("site url is empty")
	}

	output, err := exec.Command("/usr/bin/curl", "--max-time", os.Getenv("INTERVAL_HOW_LONG_FOR_SITE_TIMEOUT"), "--user-agent", "OFC_Uptime_Bot-version:CURL", "-I", "-L", s.URL).Output()

	if err != nil {
		return 500, errors.New("curl timed out")
	}

	outputLines := strings.Split(string(output), "\n")
	responseLine := outputLines[0]
	responseItems := strings.Fields(responseLine)
	if len(responseItems) < 2 {
		return 500, errors.New("cURL output: " + responseLine)
	}
	resInt, err := strconv.Atoi(responseItems[1])
	if err != nil {
		return 500, err
	}

	if resInt == 301 || resInt == 302 {
		lastHTTPLine := ""
		for _, line := range outputLines {
			if !strings.Contains(line, "HTTP") {
				continue
			}
			lastHTTPLine = line
		}

		responseItems := strings.Fields(lastHTTPLine)
		if len(responseItems) < 2 {
			return 500, errors.New("cURL output: " + responseLine)
		}
		resInt, err := strconv.Atoi(responseItems[1])
		if err != nil {
			return 500, err
		}
		return resInt, nil
	}

	return resInt, nil
}
