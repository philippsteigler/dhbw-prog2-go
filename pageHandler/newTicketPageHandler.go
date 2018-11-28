package pageHandler

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"de/vorlesung/projekt/crew/ticket"
	"fmt"
	"html/template"
	"net/http"
)

type Handler interface {
	ServeHTTP(http.Response, *http.Request)
}

// A-5.1:
// Uber eine Web-Seite soll ein Ticket erstellt werden können.
//
// https://localhost:8000/newTicketView
// anzeigen der Seite für die erstellung eines neuen Tickets

// A-5.2:
// Die Erzeugung eines Tickets soll ohne eine Authentifizierung möglich sein.
//
// Sofern der Anwender nicht eingeloggt ist, kann er über die Startseite neue Tickets erstellen und einreichen.
// Nach erfolgreicher Authentifizierung als Editor findet beim Aufruf der Startseite eine Weiterleitung zum
// internen Bereich statt. Dies basiert auf der Evaluierung des Session-Cookies.
func NewTicketViewPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		http.Redirect(response, request, "/dashboard", 302)
	} else {
		var templateFiles []string
		// TODO: GetAssetsDir() für alle PageHandler
		templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/newTicketTemplates/newTicketViewHeaderCssTemplate.html")
		templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/newTicketTemplates/newTicketTemplate.html")
		templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/newTicketTemplates/newTicketViewFooterTemplate.html")

		templates, err := template.ParseFiles(templateFiles...)
		if err != nil {
			fmt.Println(err)
		}

		templates.ExecuteTemplate(response, "outer", nil)
		templates.ExecuteTemplate(response, "newTicket", nil)
		templates.ExecuteTemplate(response, "footer", nil)

		templates.Execute(response, nil)
	}

}

// https://localhost:8000/saveTicket
// Speichert den Text aus den Textareas in mail, subject, text
func TicketSafeHandler(response http.ResponseWriter, request *http.Request) {

	if sessionHandler.IsUserLoggedIn(request) {

		ticket.NewTicket(request.FormValue("ticketSubject"), request.FormValue("ticketMail"), request.FormValue("ticketText"))
		// Zurück zu der Ticketseite
		http.Redirect(response, request, "/", http.StatusFound)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}
