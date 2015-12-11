package trojanTechTickets

import (
	"encoding/json"
	"html/template"
	"time"
)

var templates = template.Must(template.New("").Funcs(TemplateHelpers).ParseGlob("views/*.html"))

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
	"workerName": func(email, workersJSON string) string {
		// TODO A less completely *stupid* way of passing around the workers
		// list than putting it to JSON and back.
		workers := make([]Worker, 0, 10)
		if err := json.Unmarshal([]byte(workersJSON), &workers); err != nil {
			return email
		}
		for _, worker := range workers {
			if worker.Email == email {
				return worker.Name
			}
		}
		return email
	},
}
