package trojanTechTickets

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

var statusTemplate = template.Must(template.New("status.html").Funcs(TemplateHelpers).ParseFiles("assets/status.html"))

func statusHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if checkError(w, err) { return }
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
	// Make the logout URL
	logoutURL, err := user.LogoutURL(c, r.URL.String())
	if checkError(w, err) { return }
	// Get the tickets from the datastore
	query := datastore.NewQuery("Ticket").Ancestor(ticketKey(c)).Order("-Time")
	tickets := make([]Ticket, 0, 100)
	_, err = query.GetAll(c, &tickets)
	if checkError(w, err) { return }
	// Send the page
	checkError(w, statusTemplate.Execute(w, struct {
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
