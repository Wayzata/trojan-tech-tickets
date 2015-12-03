package trojanTechTickets

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/workers", workersHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found: %s", r.URL.String())
		return
	}
	if r.Method == "POST" {
		submitHandler(w, r)
	} else {
		ticketHandler(w, r)
	}
}
