package trojanTechTickets

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/user"
)

var ticketTemplate = template.Must(template.ParseFiles("assets/ticket.html"))

func ticketHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	checkError(w, ticketTemplate.Execute(w, struct {
		User string
	}{
		User: u.String(),
	}))
}
