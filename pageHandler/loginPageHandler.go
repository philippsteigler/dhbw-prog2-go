package pageHandler

import (
	"../sessionHandler"
	"net/http"
)

// A-8.1:
// Die Bearbeitung der Tickets soll ausschließlich ¨uber eine WEB-Seite erfolgen.
//
// Loginseite
// Anmeldung mit Benutzer und Passwort
func LoginPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {

		//Seite für den Angemeldeten User aufrufen
		http.Redirect(response, request, "/dashboard", 302)
	} else {

		//Neues Ticket falls kein User angemeldet ist
		http.ServeFile(response, request, "./assets/html/loginView.html")
	}
}
