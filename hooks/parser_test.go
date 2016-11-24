package hooks

import (
	"fmt"
	"testing"

	"github.com/sscheele/hallo/weather"
)

func TestWeatherParse(t *testing.T) {
	ParseWeatherHook("temp<32|precip>.5:moo")
	if len(WeatherHooks) != 1 {
		fmt.Println("Wrong number of weather hooks")
	}
	testDP1 := weather.DataPoint{
		Temperature:       40,
		PrecipProbability: .6,
	}
	testDP2 := weather.DataPoint{
		Temperature:       30,
		PrecipProbability: 0,
	}
	testDP3 := weather.DataPoint{
		Temperature:       50,
		PrecipProbability: .1,
	}
	fmt.Println("Temp: 40, Precip: .6 (should be true):", WeatherHooks[0].Verify(testDP1))
	fmt.Println("Temp: 30, Precip: 0 (should be true):", WeatherHooks[0].Verify(testDP2))
	fmt.Println("Temp: 50, Precip: .1 (should be false):", WeatherHooks[0].Verify(testDP3))
}
