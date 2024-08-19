package site

import (
	"os/exec"
	"strconv"
	"strings"
)

func (s *Website) getStatusCodeAndErrorMessage() (int, string) {
	if s.URL == "" {
		return 404, "404: URL is blank"
	}

	output, err := exec.Command("/usr/bin/curl", "--max-time", "10", "--user-agent", "OFC_Uptime_Bot-version:CURL", "-I", "-L", s.URL).Output()

	if err != nil {
		return 500, err.Error()
	}

	outputLines := strings.Split(string(output), "\n")
	responseLine := outputLines[0]
	responseItems := strings.Fields(responseLine)
	if len(responseItems) < 2 {
		return 500, "Response line from cURL: " + responseLine
	}
	resInt, err := strconv.Atoi(responseItems[1])
	if err != nil {
		return 500, "Response line from cURL: " + responseLine + " | " + err.Error()
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
			return 500, "Response line from cURL: " + responseLine
		}
		resInt, err := strconv.Atoi(responseItems[1])
		if err != nil {
			return 500, "Response line from cURL: " + responseLine + " | " + err.Error()
		}
		return resInt, ""
	}

	return resInt, ""
}
