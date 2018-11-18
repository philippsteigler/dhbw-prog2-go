package pageHandler

import (
	"../sessionHandler"
	"io/ioutil"
	"net/http"
)

type Handler interface {
	ServeHTTP(http.Response, *http.Request)
}

//NewTicketViewPageHandler
func NewTicketViewPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		http.ServeFile(response, request, "./assets/html/newTicketView.html")
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// localhost:.../saveTicket
// Speichert den Text aus den Textareas in mail, subject, text
func TicketSafeHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		mail := request.FormValue("ticketMail")
		subject := request.FormValue("ticketSubject")
		text := request.FormValue("ticketText")

		// Testdatei für die Eingabe
		inputTest := string("Mail\n" + mail + "\n\nSubject\n" + subject + "\n\nText\n" + text)
		ioutil.WriteFile("./assets/TicketTest", []byte(inputTest), 0600)

		// Zurück zu der Ticketseite
		http.Redirect(response, request, "/ticket", http.StatusFound)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}
