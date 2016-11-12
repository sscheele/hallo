package gmaps

import (
	"fmt"
	"testing"
)

func TestDepart(t *testing.T) {
	InitAPIKey()
	fmt.Println(GetTimeToLocation("Fairfax, VA", "Washington, DC", "1480543258", ""))
}
