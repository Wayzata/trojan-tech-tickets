package trojanTechTickets

import (
	"log"
	"net/http"
	"strings"
	urlPkg "net/url"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func deleteWorkerHandler(w http.ResponseWriter, r *http.Request) {
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
	// Parse the URL
	pathParts := strings.Split(r.RequestURI, "/")
	email, err := urlPkg.QueryUnescape(pathParts[len(pathParts)-1])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Find the worker
	keys, err := datastore.NewQuery("Worker").
		Ancestor(workerKey(c)).
		Filter("Email =", email).
		Limit(1).
		KeysOnly().
		GetAll(c, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// LART the sucker
	if err := datastore.DeleteMulti(c, keys); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirect to the workers page
	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}
