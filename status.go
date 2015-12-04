package trojanTechTickets

import (
	"log"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
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
	// Get the tickets from the datastore
	query := datastore.NewQuery("Ticket").Ancestor(ticketKey(c)).Order("-Time")
	tickets := make([]Ticket, 0, 100)
	_, err = query.GetAll(c, &tickets)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Send the page
	err = templates.ExecuteTemplate(w, "status.html", struct {
		BaseTemplateData
		Tickets []Ticket
	}{
		BaseTemplateData: BaseTemplateData{
			LogoutURL:   logoutURL,
			User:        u.String(),
			UserIsAdmin: u.Admin,
		},
		Tickets: tickets,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
