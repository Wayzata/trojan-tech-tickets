package trojanTechTickets

import "net/http"

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/new", newTicket)
	http.HandleFunc("/tickets", listTickets)
	http.HandleFunc("/tickets/", showTicket)
	http.HandleFunc("/workers", listWorkers)
	// http.HandleFunc("/workers/", showWorker)
}
