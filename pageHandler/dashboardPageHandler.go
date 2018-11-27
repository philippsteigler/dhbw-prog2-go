package pageHandler

import (
	"fmt"
	"html/template"
	"net/http"
	"ticketBackend/sessionHandler"
	"ticketBackend/ticket"
)

// A-8.1:
// Die Bearbeitung der Tickets soll ausschließlich ¨uber eine WEB-Seite erfolgen.
//
// Aufruf der Dashboardseite
// Der Nutzer sieht seine Tickets und kann die Liste der offenen einsehen
func DashboardViewPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {

		var templateFiles []string
		templateFiles = append(templateFiles, "./assets/html/dashboardTemplates/dashboardViewHeaderCssTemplate.html")
		templateFiles = append(templateFiles, "./assets/html/dashboardTemplates/dashboardTicketListTemplate.html")
		templateFiles = append(templateFiles, "./assets/html/dashboardTemplates/dashboardViewFooterTemplate.html")

		templates, err := template.ParseFiles(templateFiles...)
		if err != nil {
			fmt.Println(err)
		}

		templates.ExecuteTemplate(response, "outer", sessionHandler.GetSessionUser(request).Username)

		pTickets := *ticket.GetTickets(ticket.Open)

		templates.ExecuteTemplate(response, "outer", sessionHandler.GetSessionUser(request).Username)
		templates.ExecuteTemplate(response, "inner", pTickets)
		templates.ExecuteTemplate(response, "footer", nil)

		templates.Execute(response, nil)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}
