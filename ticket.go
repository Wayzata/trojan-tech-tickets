package trojanTechTickets

import (
	"log"
	"net/http"

	"appengine"
	"appengine/user"
)

func ticketHandler(w http.ResponseWriter, r *http.Request) {
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
	logoutURL, err := user.LogoutURL(c, r.URL.String())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates.ExecuteTemplate(w, "ticket.html", BaseTemplateData{
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
