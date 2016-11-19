package maps

import (
	"fmt"
	"testing"
)

func TestDepart(t *testing.T) {
	//Berkeley to San Francisco will pretty much always have traffic, so this is a good test to ensure traffic is being taken into account. The second number should be higher than the first.
	var params = map[string]string{
		"wp.0": "Berkeley, CA",
		"wp.1": "San Francisco, CA",
		//	"avoid": "tolls",
	}
	fmt.Println(GetTimeToLocation(params))
}
