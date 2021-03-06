package alclock

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/sscheele/hallo/maps"
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
	//Name is an optional parameter containing the name of the alarm (to make it easy to delete)
	Name string
}

//UpdateArriveBy updates the NextGoesOff atribute of an alarm to make sure it's in time for to arrive somewhere
func (a *Alarm) UpdateArriveBy() {
	t, err := strconv.ParseInt(a.TripInfo["arrival_time"], 10, 64)
	if err != nil {
		return
	}
	traff, noTraff := maps.GetTimeToLocation(a.TripInfo)
	if traff == 0 || noTraff == 0 {
		return
	}
	a.NextGoesOff = time.Unix(t-int64(traff), 0)
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
func NewArriveBy(arriveAt string, origin string, dest string, avoid string, weekdays []string, name string) (retVal Alarm, err error) {
	retVal, err = NewAlarm(arriveAt, weekdays, "")
	if err != nil {
		return
	}
	retVal.TripInfo = map[string]string{
		"wp.0":     origin,
		"wp.1":     dest,
		"dateTime": arriveAt,
		"avoid":    avoid,
	}
	timeUntil, _ := maps.GetTimeToLocation(retVal.TripInfo)
	retVal.NextGoesOff = time.Unix(retVal.NextGoesOff.Unix()-int64(timeUntil), 0)
	retVal.Name = name
	return
}

//NewAlarm gets a date string in a format similar to RFC3339: yyyy-mm-ddThh:MM:ss and an array of strings representing days of the week
//It returns an Alarm object based on these inputs
func NewAlarm(dateString string, weekdays []string, name string) (retVal Alarm, err error) {
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
		Name:      name,
	}
	retVal.NextGoesOff = NextRing(retVal)
	return
}

//NextRing returns a time giving the next instant the alarm will go off
//time instants can then be compared with .Before(), .Equal(), and .After()
func NextRing(a Alarm) time.Time {
	n := time.Now()
	var err error
	isEarly := false

	var se int
	if a.DateTime["second"] == "*" { //please do not do this
		//why would you do this
		se = n.Second() + 1
	} else {
		se, err = strconv.Atoi(a.DateTime["second"])
		if err != nil {
			return time.Unix(0, 0) //error value
		}
		if se < n.Second() {
			isEarly = true
		}
	}

	var mi int
	if a.DateTime["minute"] == "*" { //probably a bad idea
		mi = n.Minute()
		if isEarly {
			mi++
			isEarly = false
		}
	} else {
		mi, err = strconv.Atoi(a.DateTime["minute"])
		if err != nil {
			return time.Unix(0, 0) //error value
		}
		if mi < n.Minute() {
			isEarly = true
		}
	}

	var hr int
	if a.DateTime["hour"] == "*" {
		hr = n.Hour()
		if isEarly {
			hr++
			isEarly = false
		}
	} else {
		hr, err = strconv.Atoi(a.DateTime["hour"])
		if err != nil {
			return time.Unix(0, 0) //error value
		}
		if hr < n.Hour() {
			isEarly = true
		}
	}

	var da int
	if a.DateTime["day"] == "*" {
		da = n.Day()
		if isEarly {
			da++
			isEarly = false
		}
	} else if len(a.DayOfWeek) > 0 {
		var newStart time.Time
		for i := 0; i < 7; i++ { //when i reaches 7, we have exhausted every possible weekday
			poss := n.AddDate(0, 0, i)
			_, ok := a.DayOfWeek[poss.Weekday()]
			if ok {
				da = poss.Day()
				if i != 0 {
					isEarly = false
				}
				newStart = poss
			}
		}
		if isEarly {
			for i := 1; i < 7; i++ { //when i reaches 7, we have exhausted every possible weekday
				poss := newStart.AddDate(0, 0, i)
				_, ok := a.DayOfWeek[poss.Weekday()]
				if ok {
					da = poss.Day()
					if i != 0 {
						isEarly = false
					}
				}
			}
		}
		if isEarly {
			return time.Unix(0, 0) //if after this we're still too early, this is broken and we should error out
		}
	} else {
		da, err = strconv.Atoi(a.DateTime["day"])
		if err != nil {
			return time.Unix(0, 0) //error value
		}
		if da < n.Day() {
			isEarly = true
		}
	}

	var mo time.Month
	if a.DateTime["month"] == "*" {
		mo = n.Month()
		if isEarly {
			mo = time.Month(int(mo) + 1)
			isEarly = false
		}
	} else {
		tmp, err := strconv.Atoi(a.DateTime["month"])
		if err != nil {
			return time.Unix(0, 0) //error value
		}
		mo = time.Month(tmp) //months are 1+iota, so this should work
		if mo < n.Month() {
			isEarly = true
		}
	}

	var yr int
	if a.DateTime["year"] == "*" {
		yr = n.Year()
		if isEarly {
			yr++
			isEarly = false
		}
	} else {
		yr, err = strconv.Atoi(a.DateTime["year"])
		if err != nil {
			return time.Unix(0, 0) //error value
		}
		if yr < n.Year() {
			isEarly = true
		}
	}

	if isEarly { //YOU LIED TO ME
		//I THOUGHT YOU WERE THE CHOSEN ONE
		return time.Unix(0, 0)
	}

	return time.Date(yr, mo, da, hr, mi, se, 0, time.Local)
}
