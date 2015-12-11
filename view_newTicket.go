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

func newTicket(w http.ResponseWriter, r *http.Request) {
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
	if r.Method == "POST" {
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
			Category:     category,
			ClassName:    r.FormValue("className"),
			ClassType:    r.FormValue("classType"),
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
		// Get the last ticket number
		highestTicket := make([]Ticket, 0, 1)
		_, err = datastore.NewQuery("Ticket").
			Ancestor(ticketKey(c)).
			Order("-Number").
			Limit(1).
			GetAll(c, &highestTicket)
		if len(highestTicket) == 0 {
			log.Println("1st ticket")
			highestTicket = append(highestTicket, Ticket{})
		} else {
			log.Println(highestTicket)
		}
		ticket.Number = highestTicket[0].Number + 1
		// Put it in the datastore
		key := datastore.NewIncompleteKey(c, "Ticket", ticketKey(c))
		_, err = datastore.Put(c, key, &ticket)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Redirect to the ticket page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// Make the logout URL
		logoutURL, err := user.LogoutURL(c, r.URL.String())
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Send the page
		err = templates.ExecuteTemplate(w, "newTicket.html", BaseTemplateData{
			CurrentURL:  r.URL.Path,
			LogoutURL:   logoutURL,
			User:        u.String(),
			UserIsAdmin: u.Admin,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
