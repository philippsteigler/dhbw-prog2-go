package pageHandler

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"de/vorlesung/projekt/crew/ticket"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// A-8.3
// Bearbeiter sollen alle Tickets einsehen können, welche noch kein Bearbeiter
// übernommen hat.
//
// https://localhost:8000/ticketInsightView
// Detailansicht des Tickets
func TicketInsightPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		id := request.FormValue("TicketID")

		intId, err := strconv.Atoi(id)

		if err != nil {
			fmt.Println(err)
		}

		var templateFiles []string
		templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/ticketInsightTemplates/ticketInsightViewHeaderCssTemplate.html")
		templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/ticketInsightTemplates/ticketInsightTicketDetailsTemplate.html")

		templates, err := template.ParseFiles(templateFiles...)
		if err != nil {
			fmt.Println(err)
		}

		templates.ExecuteTemplate(response, "outer", sessionHandler.GetSessionUser(request).Username)
		templates.ExecuteTemplate(response, "inner", ticket.GetTicket(intId))
		templates.ExecuteTemplate(response, "footer", nil)

		templates.Execute(response, nil)
	} else {
		http.ServeFile(response, request, sessionHandler.GetAssetsDir()+"html/loginView.html")
	}
}

// A-8.2
// Bearbeiter sollen ein Ticket übernehmen können.
//
// https://localhost:8000/ticketTake
// Ticket übernehmen, Web Interaction
//TODO: EditorID mit geben!!!!, Redirect überarbeiten
// localhost:.../ticketTake
//Funktion ticket nehmnen
func TicketTakeHandler(response http.ResponseWriter, request *http.Request) {

	if sessionHandler.IsUserLoggedIn(request) {
		idToParse := request.FormValue("TicketID")
		ticketId, err := strconv.Atoi(idToParse)
		sessionHandler.HandleError(err)

		ticket.TakeTicket(ticketId, sessionHandler.GetSessionUser(request).ID)
		// Zurück zu der Ticketseite
		http.Redirect(response, request, "/dashboard", http.StatusFound)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}

//TODO: Redirect überarbeiten
// localhost:.../ticketSubmit
//Funktion ticket abgeben

// A-8.4
// Bearbeiter sollen Tickets nach der Übernahme auch freigeben können, so das
// diese eine anderer Bearbeiter übernehmen kann.
//
// https://localhost:8000/ticketSubmit
// Ticket abgeben, Web Interaction
func TicketSubmitHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		idToParse := request.FormValue("TicketID")
		ticketId, err := strconv.Atoi(idToParse)
		sessionHandler.HandleError(err)

		ticket.UnhandTicket(ticketId)
		// Zurück zu der Ticketseite
		http.Redirect(response, request, "/dashboard", http.StatusFound)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// A-8.5
// Ein Bearbeiter soll ein Ticket einem anderen Bearbeiter zuteilen können
//
// https://localhost:8000/ticketDelegate
// Ticket delegieren, Web Interaction
//TODO: EditorId mitgben; redirect überarbeiten
// localhost:.../ticketDelegate
//Funktion ticket Delegieren
func TicketDelegateHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		idToParse := request.FormValue("TicketID")
		ticketId, err := strconv.Atoi(idToParse)
		sessionHandler.HandleError(err)

		ticket.TakeTicket(ticketId, 0)
		// Zurück zu der Ticketseite
		http.Redirect(response, request, "/dashboard", http.StatusFound)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// A-8.6
// Kommentiert ein Bearbeiter ein Ticket, soll er wählen können, ob dieser Kommentar
// an den Kunden versendet wird, oder ob er nur für andere Bearbeiter
// sichtbar ist.
//
// Ticket Eintrag hinzufügen, Web Interaction
func TicketAppendEntry(w http.ResponseWriter, r *http.Request) {

	if sessionHandler.IsUserLoggedIn(r) {
		id := r.FormValue("TicketID")
		intId, err := strconv.Atoi(id)
		sessionHandler.HandleError(err)

		ticket.AppendEntry(intId, strconv.Itoa(sessionHandler.GetSessionUser(r).ID), r.FormValue("entryText"), false)
		// Zurück zu der Ticketseite
		http.Redirect(w, r, "/dashboard", http.StatusFound)

	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func TicketShowHistory(response http.ResponseWriter, request *http.Request) {

}
