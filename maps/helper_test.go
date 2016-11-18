package gmaps

import (
	"fmt"
	"testing"
)

func TestDepart(t *testing.T) {
	InitAPIKey()
	var params = map[string]string{
		"wp.0":  "Fairfax, VA",
		"wp.1":  "Washington, DC",
		"avoid": "tolls",
	}
	fmt.Println(GetTimeToLocation(params))
}
