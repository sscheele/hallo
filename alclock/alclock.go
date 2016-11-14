package alclock

import (
	"fmt"
	"strings"
	"time"
)

//Alarm represents an alarm
type Alarm struct {
	//DateTime should have the following indices: "year", "month", "day", "hour", "minute", "second"
	//An asterisk ("*") matches everything
	DateTime map[string]string
	//DayOfWeek contains a key for each day of the week
	DayOfWeek map[time.Weekday]struct{}
	//NextGoesOff contains the time at which the alarm will next go off
	NextGoesOff time.Time
}

//NewAlarm gets a date string in a format similar to RFC3339: yyyy-mm-ddThh:MM:ss and an array of strings representing days of the week
//It returns an Alarm object based on these inputs
func NewAlarm(dateString string, weekdays []string) Alarm {
	var (
		wd map[time.Weekday]struct{}
		yr string
		mo string
		da string
		ho string
		mi string
		se string
	)
	fmt.Fscanf(strings.NewReader(dateString), "%s-%s-%sT%s:%s:%s", &yr, &mo, &da, &ho, &mi, &se)
	for _, item := range weekdays {
		switch {
		case strings.EqualFold(item, "monday"):
			wd[time.Monday] = struct{}{}
		case strings.EqualFold(item, "tuesday"):
			wd[time.Tuesday] = struct{}{}
		case strings.EqualFold(item, "wednesday"):
			wd[time.Wednesday] = struct{}{}
		case strings.EqualFold(item, "thursday"):
			wd[time.Thursday] = struct{}{}
		case strings.EqualFold(item, "friday"):
			wd[time.Friday] = struct{}{}
		case strings.EqualFold(item, "saturday"):
			wd[time.Saturday] = struct{}{}
		case strings.EqualFold(item, "sunday"):
			wd[time.Sunday] = struct{}{}
		}
	}
	retVal := Alarm{
		DateTime:  map[string]string{"year": yr, "month": mo, "day": da, "hour": ho, "minute": mi, "second": se},
		DayOfWeek: wd,
	}
	retVal.NextGoesOff = NextRing(retVal)
	return retVal
}

//NextRing returns a time giving the next instant the alarm will go off
//time instants can then be compared with .Before(), .Equal(), and .After()
func NextRing(a Alarm) time.Time {
	n := time.Now()
	var yr int
	if a.DateTime["year"] == "*" {
		yr = n.Year()
	} else {
		fmt.Fscanf(strings.NewReader(a.DateTime["year"]), "%d", &yr)
	}

	var mo time.Month
	if a.DateTime["month"] == "*" {
		mo = n.Month()
	} else {
		var tmp int
		fmt.Fscanf(strings.NewReader(a.DateTime["month"]), "%d", &tmp)
		mo = time.Month(tmp) //months are 1+iota, so this should work
	}

	var da int
	_, isWeekday := a.DayOfWeek[n.Weekday()]
	if a.DateTime["day"] == "*" || isWeekday {
		da = n.Day()
	} else {
		fmt.Fscanf(strings.NewReader(a.DateTime["day"]), "%d", &da)
	}

	var hr int
	if a.DateTime["hour"] == "*" {
		hr = n.Hour()
	} else {
		fmt.Fscanf(strings.NewReader(a.DateTime["hour"]), "%d", &hr)
	}

	var mi int
	if a.DateTime["minute"] == "*" { //probably a bad idea
		mi = n.Minute()
	} else {
		fmt.Fscanf(strings.NewReader(a.DateTime["minute"]), "%d", &mi)
	}

	var se int
	if a.DateTime["second"] == "*" { //please do not do this
		//why would you do this
		se = n.Second()
	} else {
		fmt.Fscanf(strings.NewReader(a.DateTime["second"]), "%d", &se)
	}

	return time.Date(yr, mo, da, hr, mi, se, 0, time.Local)
}
