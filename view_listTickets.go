package trojanTechTickets

import (
	"log"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func listTickets(w http.ResponseWriter, r *http.Request) {
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
	// Make the logout URL
	logoutURL, err := user.LogoutURL(c, r.URL.String())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get the tickets from the datastore
	tickets := make([]Ticket, 0, 100)
	query := datastore.NewQuery("Ticket").
		Ancestor(ticketKey(c)).
		Order("-Time")
	if !u.Admin {
		query = query.Filter("Worker = ", u.Email)
	}
	_, err = query.GetAll(c, &tickets)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get the workers from the datastore
	workers := make([]Worker, 0, 10)
	_, err = datastore.NewQuery("Worker").
		Ancestor(workerKey(c)).
		GetAll(c, &workers)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Send the page
	err = templates.ExecuteTemplate(w, "listTickets.html", struct {
		BaseTemplateData
		Tickets []Ticket
		Workers []Worker
	}{
		BaseTemplateData: BaseTemplateData{
			CurrentURL:  r.URL.Path,
			LogoutURL:   logoutURL,
			User:        u.String(),
			UserIsAdmin: u.Admin,
		},
		Tickets: tickets,
		Workers: workers,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
