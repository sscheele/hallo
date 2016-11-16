package gmaps

import (
	"fmt"
	"testing"
)

func TestDepart(t *testing.T) {
	InitAPIKey()
	var params = map[string]string{
		"origin":       "Fairfax, VA",
		"destination":  "Washington, DC",
		"arrival_time": "1479406158",
		"avoid":        "",
	}
	fmt.Println(GetTimeToLocation(params))
}
