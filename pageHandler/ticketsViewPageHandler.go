package pageHandler

import (
	"fmt"
	"html/template"
	"net/http"
	"ticketBackend/sessionHandler"
	"ticketBackend/ticket"
)

// localhost:.../ticketsView
//anzeigen der Tickets des Users
func TicketsViewPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {

		var templateFiles []string
		templateFiles = append(templateFiles, "./assets/html/ticketsViewTemplate.html")
		templateFiles = append(templateFiles, "./assets/html/ticketListTemplate.html")

		templates, err := template.ParseFiles(templateFiles...)
		if err != nil {
			fmt.Println(err)
		}

		templates.ExecuteTemplate(response, "outer", nil)

		pTickets := *ticket.GetTickets(ticket.Open)

		for i := 0; i < len(pTickets); i++ {
			templates.ExecuteTemplate(response, "inner", pTickets[i])
		}

		templates.Execute(response, nil)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}
