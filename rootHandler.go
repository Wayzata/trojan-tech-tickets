package trojanTechTickets

import (
	"net/http"

	"appengine"
	"appengine/user"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// 404 if not a request for the root
	if r.URL.String() != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Get the user
	c := appengine.NewContext(r)
	u := user.Current(c)
	// Find the correct redirect URL
	var url string
	if u == nil {
		var err error
		url, err = user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if u.Admin {
		url = "/status"
	} else {
		url = "/ticket"
	}
	// Redirect
	http.Redirect(w, r, url, http.StatusFound)
}
