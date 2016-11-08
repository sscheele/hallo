package main

import (
	"bytes"
	"os"
	"os/exec"

	"fmt"
	"time"
)

func main() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	width, height, err := termSize()
	if err != nil {
		fmt.Println("Error getting terminal size")
	}

	for {
		t := time.Now()
		fmt.Printf("\033[%d;%dH", width/2, height/2)
		fmt.Printf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
		time.Sleep(1 * time.Second)
	}
}

func termSize() (x int, y int, err error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return
	}
	fmt.Fscanf(bytes.NewReader(out), "%d %d", &x, &y)
	return
}
