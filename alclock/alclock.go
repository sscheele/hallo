package main

import (
	"bufio"
	"fmt"
	"sync"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/sscheele/hallo/cal"
)

const numDataRows = 4
const banner = ` |_| _ | | _
 | |(_|| |(_)`

type alarm struct {
	//DateTime should have the following indices: "year", "month", "day", "hour", "minute", "second"
	//An asterisk ("*") matches everything
	DateTime map[string]string
	//DayOfWeek contains a key for each day of the week
	DayOfWeek map[time.Weekday]struct{}
}

var (
	alarmList    []alarm
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

	writeData([]string{s})
	//^ HANDLE USER INPUT INSTEAD

	return nil
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
			timer()
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

	/*
		t := time.Now()
			day := fmt.Sprintf("%02d", t.Day())
			month := fmt.Sprintf("%02d", t.Month())
			year := fmt.Sprintf("%02d", t.Year())
	*/
	calendarEvents := cal.GetEvents(getUserResponse)
	for _, e := range calendarEvents {
		retVal = append(retVal, e.Summary)
	}
	return retVal
}

func timer() {
	mainGUI.Execute(func(g *gocui.Gui) error {
		t := time.Now()
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
