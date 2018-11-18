package pageHandler

import (
	"../sessionHandler"
	"fmt"
	"html/template"
	"net/http"
)

func TicketsViewPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {

		var templateFiles []string
		templateFiles = append(templateFiles, "./assets/html/ticketsViewTemplate.html")
		templateFiles = append(templateFiles, "./assets/html/ticketsViewTicketList.html")

		templates, err := template.ParseFiles(templateFiles...)
		if err != nil {
			fmt.Println(err)
		}

		templates.ExecuteTemplate(response, "outer", nil)

		for i := 2; i <= len(ticket.GetOpenTickets())+1; i++ {
			//tmp2 := templates.Lookup("ticketsViewTicketList.html")
			//tmp2.ExecuteTemplate(response, "inner", ticket.ReadTicket(i))
			templates.ExecuteTemplate(response, "inner", ticket.ReadTicket(i))
		}

		templates.Execute(response, nil)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}
