package pageHandler

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"de/vorlesung/projekt/crew/ticket"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Matrikelnummern:
//
// 3333958
// 3880065
// 8701350

var ticketInsightTemplates *template.Template
var ticketID int

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
		ticketID = intId

		ticketInsightTemplates.ExecuteTemplate(response, "outer", sessionHandler.GetSessionUser(request).Username)

		if ticket.GetTicket(intId).Status == "offen" {
			ticketInsightTemplates.ExecuteTemplate(response, "open", ticket.GetTicket(intId))
		} else if ticket.GetTicket(intId).Status == "in Bearbeitung" {

			ticketInsightTemplates.ExecuteTemplate(response, "select", sessionHandler.GetAllOtherUserIDs(request))
			ticketInsightTemplates.ExecuteTemplate(response, "taken", ticket.GetTicket(intId))

		}

		ticketInsightTemplates.Execute(response, nil)
	} else {
		http.ServeFile(response, request, sessionHandler.GetAssetsDir()+"html/loginView.html")
	}
}

func TicketInsightInit() {
	var templateFiles []string
	var err error

	templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/ticketInsightTemplates/ticketInsightViewHeaderCssTemplate.html")
	templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/ticketInsightTemplates/ticketInsightTicketDetailsOpenTemplate.html")
	templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/ticketInsightTemplates/ticketInsightTicketDetailsTakenTemplate.html")
	templateFiles = append(templateFiles, sessionHandler.GetAssetsDir()+"html/ticketInsightTemplates/ticketInsightTicketDetailsUserSelectTemplate.html")

	ticketInsightTemplates, err = template.ParseFiles(templateFiles...)
	if err != nil {
		fmt.Println(err)
	}
}

// A-8.2
// Bearbeiter sollen ein Ticket übernehmen können.
//
// https://localhost:8000/ticketTake
// Ticket übernehmen, Web Interaction
// localhost:.../ticketTake
// Funktion ticket nehmnen
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

// localhost:.../ticketSubmit
// Funktion ticket abgeben

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
// localhost:.../ticketDelegate
// Funktion ticket Delegieren
func TicketDelegateHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		eID := request.FormValue("select")
		editorID, err := strconv.Atoi(eID)
		sessionHandler.HandleError(err)

		ticket.TakeTicket(ticketID, editorID)
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

		sendingMail := false
		if r.FormValue("mail") == "mail" {
			sendingMail = true
		}

		ticket.AppendEntry(intId, sessionHandler.GetUsername(sessionHandler.GetSessionUser(r).ID), r.FormValue("entryText"), sendingMail)
		// Zurück zu der Ticketseite
		http.Redirect(w, r, "/dashboard", http.StatusFound)

	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func TicketClose(w http.ResponseWriter, r *http.Request) {
	if sessionHandler.IsUserLoggedIn(r) {
		id := r.FormValue("TicketID")
		intId, err := strconv.Atoi(id)
		sessionHandler.HandleError(err)

		ticket.CloseTicket(intId)
		// Zurück zu der Ticketseite
		http.Redirect(w, r, "/dashboard", http.StatusFound)

	} else {
		http.Redirect(w, r, "/", 302)
	}
}
