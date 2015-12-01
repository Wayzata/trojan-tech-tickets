package trojanTechTickets

import (
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func status(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
	// Make the logout URL
	logoutURL, err := user.LogoutURL(c, r.URL.String())
	if checkError(w, err) { return }
	// Get the tickets from the datastore
	query := datastore.NewQuery("Ticket").Ancestor(ticketKey(c)).Order("-Time")
	tickets := make([]Ticket, 0, 100)
	_, err = query.GetAll(c, &tickets)
	if checkError(w, err) { return }
	// Send the page
	checkError(w, templates.ExecuteTemplate(w, "status.html", struct {
		BaseTemplateData
		Tickets []Ticket
	}{
		BaseTemplateData: BaseTemplateData{
			LogoutURL: logoutURL,
			User: u.String(),
		},
		Tickets: tickets,
	}))
}
