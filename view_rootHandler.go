package trojanTechTickets

import (
	"net/http"

	"appengine"
	"appengine/user"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	} else {
		c := appengine.NewContext(r)
		u := user.Current(c)
		url := "/new"
		if u != nil && u.Admin {
			url = "/tickets"
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
