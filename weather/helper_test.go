package weather

import (
	"fmt"
	"testing"
)

func TestNHours(t *testing.T) {
	dps, err := GetNHours(6, "38.8522392", "-77.3368576")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(dps[0].PrecipProbability, dps[1].PrecipProbability)
}
