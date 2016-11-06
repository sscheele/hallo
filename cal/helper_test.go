package cal

import (
	"fmt"
	"strings"
	"testing"
)

func TestCal(t *testing.T) {
	srv, err := GetCalendar()
	if err != nil {
		fmt.Printf("Unable to retrieve calendar Client %v\n", err)
		t.Error(err)
	}

	events, err := RetrieveEvents(srv, 10)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Printf("No upcoming events found.\n")
		return
	}
	var (
		isAllDay bool
		hour     int
		minute   int
		second   int
	)
	for _, i := range events.Items {
		var when string
		// If the DateTime is an empty string the Event is an all-day Event.
		// So only Date is available.
		if i.Start.DateTime != "" {
			when = i.Start.DateTime
			//Dates are formatted according to RFC3339
			weNeed := strings.Split(strings.Split(when, "T")[1], "-")[0]
			fmt.Fscanf(strings.NewReader(weNeed), "%d:%d:%d", &hour, &minute, &second)
			isAllDay = false
			fmt.Printf("Hour: %d, Minute: %d, Second: %d ", hour, minute, second)
		} else {
			when = i.Start.Date
			isAllDay = true
		}
		fmt.Printf("Is all day: %v ", isAllDay)
		fmt.Printf("%s (%s)\n", i.Summary, when)
	}
}
