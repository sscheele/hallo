package gmaps

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var apiKey string

//InitAPIKey reads in the API key from the default file
func InitAPIKey() error {
	f, err := os.Open("api-key.txt")
	if err != nil {
		return err
	}
	reader := bufio.NewReader(f)
	apiKey, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	apiKey = apiKey[:len(apiKey)-1]
	return nil
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
	var (
		rv1 int
		rv2 int
	)
	traffInd := strings.LastIndex(respBody, `"travelDurationTraffic"`)
	noTraffInd := strings.LastIndex(respBody, `"travelDuration"`)
	if traffInd == -1 || noTraffInd == -1 {
		return -1, -1
	}
	noTraffInd += 17 //len(`"travelDuration"`)
	traffInd += 24   //len(`"travelDurationTraffic"`)
	var i int
	for i = noTraffInd; isNumberChar(respBody[i]); i++ {
		//do nothing; we only care about i
	}
	noTraffStr := respBody[noTraffInd:i]
	fmt.Fscanf(strings.NewReader(noTraffStr), "%d", &rv2)
	for i = traffInd; isNumberChar(respBody[i]); i++ {
		//do nothing; we only care about i
	}
	traffStr := respBody[traffInd:i]
	fmt.Fscanf(strings.NewReader(traffStr), "%d", &rv1)
	return rv1, rv2 //TODO: ACTUALLY DO THIS
}

func isNumberChar(r byte) bool {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
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
	fmt.Println(req.URL.RawQuery)
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	return string(respBody)
}