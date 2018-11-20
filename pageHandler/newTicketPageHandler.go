package pageHandler

import (
	"../sessionHandler"
	"fmt"
	"html/template"
	"io/ioutil"
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
// Keine Authentifizierung nötig
// die anmeldung wird nicht abgefragt (kein: if sessionHandler.IsUserLoggedIn(request) {})
func NewTicketViewPageHandler(response http.ResponseWriter, request *http.Request) {
	var templateFiles []string
	templateFiles = append(templateFiles, "./assets/html/ticketInsightTemplates/ticketInsightViewHeaderCssTemplate.html") //TODO: Falsche datei, eigene erstellen
	templateFiles = append(templateFiles, "./assets/html/newTicketViewTemplate.html")

	templates, err := template.ParseFiles(templateFiles...)
	if err != nil {
		fmt.Println(err)
	}

	templates.ExecuteTemplate(response, "outer", sessionHandler.GetSessionUserName(request))
	templates.ExecuteTemplate(response, "newTicket", nil)

	templates.Execute(response, nil)
}

// https://localhost:8000/saveTicket
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
