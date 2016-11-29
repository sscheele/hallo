package parseCmd

import (
	"bytes"
	"flag"

	"github.com/mattn/go-shellwords"
)

//AlarmVars will contain the parsed variables of the add-alarm command
type AlarmVars struct {
	DateString string
	Weekdays   string
	Name       string
}

//ArriveByVars will contain the parsed variables of the add-arrive-by command
type ArriveByVars struct {
	DateString  string
	Weekdays    string
	Origin      string
	Destination string
	Avoid       string
	Name        string
}

//ParseAlarm parses a command to add an alarm
func ParseAlarm(s string) (a AlarmVars, err error) {
	addAlarmFlags := flag.NewFlagSet("add-alarm", flag.ContinueOnError)
	addAlarmFlags.StringVar(&a.DateString, "date", "", "")
	addAlarmFlags.StringVar(&a.Weekdays, "weekdays", "", "")
	addAlarmFlags.StringVar(&a.Name, "name", "", "")
	//to prevent help message from breaking UI, output redirects to a junk buffer
	addAlarmFlags.SetOutput(bytes.NewBuffer([]byte{}))

	args, err := shellwords.Parse(s)
	if err != nil || len(args) < 2 {
		return
	}

	err = addAlarmFlags.Parse(args[1:])
	return
}

//ParseArriveBy parses a command to arrive at a time
func ParseArriveBy(s string) (a ArriveByVars, err error) {
	addArriveByFlags := flag.NewFlagSet("add-arrive-by", flag.ContinueOnError)
	addArriveByFlags.StringVar(&a.DateString, "date", "", "")
	addArriveByFlags.StringVar(&a.Weekdays, "weekdays", "", "")
	addArriveByFlags.StringVar(&a.Origin, "start", "", "")
	addArriveByFlags.StringVar(&a.Destination, "end", "", "")
	addArriveByFlags.StringVar(&a.Avoid, "avoid", "tolls", "")
	addArriveByFlags.StringVar(&a.Name, "name", "", "")
	//to prevent help message from breaking UI, output redirects to a junk buffer
	addArriveByFlags.SetOutput(bytes.NewBuffer([]byte{}))

	args, err := shellwords.Parse(s)
	if err != nil || len(args) < 2 {
		return
	}

	err = addArriveByFlags.Parse(args[1:])
	return
}
