package main

import (
	"fmt"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

func main() {
	c, err := maps.NewClient(maps.WithAPIKey(" AIzaSyCO8qw67aswZVM-3-GPIxPEBJwDIArxjsY"))
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
