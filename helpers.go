package trojanTechTickets

import (
	"time"

	"appengine"
	"appengine/datastore"
)

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
