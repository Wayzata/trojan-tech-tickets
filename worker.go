package trojanTechTickets

import (
	"appengine"
	"appengine/datastore"
)

type Worker struct {
	Block int
	Email string
	Name  string
	Term  int
}

func (w *Worker) Stats(c *appengine.Context) (int, int, error) {
	x := 0
	n := 0
	iter := datastore.NewQuery("Ticket").
		Ancestor(ticketKey(*c)).
		Filter("Worker =", w.Email).
		Order("-Time").
		Run(*c)
	var err error
	var ticket Ticket
	for {
		_, err = iter.Next(&ticket)
		if err != nil {
			break
		}
		if ticket.EscalatedTo == "" {
			x++
		}
		n++
	}
	if err == datastore.Done {
		err = nil
	}
	return int(float64(x * 100) / float64(n)), n, err
}
