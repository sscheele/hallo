package weather

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	//ErrNoData will be returned when there is no hourly data for the time specified
	ErrNoData = errors.New("No hourly data for specified time")
	//ErrBadResp will be returned when unable to correctly parse the response
	ErrBadResp = errors.New("Badly formatted JSON response")
	//ErrBadTime will be returned when you give the program a bad number of hours
	ErrBadTime = errors.New("Bad number of hours")
	apiKey     string
)

//DataPoint contains weather information for a time
//From github.com/mlbright/darksky
type DataPoint struct {
	Time                   float64 `json:"time"`
	Summary                string  `json:"summary"`
	Icon                   string  `json:"icon"`
	SunriseTime            float64 `json:"sunriseTime"`
	SunsetTime             float64 `json:"sunsetTime"`
	PrecipIntensity        float64 `json:"precipIntensity"`
	PrecipIntensityMax     float64 `json:"precipIntensityMax"`
	PrecipIntensityMaxTime float64 `json:"precipIntensityMaxTime"`
	PrecipProbability      float64 `json:"precipProbability"`
	PrecipType             string  `json:"precipType"`
	PrecipAccumulation     float64 `json:"precipAccumulation"`
	Temperature            float64 `json:"temperature"`
	TemperatureMin         float64 `json:"temperatureMin"`
	TemperatureMinTime     float64 `json:"temperatureMinTime"`
	TemperatureMax         float64 `json:"temperatureMax"`
	TemperatureMaxTime     float64 `json:"temperatureMaxTime"`
	ApparentTemperature    float64 `json:"apparentTemperature"`
	DewPoint               float64 `json:"dewPoint"`
	WindSpeed              float64 `json:"windSpeed"`
	WindBearing            float64 `json:"windBearing"`
	CloudCover             float64 `json:"cloudCover"`
	Humidity               float64 `json:"humidity"`
	Pressure               float64 `json:"pressure"`
	Visibility             float64 `json:"visibility"`
	Ozone                  float64 `json:"ozone"`
	MoonPhase              float64 `json:"moonPhase"`
}

func init() {
	f, err := os.Open("api-key.txt")
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

//GetNHours returns an array of DataPoints containing weather information for now and for n hours from now
func GetNHours(n int, lat string, long string) ([2]DataPoint, error) {
	var retVal [2]DataPoint
	if n < 1 {
		return retVal, ErrBadTime
	}
	var hourlyWeather []DataPoint
	var currWeather DataPoint
	res, err := getResponse(lat, long)
	if err != nil {
		return retVal, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return retVal, err
	}
	currText, hourly := parseWeatherString(string(body))
	err = json.Unmarshal([]byte(hourly), &hourlyWeather)
	if err != nil {
		return retVal, err
	}
	err = json.Unmarshal([]byte(currText), &currWeather)
	if err != nil {
		return retVal, err
	}
	retVal[0] = currWeather
	if len(hourlyWeather) < n {
		return retVal, ErrNoData
	}
	retVal[1] = hourlyWeather[n]
	return retVal, nil
}

func getResponse(lat string, long string) (*http.Response, error) {
	return http.Get(fmt.Sprintf("https://api.darksky.net/forecast/%s/%s,%s", apiKey, lat, long))
}

//returns only the JSON of the current and hourly data
func parseWeatherString(s string) (string, string) {
	ind := strings.Index(s, `"currently"`)
	if ind == -1 || ind+12 >= len(s) {
		return "", ""
	}
	s = s[ind+12:]
	ind = strings.Index(s, "}")
	if ind == -1 || ind+1 > len(s) {
		return "", ""
	}
	currentStr := s[:ind+1]
	ind = strings.Index(s, `"hourly"`)
	if ind == -1 || ind+9 >= len(s) {
		return "", ""
	}
	s = s[ind+9:]
	ind = strings.Index(s, `"data"`)
	if ind == -1 || ind+7 >= len(s) {
		return "", ""
	}
	s = s[ind+7:]
	ind = strings.Index(s, "]")
	if ind == -1 || ind+1 > len(s) {
		return "", ""
	}
	return currentStr, s[:ind+1]
}
