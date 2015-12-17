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

func listWorkers(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	} else if !u.Admin {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}
	// Switch functionality based on request type
	if r.Method == "POST" { // Adding new worker
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
	} else { // Getting list of workers
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
		// Get worker stats
		workerStats := make(map[string][]int)
		for _, worker := range workers {
			rr, tc, err := (&worker).Stats(&c)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			workerStats[worker.Email] = []int{rr, tc}
		}
		// Send the page
		err = templates.ExecuteTemplate(w, "listWorkers.html", struct {
			BaseTemplateData
			Stats map[string][]int
			Workers []Worker
		}{
			BaseTemplateData: BaseTemplateData{
				CurrentURL:  r.URL.Path,
				LogoutURL:   logoutURL,
				User:        u.String(),
				UserIsAdmin: u.Admin,
			},
			Stats: workerStats,
			Workers: workers,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
