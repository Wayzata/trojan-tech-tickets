package trojanTechTickets

import (
	"appengine"
	"appengine/datastore"
)

func ticketKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Ticket", "default", 0, nil)
}
func workerKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Worker", "default", 0, nil)
}
