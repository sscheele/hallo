package main

import (
	"bufio"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/sscheele/hallo/alclock"
	"github.com/sscheele/hallo/audio"
	"github.com/sscheele/hallo/cal"
	"github.com/sscheele/hallo/config"
	"github.com/sscheele/hallo/parseCmd"
	"github.com/sscheele/hallo/weather"
)

const numDataRows = 4
const banner = ` |_| _ | | _
 | |(_|| |(_)`

var (
	alarmList      []alclock.Alarm
	nextAlarm      alclock.Alarm
	nextAlarmIndex int
	dataFieldMut   sync.Mutex
	inReader       *bufio.Reader
	inputChan      chan string
	mainGUI        *gocui.Gui
	useGUI         = true
)

func updateNextAlarm() {
	if !(len(alarmList) > 0) {
		return
	}
	tempMin := alarmList[0]
	for i, al := range alarmList {
		if al.NextGoesOff.Before(tempMin.NextGoesOff) {
			tempMin = al
			nextAlarmIndex = i
		}
	}
	nextAlarm = tempMin
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

//getUserResponse echoes a string to the "data" field and returns the user's response
func getUserResponse(s string) string {
	inputChan := make(chan string)
	defer func() {
		inputChan = nil //redirect input to regular handling function after we're done
	}()
	writeData([]string{s})
	return <-inputChan
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("title", 0, 0, maxX-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Wrap = true
		fmt.Fprint(v, banner)
	}
	if v, err := g.SetView("left-bg", 0, 3, maxX/2, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Wrap = true
		v.BgColor = gocui.ColorBlue
	}
	if v, err := g.SetView("right-bg", maxX/2, 3, maxX-1, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Wrap = true
		v.BgColor = gocui.ColorRed
	}
	if v, err := g.SetView("time", maxX/2-5, maxY/2-2, maxX/2+4, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
	}
	g.SetViewOnTop("time")
	if v, err := g.SetView("data", 0, maxY-8, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Wrap = true
		v.Autoscroll = true
	}

	if v, err := g.SetView("input", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Frame = false

		if _, err = setCurrentViewOnTop(g, "input"); err != nil {
			return err
		}

		inReader = bufio.NewReader(v)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func getUserInput(g *gocui.Gui, v *gocui.View) error {
	s := v.Buffer()
	v.Clear()
	err := v.SetCursor(0, 0)
	if err != nil {
		writeData([]string{fmt.Sprint(err)})
	}
	if len(s) <= 1 {
		return nil //user just pressed enter for whatever reason, ignore it
	}
	s = s[:len(s)-1] //shave off newline

	if inputChan != nil {
		inputChan <- s
		return nil
	}

	handleInput(s)
	//writeData([]string{s})
	//^ HANDLE USER INPUT INSTEAD

	return nil
}

func handleInput(s string) {
	args := strings.Split(s, " ")
	switch args[0] {
	case "add-alarm":
		alarmVars, err := parseCmd.ParseAlarm(s)
		if err != nil {
			if err == flag.ErrHelp {
				writeData([]string{"Required arguments: 'date': date to go off in form yyyy-mm-ddThh:MM:ss, with asterisks (*) meaning 'every'. Optional arguments: 'weekdays': space separated list of weekdays the alarm should go off on."})
				return
			}
			writeData([]string{"Error parsing arguments"})
			return
		}
		a, err := alclock.NewAlarm(alarmVars.DateString, strings.Split(alarmVars.Weekdays, " "))
		if err != nil {
			if err == alclock.ErrDateString {
				writeData([]string{"add-alarm: date string improperly formatted"})
				return
			}
			writeData([]string{"add-alarm: unknown error"})
		}
		alarmList = append(alarmList, a)
		updateNextAlarm()
		//DEBUG
		writeData([]string{fmt.Sprintf("Length of alarmList is: %d", len(alarmList))})
		if len(alarmList) > 0 {
			writeData([]string{fmt.Sprintf("Unix time of next alarm: %d", alarmList[0].NextGoesOff.Unix())})
		}
	case "add-arrive-by":
		arriveVars, err := parseCmd.ParseArriveBy(s)
		if err != nil {
			if err == flag.ErrHelp {
				writeData([]string{"Required arguments: 'date': date when you want to arrive in form yyyy-mm-ddThh:MM:ss, with asterisks (*) meaning 'every', 'start': starting location, 'end': destination. Optional arguments: 'weekdays': space separated list of weekdays the alarm should go off on, 'avoid': pipe (|) separated list of things to avoid (valid values are 'ferries', 'tolls', and 'highways', default is 'tolls|ferries'."})
				return
			}
			writeData([]string{"Error parsing arguments"})
			return
		}
		a, err := alclock.NewArriveBy(
			arriveVars.DateString,
			arriveVars.Origin,
			arriveVars.Destination,
			arriveVars.Avoid,
			strings.Split(arriveVars.Weekdays, " "),
		)
		if err != nil {
			if err == alclock.ErrDateString {
				writeData([]string{"add-alarm: date string improperly formatted"})
				return
			}
			writeData([]string{"add-alarm: unknown error"})
		}
		alarmList = append(alarmList, a)
		updateNextAlarm()
		//DEBUG
		writeData([]string{fmt.Sprintf("Length of alarmList is: %d", len(alarmList))})
		if len(alarmList) > 0 {
			writeData([]string{fmt.Sprintf("Unix time of next alarm: %d", alarmList[0].NextGoesOff.Unix())})
		}
	}
}

func main() {
	var g *gocui.Gui
	if useGUI {
		var err error
		g, err = gocui.NewGui()
		mainGUI = g
		if err != nil {
			return
		}
		defer g.Close()

		g.Highlight = true
		g.Cursor = true
		g.SelFgColor = gocui.ColorGreen

		g.SetManagerFunc(layout)

		if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
			return
		}

		if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, getUserInput); err != nil {
			return
		}
	}

	//timer
	go func() {
		for {
			t := time.Now()
			updateTime(t)
			go func() {
				//writeData([]string{fmt.Sprintf("Current Unix time: %d", t.Unix())})
				if len(alarmList) > 0 && t.Unix() == nextAlarm.NextGoesOff.Unix() {
					writeData([]string{"ALARM!"})
					if useGUI {
						audioSig := make(chan byte, 1)
						if err := g.SetKeybinding("input",
							gocui.KeySpace,
							gocui.ModNone,
							func(g *gocui.Gui, v *gocui.View) error {
								audioSig <- 1
								return nil
							}); err != nil {
							return
						}
						for {
							err := audio.PlayFile(config.Cfg.AudioFilePath, audioSig)
							if err == audio.ErrInterrupt {
								break
							}
						}
						g.DeleteKeybinding("input", gocui.KeySpace, gocui.ModNone)
					}
					go updateWeather() //make sure the weather is up-to-date after each alarm
					newTime := alclock.NextRing(nextAlarm)
					if newTime.Equal(nextAlarm.NextGoesOff) {
						//alarm does not repeat, delete it
						alarmList = append(alarmList[:nextAlarmIndex], alarmList[nextAlarmIndex+1:]...)
					}
					updateNextAlarm()
				}
			}()
			time.Sleep(1 * time.Second)
		}
	}()

	//calendar

	go func() {
		for {
			c := getCalendarUpdate()
			writeData(c)
			time.Sleep(time.Duration(config.Cfg.CalUpdatePeriod) * time.Minute)
		}
	}()

	//maps
	go func() {
		for {
			if len(alarmList) > 0 && alarmList[0].TripInfo != nil {
				alarmList[0].UpdateArriveBy()
			}
			time.Sleep(time.Duration(config.Cfg.MapUpdatePeriod) * time.Minute)
		}
	}()

	//weather
	go func() {
		for {
			updateWeather()
			time.Sleep(time.Duration(config.Cfg.WeatherUpdatePeriod) * time.Minute)
		}
	}()

	if useGUI {
		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			fmt.Println(err)
			return
		}
	} else {
		for {
			time.Sleep(30 * time.Second)
		}
	}
}

func getCalendarUpdate() []string {
	var retVal = []string{"Calendar updates:"}

	calendarEvents := cal.GetEvents(getUserResponse)
	for _, e := range calendarEvents {
		retVal = append(retVal, e.Summary)
	}
	return retVal
}

func updateWeather() {
	writeData([]string{fmt.Sprintf("GetWeather(%d, %s, %s)", config.Cfg.WeatherLookAhead, config.Cfg.Location[0], config.Cfg.Location[1])})
	wData, err := weather.GetNHours(config.Cfg.WeatherLookAhead, config.Cfg.Location[0], config.Cfg.Location[1])
	if err != nil {
		writeData([]string{fmt.Sprintf("Error getting weather data: %#v", err)})
		return
	}
	writeData([]string{fmt.Sprintf("Current chance of precipitation: %v", wData[0].PrecipProbability)})
	//write the current weather to the left panel and the forecast to the right panel
	if useGUI {
		mainGUI.Execute(func(g *gocui.Gui) error {
			v, err := g.View("left-bg")
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintf(v, "Current chance of precipitation: %v", wData[0].PrecipProbability)
			v, err = g.View("right-bg")
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintf(v, "Chance of precipitation %d hours from now: %v", config.Cfg.WeatherLookAhead, wData[1].PrecipProbability)
			return nil
		})
	}
}

func updateTime(t time.Time) {
	if useGUI {
		mainGUI.Execute(func(g *gocui.Gui) error {
			v, err := g.View("time")
			if err != nil {
				return err
			}
			v.Clear()
			fmt.Fprintf(v, "%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
			return nil
		})
	}
}

func writeData(sArr []string) {
	dataFieldMut.Lock()
	for _, s := range sArr {
		if useGUI {
			mainGUI.Execute(func(g *gocui.Gui) error {
				v, err := g.View("data")
				if err != nil {
					return err
				}

				fmt.Fprintln(v, s)
				return nil
			})
		} else {
			fmt.Println(s)
		}
		time.Sleep(500 * time.Millisecond)
	}
	dataFieldMut.Unlock()
}
