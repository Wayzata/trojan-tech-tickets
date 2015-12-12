package trojanTechTickets

import (
	"log"
	"net/http"
	"strconv"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func showTicket(w http.ResponseWriter, r *http.Request) {
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
	} else if !u.Admin {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}
	// Make the logout URL
	logoutURL, err := user.LogoutURL(c, r.URL.String())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get the ticket number
	ticketNum, err := strconv.Atoi(r.URL.Path[9:])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get the tickets from the datastore
	ticket := make([]Ticket, 0, 1)
	_, err = datastore.NewQuery("Ticket").
		Ancestor(ticketKey(c)).
		Filter("Number =", ticketNum).
		Limit(1).
		GetAll(c, &ticket)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(ticket) == 0 {
		http.Error(w, "Unknown ticket", http.StatusNotFound)
		return
	}
	// Get the worker from the datastore
	worker := make([]Worker, 0, 1)
	_, err = datastore.NewQuery("Worker").
		Ancestor(workerKey(c)).
		Filter("Email =", ticket[0].Worker).
		Limit(1).
		GetAll(c, &worker)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(worker) != 0 {
		ticket[0].Worker = worker[0].Name
	}
	// Send the page
	err = templates.ExecuteTemplate(w, "showTicket.html", struct {
		BaseTemplateData
		Ticket Ticket
	}{
		BaseTemplateData: BaseTemplateData{
			CurrentURL:  r.URL.Path,
			LogoutURL:   logoutURL,
			User:        u.String(),
			UserIsAdmin: u.Admin,
		},
		Ticket: ticket[0],
	})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
