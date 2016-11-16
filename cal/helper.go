package cal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

//Event contains a description and time for an event
type Event struct {
	Summary  string
	DateTime EventTime
}

//EventTime contains the time of an event
//In the case of an  all-day event, hour, minute, and second will be empty strings
//Days and months should always be two digits (ie, March is 03)
type EventTime struct {
	Day      string
	Month    string
	Year     string
	Hour     string
	Minute   string
	Second   string
	IsAllDay bool
}

//NewEventTime returns an EventTime with default values
func NewEventTime() EventTime {
	return EventTime{
		Day:      "",
		Month:    "",
		Year:     "",
		Hour:     "",
		Minute:   "",
		Second:   "",
		IsAllDay: false,
	}
}

//GetEvents returns the next ten calendar events
func GetEvents(f func(string) string) (retVal []Event) {
	srv, err := GetCalendar(f)
	if err != nil {
		return
	}

	events, err := RetrieveEvents(srv, 10)
	if err != nil {
		return
	}

	if len(events.Items) == 0 {
		return
	}

	for _, i := range events.Items {
		var when string
		t := NewEventTime()
		// If the DateTime is an empty string the Event is an all-day Event.
		// So only Date is available.
		if i.Start.DateTime != "" {
			when = i.Start.DateTime
			//Dates are formatted according to RFC3339
			when = when[:19]
			fmt.Fscanf(strings.NewReader(when), "%s-%s-%sT%s:%s:%s", &t.Year, &t.Month, &t.Day, &t.Hour, &t.Minute, &t.Second)
			t.IsAllDay = false
			//fmt.Printf("Hour: %d, Minute: %d, Second: %d ", hour, minute, second)
		} else {
			when = i.Start.Date
			fmt.Fscanf(strings.NewReader(when), "%s-%s-%s", &t.Year, &t.Month, &t.Day)
			t.IsAllDay = true
		}
		retVal = append(retVal, Event{Summary: i.Summary, DateTime: t})
	}

	return retVal
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config, f func(string) string) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config, f)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config, f func(string) string) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	code := f(fmt.Sprintf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL))

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("calendar-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

//GetCalendar authenticates to the google api and returns a calendar service object
func GetCalendar(f func(string) string) (*calendar.Service, error) {
	ctx := context.Background()

	b, err := ioutil.ReadFile("/home/sam/Projects/Go/Gopath/src/github.com/sscheele/hallo/cal/client_secret.json")
	if err != nil {
		//log.Fatalf("Unable to read client secret file: %v", err)
		return nil, err
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/calendar-go-quickstart.json
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		//log.Fatalf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}
	client := getClient(ctx, config, f)

	srv, err := calendar.New(client)

	return srv, err
}

//RetrieveEvents accepts a calendar service object and returns a list of calendar events (up to MaxEvents long)
func RetrieveEvents(srv *calendar.Service, maxEvents int64) (list *calendar.Events, err error) {
	t := time.Now().Format(time.RFC3339)
	return srv.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(maxEvents).OrderBy("startTime").Do()
}
