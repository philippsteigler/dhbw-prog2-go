package pageHandler

import (
	"fmt"
	"html/template"
	"net/http"
	"ticketBackend/sessionHandler"
	"ticketBackend/ticket"
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
		templateFiles = append(templateFiles, "./assets/html/ticketsTemplates/ticketsViewHeaderCssTemplate.html")
		templateFiles = append(templateFiles, "./assets/html/ticketsTemplates/ticketsTicketListTemplate.html")
		templateFiles = append(templateFiles, "./assets/html/ticketsTemplates/ticketsViewFooterTemplate.html")

		templates, err := template.ParseFiles(templateFiles...)
		if err != nil {
			fmt.Println(err)
		}

		pTickets := *ticket.GetTickets(ticket.Open)

		templates.ExecuteTemplate(response, "outer", sessionHandler.GetSessionUser().Username)
		templates.ExecuteTemplate(response, "inner", pTickets)
		templates.ExecuteTemplate(response, "footer", nil)

		templates.Execute(response, nil)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}
