package trojanTechTickets

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/user"
)

var ticketTemplate = template.Must(template.New("ticket.html").ParseFiles("assets/ticket.html"))

func ticket(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
	logoutURL, err := user.LogoutURL(c, r.URL.String())
	if checkError(w, err) { return }
	checkError(w, ticketTemplate.Execute(w, BaseTemplateData{
		LogoutURL: logoutURL,
		User: u.String(),
	}))
}
