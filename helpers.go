package trojanTechTickets

import (
	"encoding/json"
	"html/template"
	"time"

	"appengine"
	"appengine/datastore"
)

type BaseTemplateData struct {
	LogoutURL string
	User string
}

type Step struct {
	Action     string
	Diagnostic string
}

type Ticket struct {
	Customer    string
	Description string
	EscalatedTo string
	Steps       []Step
	Time        time.Time
	Worker      string
}

func ticketKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Ticket", "default", 0, nil)
}

var TemplateHelpers = template.FuncMap{
	"json": func(v interface{}) (string, error) {
		b, err := json.Marshal(v)
		return string(b), err
	},
	"timeFormat": func(format string, time time.Time) string {
		return time.Format(format)
	},
	"truncate": func(maxLen int, str string) string {
		if len(str) < maxLen {
			return str
		}
		return str[:maxLen - 3] + "..."
	},
}
