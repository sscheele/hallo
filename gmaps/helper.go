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
func InitAPIKey() {
	f, err := os.Open("api-key.txt")
	check(err, "opening file")
	reader := bufio.NewReader(f)
	apiKey, err = reader.ReadString('\n')
	check(err, "reading api key from file")
	apiKey = apiKey[:len(apiKey)-1]
}

func check(err error, description string) {
	if err != nil {
		fmt.Printf("fatal error in %s: %s", description, err)
	}
}

//GetTimeToLocation is given an origin, destination, target arrival time, and, optionally, points to avoid
//it returns the estimated time, in seconds
func GetTimeToLocation(origin, destination, arrival, avoid string) int {
	if avoid == "" {
		avoid = "tolls|ferries"
	}
	var params = map[string]string{
		"key":          "",
		"origin":       origin,      //place of origin
		"destination":  destination, //destination coordinates/address
		"arrival_time": arrival,     //arrival time (seconds since epoch)
		"avoid":        avoid,       //valid values are tolls, highways, ferries, | separated
	}

	respBody := getDirs(params)
	return getTripLen(respBody)
}

//returns the trip length in seconds
func getTripLen(respBody string) int {
	//first, get the part that describes the whole trip
	i := strings.Index(respBody, `         "legs" : [`)
	respBody = respBody[i:]
	i = strings.Index(respBody, `               "steps" : [`)
	respBody = respBody[:i]

	//next, the part that describes only the duration
	i = strings.Index(respBody, `"duration" : {`)
	respBody = respBody[i:]
	i = strings.Index(respBody, `"value" :`)
	respBody = respBody[i:]
	i = strings.Index(respBody, "\n")
	respBody = respBody[:i]
	respBody = strings.Split(respBody, " : ")[1]

	fmt.Fscanf(strings.NewReader(respBody), "%d", &i)
	return i
}

func getDirs(params map[string]string) string {
	params["key"] = apiKey //shave off trailing newline

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://maps.googleapis.com/maps/api/directions/json", nil)
	check(err, "creating http request")
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	resp, err := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)

	return string(respBody)
}