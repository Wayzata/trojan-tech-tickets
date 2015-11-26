package trojanTechTickets

import (
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func submitHandler(w http.ResponseWriter, r *http.Request) {
	// Check for user login
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	// Create the data item
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
		return
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
		return
	}
	// Redirect
	http.Redirect(w, r, "/ticket", http.StatusFound)
}
