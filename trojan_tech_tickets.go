package trojanTechTickets

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ticket", ticketHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/submit", submitHandler)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	// c := appengine.NewContext(r)
	fmt.Fprintf(w, "Hello (TODO Show the status of every ticket)")
}
