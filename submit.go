package trojanTechTickets

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func submitHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	}
	// Get the category.
	category := r.FormValue("category")
	if category == "Other" {
		category = r.FormValue("categoryOther")
	}
	// Parse grade into an int.
	grade, err := strconv.Atoi(r.FormValue("customerGrade"))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build ticket
	ticket := Ticket{
		Category: category,
		ClassName: r.FormValue("className"),
		ClassType: r.FormValue("classType"),
		ClassTeacher: r.FormValue("classTeacher"),
		Customer: Customer{
			Email: strings.ToLower(r.FormValue("customerEmail")),
			Grade: grade,
			Name:  r.FormValue("customerName"),
		},
		Description: r.FormValue("description"),
		EscalatedTo: r.FormValue("escalatedTo"),
		Time:        time.Now(),
		Worker:      strings.ToLower(u.String()),
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
	id, _, err := datastore.AllocateIDs(c, "Ticket", ticketKey(c), 1)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	key := datastore.NewKey(c, "Ticket", "", id, ticketKey(c))
	_, err = datastore.Put(c, key, &ticket)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirect to the ticket page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
