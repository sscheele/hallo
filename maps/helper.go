package maps

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/sscheele/hallo/config"
)

var apiKey string

func init() {
	f, err := os.Open(config.Cfg.MapsAPIPath)
	if err != nil {
		return
	}
	reader := bufio.NewReader(f)
	apiKey, err = reader.ReadString('\n')
	if err != nil {
		return
	}
	apiKey = apiKey[:len(apiKey)-1]
	return
}

//GetTimeToLocation is given an origin, destination, target arrival time, and, optionally, points to avoid
//it returns the estimated time, in seconds, both with and without traffic
func GetTimeToLocation(params map[string]string) (int, int) {
	respBody := getDirs(params)
	return getTripLen(respBody)
}

//returns the trip length in seconds, both with and without traffic
func getTripLen(respBody string) (int, int) {
	//LOOK FOR "travelDurationTraffic"
	traffInd := strings.LastIndex(respBody, `"travelDurationTraffic"`)
	noTraffInd := strings.LastIndex(respBody, `"travelDuration"`)
	if traffInd == -1 || noTraffInd == -1 {
		return -1, -1
	}
	noTraffInd += 17 //len(`"travelDuration"`)
	traffInd += 24   //len(`"travelDurationTraffic"`)
	var i int
	for i = noTraffInd; !isSeparator(respBody[i]); i++ {
		//do nothing; we only care about i
	}
	noTraffStr := respBody[noTraffInd:i]
	rv2, err := strconv.Atoi(noTraffStr)
	if err != nil {
		return -1, -1
	}
	for i = traffInd; !isSeparator(respBody[i]); i++ {
		//do nothing; we only care about i
	}
	traffStr := respBody[traffInd:i]
	rv1, err := strconv.Atoi(traffStr)
	if err != nil {
		return -1, -1 //error value
	}

	return rv1, rv2
}

func isSeparator(r byte) bool {
	switch r {
	case ',', '}', ']', ':', '(', ')', '[', '{':
		return true
	}
	return false
}

func getDirs(params map[string]string) string {
	params["key"] = apiKey //shave off trailing newline
	params["optimize"] = "timeWithTraffic"

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://dev.virtualearth.net/REST/v1/Routes", nil)
	if err != nil {
		return ""
	}
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	//fmt.Println(req.URL.RawQuery)
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	return string(respBody)
}
