package main

import (
	"bytes"
	"os"
	"os/exec"

	"fmt"
	"time"
)

const numDataRows = 4
const banner = `
 |_| _ | | _
 | |(_|| |(_)
`

func main() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Print(banner)

	dataStrs := make([]string, numDataRows)
	for i := 0; ; i++ {
		t := time.Now()

		dataStrs = append(dataStrs[1:], fmt.Sprintf("%d", i))

		updateScreen(fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second()), dataStrs)
		time.Sleep(1 * time.Second)
	}
}

func updateScreen(timeField string, dataFields []string) {
	width, height, err := termSize()
	if err != nil {
		fmt.Println("Error getting terminal size")
	}
	fmt.Printf("\033[%d;%dH", height/2, width/2)
	fmt.Print(timeField)

	for i, s := range dataFields {
		fmt.Printf("\033[%d;0H", (height-numDataRows)+i+1)
		fmt.Printf(fmt.Sprintf("%% -%ds", width-1), s)
	}
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
