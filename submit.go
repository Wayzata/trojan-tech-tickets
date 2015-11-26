package trojanTechTickets

import (
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func submit(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) bool {
	ticket := Ticket{
		Customer:    r.FormValue("customer"),
		Description: r.FormValue("description"),
		EscalatedTo: r.FormValue("escalatedTo"),
		Time:        time.Now(),
		Worker:      u.String(),
	}
	// Put in the steps
	diagnostics := r.Form["stepDiagnostic[]"]
	actions := r.Form["stepAction[]"]
	if len(diagnostics) != len(actions) {
		http.Error(w, "Invalid steps.", http.StatusBadRequest)
		return true
	}
	ticket.Steps = make([]Step, len(diagnostics))
	for i, diagnostic := range diagnostics {
		action := actions[i]
		ticket.Steps[i] = Step{
			Action:     action,
			Diagnostic: diagnostic,
		}
	}
	// Put it in the datastore
	key := datastore.NewIncompleteKey(c, "Ticket", ticketKey(c))
	_, err := datastore.Put(c, key, &ticket)
	if checkError(w, err) {
		return true
	}
	// Return
	return false
}
