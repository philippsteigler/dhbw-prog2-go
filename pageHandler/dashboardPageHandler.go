package pageHandler

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"de/vorlesung/projekt/crew/ticket"
	"fmt"
	"html/template"
	"net/http"
)

var dashboardViewTemplates *template.Template

// A-8.1:
// Die Bearbeitung der Tickets soll ausschließlich ¨uber eine WEB-Seite erfolgen.
//
// Aufruf der Dashboardseite
// Der Nutzer sieht seine Tickets und kann die Liste der offenen einsehen
func DashboardViewPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {

		pTickets := *ticket.GetTicketsByEditorId(sessionHandler.GetSessionUser(request).ID)

		dashboardViewTemplates.ExecuteTemplate(response, "outer", sessionHandler.GetSessionUser(request).Username)
		dashboardViewTemplates.ExecuteTemplate(response, "inner", pTickets)
		dashboardViewTemplates.ExecuteTemplate(response, "footer", nil)

		dashboardViewTemplates.Execute(response, nil)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func DashboardViewInit() {
	var templateFiles []string
	var err error

	templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/dashboardTemplates/dashboardViewHeaderCssTemplate.html")
	templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/dashboardTemplates/dashboardTicketListTemplate.html")

	dashboardViewTemplates, err = template.ParseFiles(templateFiles...)
	if err != nil {
		fmt.Println(err)
	}
}
