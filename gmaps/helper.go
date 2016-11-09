package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

func main() {
	f, err := os.Open("api-key.txt")
	if err != nil {
		fmt.Println("Error reading in API key!")
		return
	}
	reader := bufio.NewReader(f)
	text, _ := reader.ReadString('\n')
	c, err := maps.NewClient(maps.WithAPIKey(text))
	if err != nil {
		fmt.Println("fatal error: %s", err)
		return
	}
	r := &maps.DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Perth",
	}
	resp, _, err := c.Directions(context.Background(), r)
	if err != nil {
		fmt.Println("fatal error: %s", err)
		return
	}

	fmt.Println(resp)
}
