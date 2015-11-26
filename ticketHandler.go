package trojanTechTickets

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/user"
)

var ticketTemplate = template.Must(template.New("ticket.html").ParseFiles("assets/ticket.html"))

func ticketHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if checkError(w, err) { return }
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
	logoutURL, err := user.LogoutURL(c, r.URL.String())
	if checkError(w, err) { return }
	checkError(w, ticketTemplate.Execute(w, BaseTemplateData{
		LogoutURL: logoutURL,
		User: u.String(),
	}))
}
