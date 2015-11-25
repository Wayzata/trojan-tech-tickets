package ticket_system

import (
	"fmt"
	_ "html/template"
	"net/http"

	"appengine"
	"appengine/user"
)

func init() {
	http.HandleFunc("/", userHandler)
	// http.HandleFunc("/list", listHandler)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprintf(w, "Hello, %v!", u)
}
