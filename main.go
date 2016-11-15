package main

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/sscheele/hallo/alclock"
	"github.com/sscheele/hallo/audio"
	"github.com/sscheele/hallo/cal"
)

const numDataRows = 4
const banner = ` |_| _ | | _
 | |(_|| |(_)`

var (
	alarmList    []alclock.Alarm
	dataFieldMut sync.Mutex
	inReader     *bufio.Reader
	inputChan    chan string
	mainGUI      *gocui.Gui
)

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
	if v, err := g.SetView("time", maxX/2-5, maxY/2-2, maxX-1, maxY/2+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
	}
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
		//SYNTAX: add-alarm yyyy-mm-ddThh:MM:ss weekday1 weekday2 weekday3
		if len(args) < 2 {
			writeData([]string{"add-alarm: not enough arguments"})
			return
		}
		args = append(args, "") //avoid index out of bounds while calling NewAlarm
		a, err := alclock.NewAlarm(args[1], args[2:])
		if err != nil {
			if err == alclock.ErrDateString {
				writeData([]string{"add-alarm: date string improperly formatted"})
				return
			}
			writeData([]string{"add-alarm: unknown error"})
		}
		addAlarm(a)
		//DEBUG
		writeData([]string{fmt.Sprintf("Length of alarmList is: %d", len(alarmList))})
		if len(alarmList) > 0 {
			writeData([]string{fmt.Sprintf("Unix time of next alarm: %d", alarmList[0].NextGoesOff.Unix())})
		}
	}
}

func addAlarm(a alclock.Alarm) {
	//keep alarmList sorted so that alarmList[0] is always the next to go off
	if len(alarmList) == 0 {
		alarmList = append(alarmList, a)
		return
	}
	i := 0
	for i < len(alarmList) && a.NextGoesOff.After(alarmList[i].NextGoesOff) {
		i++ //could do this in one line, but it's more clear this way
	}
	if i == len(alarmList)-1 { //avoid needlessly copying
		alarmList = append(alarmList, a)
		return
	}
	alarmList = append(alarmList, alclock.EmptyAlarm())
	copy(alarmList[i+1:], alarmList[i:]) //golang's copy is dst, src
	alarmList[i] = a
}

func main() {
	g, err := gocui.NewGui()
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

	//timer
	go func() {
		for {
			t := time.Now()
			updateTime(t)
			go func() {
				//writeData([]string{fmt.Sprintf("Current Unix time: %d", t.Unix())})
				if len(alarmList) > 0 && t.Unix() == alarmList[0].NextGoesOff.Unix() {
					writeData([]string{"ALARM!"})
					sig := make(chan byte, 1)
					if err := g.SetKeybinding("input",
						gocui.KeySpace,
						gocui.ModNone,
						func(g *gocui.Gui, v *gocui.View) error {
							sig <- 1
							return nil
						}); err != nil {
						return
					}
					for {
						err := audio.PlayFile("/home/sam/Projects/Go/Gopath/src/github.com/sscheele/hallo/audio/bell.aiff", sig)
						if err == audio.ErrInterrupt {
							break
						}
					}
					g.DeleteKeybinding("input", gocui.KeySpace, gocui.ModNone)
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
			time.Sleep(15 * time.Minute) //TODO: Change this value for real purposes
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		fmt.Println(err)
		return
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

func updateTime(t time.Time) {
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

func writeData(sArr []string) {
	dataFieldMut.Lock()
	for _, s := range sArr {
		mainGUI.Execute(func(g *gocui.Gui) error {
			v, err := g.View("data")
			if err != nil {
				return err
			}

			fmt.Fprintln(v, s)
			return nil
		})
		time.Sleep(500 * time.Millisecond)
	}
	dataFieldMut.Unlock()
}
