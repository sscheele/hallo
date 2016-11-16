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
//it returns the estimated time, in seconds
func GetTimeToLocation(params map[string]string) int {
	if params["avoid"] == "" {
		params["avoid"] = "tolls|ferries"
	}

	respBody := getDirs(params)
	return getTripLen(respBody)
}

//returns the trip length in seconds
func getTripLen(respBody string) int {
	//first, get the part that describes the whole trip
	i := strings.Index(respBody, `"steps" : [`)
	respBody = respBody[:i]
	i = strings.Index(respBody, `"legs" : [`)
	respBody = respBody[i:]

	var tmp string
	var dur int

	fmt.Fscanf(strings.NewReader(respBody), `"duration" : {%s}`, &tmp)
	fmt.Fscanf(strings.NewReader(tmp), `"value" : %d`, &dur)

	return dur
}

func getDirs(params map[string]string) string {
	params["key"] = apiKey //shave off trailing newline

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://maps.googleapis.com/maps/api/directions/json", nil)
	if err != nil {
		return ""
	}
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	return string(respBody)
}
