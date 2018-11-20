package pageHandler

import (
	"net/http"
	"ticketBackend/sessionHandler"
	"ticketBackend/ticket"
)

type Handler interface {
	ServeHTTP(http.Response, *http.Request)
}

// localhost:.../newTicketView
//anzeigen der neues Ticket Seite
func NewTicketViewPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		http.ServeFile(response, request, "./assets/html/newTicketViewTemplate.html")
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// localhost:.../saveTicket
// Speichert den Text aus den Textareas in mail, subject, text
func TicketSafeHandler(response http.ResponseWriter, request *http.Request) {

	if sessionHandler.IsUserLoggedIn(request) {

		ticket.NewTicket(request.FormValue("ticketSubject"), request.FormValue("ticketMail"), request.FormValue("ticketText"))
		// Zurück zu der Ticketseite
		http.Redirect(response, request, "/ticket", http.StatusFound)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}
