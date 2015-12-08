package trojanTechTickets

import (
	"encoding/json"
	"html/template"
	"time"

	"appengine"
	"appengine/datastore"
)

type BaseTemplateData struct {
	LogoutURL   string
	User        string
	UserIsAdmin bool
}

type Customer struct {
	Email string
	Grade int
	Name  string
}

type Step struct {
	Action     string
	Diagnostic string
}

type Ticket struct {
	Category    string
	ClassName string
	ClassType string
	ClassTeacher string
	Customer    Customer
	Description string
	EscalatedTo string
	Steps       []Step
	Time        time.Time
	Worker      string
}

func ticketKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Ticket", "default", 0, nil)
}
func workerKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Worker", "default", 0, nil)
}

var templates = template.Must(template.New("").Funcs(TemplateHelpers).ParseGlob("assets/*.html"))

var TemplateHelpers = template.FuncMap{
	"json": func(v interface{}) (string, error) {
		b, err := json.Marshal(v)
		return string(b), err
	},
	"timeFormat": func(format string, time time.Time) string {
		return time.Format(format)
	},
	"timeIn": func(tz string, t time.Time) (time.Time, error) {
		location, err := time.LoadLocation(tz)
		if err != nil {
			return time.Time{}, err
		}
		return t.In(location), nil
	},
	"truncate": func(maxLen int, str string) string {
		if len(str) < maxLen {
			return str
		}
		return str[:maxLen-3] + "..."
	},
}
