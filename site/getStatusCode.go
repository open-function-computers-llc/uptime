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

	output, err := exec.Command("/usr/bin/curl", "--max-time", "10", "--user-agent", "OFC_Uptime_Bot-version:CURL", "-I", s.URL).Output()

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

	return resInt, ""

	// req, err := http.NewRequest("GET", s.URL, nil)

	// if err != nil {
	// 	if strings.HasSuffix(err.Error(), "(Client.Timeout exceeded while awaiting headers)") {
	// 		// s.Logger.Error(s.URL + " took too long to respond, and the URL was different!")
	// 		return s.clientTimeout, err.Error()
	// 	}
	// 	return s.clientTimeout, err.Error()
	// }
	// req.Header.Set("User-Agent", "OFC_Uptime_Bot-version:"+os.Getenv("VERSION"))

	// client := *(s.httpClient)
	// resp, err := client.Do(req)
	// if err != nil {
	// 	if strings.HasSuffix(err.Error(), "(Client.Timeout exceeded while awaiting headers)") {
	// 		// s.Logger.Error(s.URL + " took too long to respond!")
	// 		req.Header.Set("Connection", "Close")
	// 		client.Do(req)
	// 		return s.clientTimeout, err.Error()
	// 	}
	// 	req.Header.Set("Connection", "Close")
	// 	client.Do(req)
	// 	return 500, err.Error()
	// }
	// defer resp.Body.Close()

	// return 200, ""
}
