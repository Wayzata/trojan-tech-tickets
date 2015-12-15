package trojanTechTickets

import (
	"encoding/json"
	"html/template"
	"strconv"
	"time"
)

var templates = template.Must(template.New("").Funcs(TemplateHelpers).ParseGlob("views/*.html"))

var TemplateHelpers = template.FuncMap{
	"gradeString": func(grade int) string {
		switch grade {
		case 9:
			return "Freshman"
		case 10:
			return "Sophomore"
		case 11:
			return "Junior"
		case 12:
			return "Senior"
		default:
			return strconv.Itoa(grade)
		}
	},
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
