package trojanTechTickets

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func newWorkerHandler(w http.ResponseWriter, r *http.Request) {
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
	} else if !u.Admin {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}
	// Parse block and term into ints.
	block, err := strconv.Atoi(r.FormValue("block"))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	term, err := strconv.Atoi(r.FormValue("term"))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Build worker
	worker := Worker{
		Block: block,
		Email: strings.ToLower(r.FormValue("email")),
		Name:  r.FormValue("name"),
		Term:  term,
	}
	// Put it in the datastore
	key := datastore.NewIncompleteKey(c, "Worker", workerKey(c))
	_, err = datastore.Put(c, key, &worker)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirect to the ticket page
	http.Redirect(w, r, "/workers", http.StatusSeeOther)
}
