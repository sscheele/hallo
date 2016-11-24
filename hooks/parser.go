package hooks

import (
	"strconv"
	"strings"

	"github.com/sscheele/hallo/cal"
	"github.com/sscheele/hallo/weather"
)

//all hooks should be formatted like this (only OR is supported):
//condition_1|condition_2:action

//WeatherHook changes the behavior of the alarm based on the weather
type WeatherHook struct {
	//Verify may analyze any field of a DataPoint, and returns a bool based on that analysis
	Verify func([]func(weather.DataPoint) bool) bool
	//Action is performed if Verify returns true
	Action string
}

//CalendarHook changes the behavior of the alarm based on the calendar
type CalendarHook struct {
	//Verify may analyze a calendar event's time and summary, and returns a bool based on that analysis
	Verify func([]func(cal.Event) bool) bool
	//Action is performed if Verify returns true
	Action string
}

//ScriptHook changes the behavior of the alarm based on the output of a script
type ScriptHook struct {
	//Verify accepts the command to run and returns a bool based on its output
	Verify func([]func(string) bool) bool
	//Action is performed if Verify returns true
	Action string
}

var (
	//WeatherHooks is the list of every weather hook
	WeatherHooks []WeatherHook
	//CalendarHooks is the list of every calendar hook
	CalendarHooks []CalendarHook
	//ScriptHooks is the list of every script hook
	ScriptHooks []ScriptHook
)

//ParseWeatherHook parses a string to add a hook to WeatherHooks
func ParseWeatherHook(s string) {
	//Acceptable values for weather hooks: Temperature[<=>]%d, rain, snow, sleet, wind, fog, PrecipProbability[<=>]%d
	var hooks []func(weather.DataPoint) bool
	cam := strings.Split(s, ":")
	if len(cam) != 2 {
		return
	}
	conditions := strings.Split(cam[0], "|")
	for _, cond := range conditions {
		cond = strings.ToLower(cond)
		switch {
		case cond == "rain", cond == "snow", cond == "sleet", cond == "wind", cond == "fog":
			hooks = append(hooks, func(dp weather.DataPoint) bool {
				return dp.Icon == cond
			})
		case len(cond) > 5 && cond[:4] == "temp":
			tFields, comp := splitComp(cond)
			if len(tFields) != 2 {
				continue
			}
			temp, err := strconv.ParseFloat(tFields[1], 64)
			if err != nil {
				continue
			}
			switch comp {
			case "<":
				hooks = append(hooks, func(dp weather.DataPoint) bool {
					return dp.Temperature < t
				})
			case "=":
				hooks = append(hooks, func(dp weather.DataPoint) bool {
					return dp.Temperature == t
				})
			case ">":
				hooks = append(hooks, func(dp weather.DataPoint) bool {
					return dp.Temperature > t
				})
			}
		case len(cond) > 7 && cond[:6] == "precip":
			tFields, comp := splitComp(cond)
			if len(tFields) != 2 {
				continue
			}
			temp, err := strconv.ParseFloat(tFields[1], 64)
			if err != nil {
				continue
			}
			switch comp {
			case "<":
				hooks = append(hooks, func(dp weather.DataPoint) bool {
					return dp.PrecipProbability < t
				})
			case "=":
				hooks = append(hooks, func(dp weather.DataPoint) bool {
					return dp.PrecipProbability == t
				})
			case ">":
				hooks = append(hooks, func(dp weather.DataPoint) bool {
					return dp.PrecipProbability > t
				})
			}
		}
	}
}

//splitComp splits s at the first instance of a comparison (<, =, >) and returns the fields and the comparison character it found
func splitComp(s string) ([]string, string) {
	var chr string
	for i := 0; i < len(s); i++ {
		chr = s[i : i+1]
		if chr == "<" || chr == "=" || chr == ">" {
			return []string{s[:i], s[i+1:]}, chr
		}
	}
	return []string{s}, ""
}
