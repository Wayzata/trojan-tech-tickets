package trojanTechTickets

import (
	"log"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func workersHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	}
	// Make the logout URL
	logoutURL, err := user.LogoutURL(c, r.URL.String())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get the tickets from the datastore
	query := datastore.NewQuery("Worker").Ancestor(workerKey(c)).Order("Name")
	workers := make([]Worker, 0, 100)
	_, err = query.GetAll(c, &workers)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Send the page
	err = templates.ExecuteTemplate(w, "workers.html", struct {
		BaseTemplateData
		Workers []Worker
	}{
		BaseTemplateData: BaseTemplateData{
			LogoutURL:   logoutURL,
			User:        u.String(),
			UserIsAdmin: u.Admin,
		},
		Workers: workers,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
