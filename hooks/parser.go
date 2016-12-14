package hooks

import (
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sscheele/hallo/cal"
	"github.com/sscheele/hallo/weather"
)

//all hooks should be formatted like this except script hooks (only OR is supported):
//condition_1|condition_2:action
//script hooks may only have one condition, and their only action is to disable the alarm

//WeatherHook changes the behavior of the alarm based on the weather
type WeatherHook struct {
	//Verify may analyze any field of a DataPoint, and returns a bool based on that analysis
	Verify func(weather.DataPoint) bool
	//Action is performed if Verify returns true
	Action string
}

//CalendarHook creates or deletes alarms based on the calendar
type CalendarHook struct {
	//Verify may analyze a calendar event's time and summary, and returns a bool based on that analysis
	Verify func(cal.Event) bool
	//Action is performed if Verify returns true
	Action string
}

//ScriptHook may "disable" an alarm based on the output of a script
type ScriptHook struct {
	//Verify returns a bool based on the output of a command
	Verify func() bool
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

//ParseScriptHook parses a string to add a hook to ScriptHooks
func ParseScriptHook(s string) {
	ScriptHooks = append(ScriptHooks, ScriptHook{
		Verify: func() bool {
			cmd := exec.Command(s)
			sout, err := cmd.StdoutPipe()
			if err != nil {
				return false
			}
			_, err := cmd.StderrPipe()
			if err != nil {
				return false
			}
			if err := cmd.Start(); err != nil {
				return false
			}
			retVal, err := ioutil.ReadAll(sout)
			if len(retVal) > 0 && (retVal[0] == "t" || retVal[0] == "T") {
				return true
			}
			return false
		},
	})
}

//ParseCalendarHook parses a string to add a hook to CalendarHooks
func ParseCalendarHook(s string) {
	//Calendar hooks allow you to do something if a calendar event contains a certain word
	var retVal CalendarHook
	var hooks []func(cal.Event) bool
	ca := strings.Split(s)
	if len(ca) != 2 {
		return
	}
	conditions := strings.Split(ca[0], "|")
	for _, cond := range conditions {
		hooks = append(hooks, func(c cal.Event) bool {
			return strings.Contains(c.Summary, cond)
		})
	}
	retVal.Verify = func(c cal.Event) bool {
		for _, hook := range hooks {
			if hook(c) {
				return true
			}
		}
		return false
	}
	retVal.Action = ca[1]
	CalendarHooks = append(CalendarHooks, retVal)
}

//ParseWeatherHook parses a string to add a hook to WeatherHooks
func ParseWeatherHook(s string) {
	//Acceptable values for weather hooks: Temperature[<=>]%d, rain, snow, sleet, wind, fog, PrecipProbability[<=>]%d
	var retVal WeatherHook
	var hooks []func(weather.DataPoint) bool
	ca := strings.Split(s, ":")
	if len(ca) != 2 {
		return
	}
	conditions := strings.Split(ca[0], "|")
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
					return dp.Temperature < temp
				})
			case "=":
				hooks = append(hooks, func(dp weather.DataPoint) bool {
					return dp.Temperature == temp
				})
			case ">":
				hooks = append(hooks, func(dp weather.DataPoint) bool {
					return dp.Temperature > temp
				})
			}
		case len(cond) > 7 && cond[:6] == "precip":
			tFields, comp := splitComp(cond)
			if len(tFields) != 2 {
				continue
			}
			precip, err := strconv.ParseFloat(tFields[1], 64)
			if err != nil {
				continue
			}
			switch comp {
			case "<":
				hooks = append(hooks, func(dp weather.DataPoint) bool {
					return dp.PrecipProbability < precip
				})
			case "=":
				hooks = append(hooks, func(dp weather.DataPoint) bool {
					return dp.PrecipProbability == precip
				})
			case ">":
				hooks = append(hooks, func(dp weather.DataPoint) bool {
					return dp.PrecipProbability > precip
				})
			}
		}
	}
	retVal.Verify = func(dp weather.DataPoint) bool {
		for _, hook := range hooks {
			if hook(dp) {
				return true
			}
		}
		return false
	}
	retVal.Action = ca[1]
	WeatherHooks = append(WeatherHooks, retVal)
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
