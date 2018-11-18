package pageHandler

import (
	"../sessionHandler"
	"../ticket"
	"fmt"
	"html/template"
	"net/http"
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

		for i := 2; i <= len(*ticket.GetTickets(ticket.Open))+1; i++ {
			templates.ExecuteTemplate(response, "inner", ticket.GetTicket(i))
		}

		templates.Execute(response, nil)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}
