package alclock

import (
	"fmt"
	"testing"
)

func TestNewAlarm(t *testing.T) {
	a, err := NewAlarm("2016-11-14T22:23:00", []string{""})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(a.NextGoesOff)
}
