package trojanTechTickets

import (
	"fmt"
	"net/http"

	"appengine"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found: %s", r.URL.String())
		return
	}
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	} else if r.Method == "POST" {
		if !submit(w, r, c, u) {
			ticket(w, r, c, u)
		}
	} else if u.Admin {
		status(w, r, c, u)
	} else {
		ticket(w, r, c, u)
	}
}
