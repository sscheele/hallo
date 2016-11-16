package alclock

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sscheele/hallo/gmaps"
)

//ErrDateString will be returned when trying to parse an invalid date string
var ErrDateString = errors.New("Error: date string invalid")

//Alarm represents an alarm
type Alarm struct {
	//DateTime should have the following indices: "year", "month", "day", "hour", "minute", "second"
	//An asterisk ("*") matches everything
	DateTime map[string]string
	//DayOfWeek contains a key for each day of the week
	DayOfWeek map[time.Weekday]struct{}
	//NextGoesOff contains the time at which the alarm will next go off
	NextGoesOff time.Time
	//TripInfo contains information about whether the alarm is an "in time for" alarm
	TripInfo map[string]string
}

//UpdateArriveBy updates the NextGoesOff atribute of an alarm to make sure it's in time for to arrive somewhere
func (a *Alarm) UpdateArriveBy() {
	var t int64
	fmt.Fscanf(strings.NewReader(a.TripInfo["arrival_time"]), "%d", &t)
	dur := gmaps.GetTimeToLocation(a.TripInfo)
	if dur == 0 {
		return
	}
	a.NextGoesOff = time.Unix(t-int64(dur), 0)
}

//EmptyAlarm returns an empty alarm which is about to go off
//primarily useful as a junk value
func EmptyAlarm() Alarm {
	return Alarm{
		DateTime:    nil,
		DayOfWeek:   nil,
		NextGoesOff: time.Now(),
		TripInfo:    nil,
	}
}

//NewArriveBy returns an alarm that changes as google changes its estimate for arrival time
func NewArriveBy(arriveAt string, origin string, dest string, avoid string, weekdays []string) (retVal Alarm, err error) {
	retVal, err = NewAlarm(arriveAt, weekdays)
	if err != nil {
		return
	}
	retVal.TripInfo = map[string]string{
		"origin":       origin,
		"destination":  dest,
		"arrival_time": fmt.Sprintf("%d", retVal.NextGoesOff.Unix()),
		"avoid":        avoid,
	}
	retVal.NextGoesOff = time.Unix(retVal.NextGoesOff.Unix()-int64(gmaps.GetTimeToLocation(retVal.TripInfo)), 0)
	return
}

//NewAlarm gets a date string in a format similar to RFC3339: yyyy-mm-ddThh:MM:ss and an array of strings representing days of the week
//It returns an Alarm object based on these inputs
func NewAlarm(dateString string, weekdays []string) (retVal Alarm, err error) {
	retVal = EmptyAlarm()
	var wd map[time.Weekday]struct{}

	dt := strings.FieldsFunc(dateString, func(r rune) bool {
		switch r {
		case '-', 'T', ':':
			return true
		}
		return false
	})

	if len(dt) != 6 {
		err = ErrDateString
		return
	}

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
	retVal = Alarm{
		DateTime:  map[string]string{"year": dt[0], "month": dt[1], "day": dt[2], "hour": dt[3], "minute": dt[4], "second": dt[5]},
		DayOfWeek: wd,
		TripInfo:  nil,
	}
	retVal.NextGoesOff = NextRing(retVal)
	return
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
