package pageHandler

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"de/vorlesung/projekt/crew/ticket"
	"fmt"
	"html/template"
	"net/http"
)

//TODO: Selektor für die Ticket Anzeige

// A-8.1
// Die Bearbeitung der Tickets soll ausschließlich ¨uber eine WEB-Seite erfolgen.
//
// https://localhost:8000/ticketsView
//anzeigen der Tickets des Users
func TicketsViewPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {

		var templateFiles []string
		templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/ticketsTemplates/ticketsViewHeaderCssTemplate.html")
		templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/ticketsTemplates/ticketsTicketListTemplate.html")
		templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/ticketsTemplates/ticketsViewFooterTemplate.html")

		templates, err := template.ParseFiles(templateFiles...)
		if err != nil {
			fmt.Println(err)
		}

		pTickets := *ticket.GetAllOpenTickets()

		templates.ExecuteTemplate(response, "outer", sessionHandler.GetSessionUser(request).Username)
		templates.ExecuteTemplate(response, "inner", pTickets)
		templates.ExecuteTemplate(response, "footer", nil)

		templates.Execute(response, nil)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}
