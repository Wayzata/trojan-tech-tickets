package trojanTechTickets

import (
	"log"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

func showWorker(w http.ResponseWriter, r *http.Request) {
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
	// Get the worker email
	email := r.URL.Path[9:]
	if len(email) == 0 {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	// Get the worker from the datastore
	worker := make([]Worker, 0, 1)
	keys, err := datastore.NewQuery("Worker").
		Ancestor(workerKey(c)).
		Filter("Email =", email).
		Limit(1).
		GetAll(c, &worker)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(worker) == 0 {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	// If GET
	if r.Method == "GET" {
		// Make the logout URL
		logoutURL, err := user.LogoutURL(c, r.URL.String())
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Get the number of tickets and resolve rate.
		resolveRate, ticketCount, err := (&worker[0]).Stats(&c)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Send the page
		err = templates.ExecuteTemplate(w, "showWorker.html", struct {
			BaseTemplateData
			ResolveRate int
			TicketCount int
			Worker *Worker
		}{
			BaseTemplateData: BaseTemplateData{
				CurrentURL:  r.URL.Path,
				LogoutURL:   logoutURL,
				User:        u.String(),
				UserIsAdmin: u.Admin,
			},
			ResolveRate: resolveRate,
			TicketCount: ticketCount,
			Worker: &worker[0],
		})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "DELETE" {
		if err := datastore.Delete(c, keys[0]); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
