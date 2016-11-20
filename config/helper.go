package config

import (
	"github.com/BurntSushi/toml"
)

//Cfg contains the current configuration
var Cfg Config

//Config can be used to store a configuration
type Config struct {
	//Location contains the location for the user
	Location []string
	//CalUpdatePeriod contains the number of minutes to wait before updating the calendar
	CalUpdatePeriod int
	//WeatherUpdatePeriod contains the number of minutes to wait before updating the weather
	WeatherUpdatePeriod int
	//WeatherLookAhead is the number of hours to look ahead for the weather
	WeatherLookAhead int
	//MapUpdatePeriod contains the number of minutes to wait before updating a maps location
	MapUpdatePeriod int
	//TwelveHour tells whether to use a 12-hour or a 24-hour clock
	TwelveHour bool
	//NumCalEvents gives the number of calendar events to read
	NumCalEvents int
	//AudioFilePath contains the path of the audio file to use for the alarm
	AudioFilePath string
	//TimeBeforeLeave is the number of seconds to give you before you have to leave
	TimeBeforeLeave int
}

func init() {
	_, err := toml.DecodeFile("config.toml", &Cfg)
	if err != nil {
		useDefaults()
	}
}

func useDefaults() {
	Cfg.Location = []string{"38.8522392", "-77.3368576"}
	Cfg.CalUpdatePeriod = 15
	Cfg.WeatherUpdatePeriod = 45
	Cfg.WeatherLookAhead = 6
	Cfg.MapUpdatePeriod = 10
	Cfg.TwelveHour = false
	Cfg.NumCalEvents = 10
	Cfg.AudioFilePath = "audio/bell.aiff"
	Cfg.TimeBeforeLeave = 1800 //thirty minutes
}
