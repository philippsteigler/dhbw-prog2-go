package pageHandler

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"de/vorlesung/projekt/crew/ticket"
	"fmt"
	"html/template"
	"net/http"
)

// Matrikelnummern:
//
// 3333958
// 3880065
// 8701350

var ticketsViewTemplates *template.Template

// A-8.1
// Die Bearbeitung der Tickets soll ausschließlich ¨uber eine WEB-Seite erfolgen.
//
// https://localhost:8000/ticketsView
// anzeigen der Tickets des Users
func TicketsViewPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		pTickets := *ticket.GetAllOpenTickets()

		ticketsViewTemplates.ExecuteTemplate(response, "outer", sessionHandler.GetSessionUser(request).Username)
		ticketsViewTemplates.ExecuteTemplate(response, "inner", pTickets)
		ticketsViewTemplates.ExecuteTemplate(response, "footer", nil)

		ticketsViewTemplates.Execute(response, nil)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func TicketsViewInit() {
	var templateFiles []string
	var err error

	templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/ticketsTemplates/ticketsViewHeaderCssTemplate.html")
	templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/ticketsTemplates/ticketsTicketListTemplate.html")

	ticketsViewTemplates, err = template.ParseFiles(templateFiles...)
	if err != nil {
		fmt.Println(err)
	}
}
