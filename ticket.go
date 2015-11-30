package trojanTechTickets

import (
	"net/http"

	"appengine"
	"appengine/user"
)

func ticket(w http.ResponseWriter, r *http.Request, c appengine.Context, u *user.User) {
	logoutURL, err := user.LogoutURL(c, r.URL.String())
	if checkError(w, err) {
		return
	}
	checkError(w, templates.ExecuteTemplate(w, "ticket.html", BaseTemplateData{
		LogoutURL: logoutURL,
		User:      u.String(),
	}))
}
