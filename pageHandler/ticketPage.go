package pageHandler

import (
	"net/http"
	"ticketBackend/sessionHandler"
	"ticketBackend/ticket"
)

//localhost:.../saveTicket
//Speichert den Text aus den Textareas in mail, subject, text
func SaveTicketHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		// Neues Ticket erzeugen
		ticket.NewTicket(request.FormValue("ticketSubject"), request.FormValue("ticketMail"), request.FormValue("ticketText"))

		//Zurück zu der Ticketseite
		http.Redirect(response, request, "/ticket", http.StatusFound)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}
