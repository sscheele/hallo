package main

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"

	"fmt"
	"time"
)

const numDataRows = 4

var banner = []string{" |_| _ | | _ ", " | |(_|| |(_)"}

type Alarm struct {
	//DateTime should have the following indices: "year", "month", "day", "hour", "minute", "second"
	//An asterisk ("*") matches everything
	DateTime map[string]string
	//DayOfWeek contains a key for each day of the week
	DayOfWeek map[time.Weekday]struct{}
}

var Alarms []Alarm
var ch chan string

func main() {
	//initScreen()

	ch = make(chan string)
	runClock(ch)
}

func runClock(ch chan string) {
	dataStrs := make([]string, numDataRows)
	for {
		t := time.Now()
		select {
		case x, ok := <-ch:
			if ok {
				dataStrs = append(dataStrs[1:], x)
			} else {
				return
			}
		default:

		}

		updateScreen(fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second()), dataStrs)
		time.Sleep(1 * time.Second)
	}
}

func moveCursorTo(x, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

func initScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	/*
		moveCursorTo(1, 1)
		fmt.Print(banner[0])
		moveCursorTo(1, 2)
		fmt.Print(banner[1])
	*/
}

func updateScreen(timeField string, dataFields []string) {
	width, height, err := termSize()

	x, _ := cursorPos(height, width) //cursorPos(height, width)

	if err != nil {
		fmt.Println("Error getting terminal size")
	}
	moveCursorTo(width/2, height/2)
	fmt.Print(timeField)

	for i, s := range dataFields {
		moveCursorTo(0, (height-numDataRows)+i)
		fmt.Print(s)
	}

	moveCursorTo(x, height)
}

func termSize() (x int, y int, err error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return
	}
	fmt.Fscanf(bytes.NewReader(out), "%d %d", &y, &x)
	return
}

func cursorPos(height, width int) (x int, err error) {
	if err != nil {
		return
	}
	enableRaw := exec.Command("stty", "raw")
	enableRaw.Stdin = os.Stdin
	enableRaw.Run()

	cmd := exec.Command("echo", "-e", fmt.Sprintf("%c[6n", 27))
	randomBytes := &bytes.Buffer{}
	cmd.Stdout = randomBytes

	// Start command asynchronously
	_ = cmd.Start()

	// capture keyboard output from echo command
	reader := bufio.NewReader(os.Stdin)
	cmd.Wait()

	// by printing the command output, we are triggering input
	fmt.Print(randomBytes)
	text, _ := reader.ReadString('R')
	text = trimPreJunk(text)

	disableRaw := exec.Command("stty", "-raw")
	disableRaw.Stdin = os.Stdin
	disableRaw.Run()

	for i := 3; i < height; i++ {
		moveCursorTo(0, i)
		fmt.Printf(fmt.Sprintf("%% -%ds", width), "")
	}
	/*
		moveCursorTo(1, 1)
		fmt.Print(banner[0])
		moveCursorTo(1, 2)
		fmt.Print(banner[1])
	*/

	fmt.Fscanf(strings.NewReader(text), "%dR", &x)
	moveCursorTo(x, height)
	fmt.Printf(fmt.Sprintf("%% -%ds", width-x), "")

	go func() {
		ch <- fmt.Sprintf("Cursor pos is (%d, %d)", x, height)
	}()

	return
}

func trimPreJunk(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == ';' {
			return s[i+1:]
		}
	}
	return ""
}

func addAlarm(s string) {

}
