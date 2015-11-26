package trojanTechTickets

import "net/http"

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ticket", ticketHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/submit", submitHandler)
}
