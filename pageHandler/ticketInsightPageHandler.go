package pageHandler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"ticketBackend/sessionHandler"
	"ticketBackend/ticket"
)

// localhost:.../ticketInsight
//anzeigen der Detailansicht Seite
func TicketInsightPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {

		id := request.FormValue("TicketID")

		intId, err := strconv.Atoi(id)

		if err != nil {
			fmt.Println(err)
		}

		internal, err := template.ParseFiles("./assets/html/ticketInsightViewTemplate.html")
		if err != nil {
			fmt.Println(err)
		}

		internal.ExecuteTemplate(response, "internal", ticket.GetTicket(intId))

		internal.Execute(response, nil)

	} else {
		http.ServeFile(response, request, "./assets/html/index.html")
	}
}

//TODO: EditorID mit geben!!!!, Redirect überarbeiten
// localhost:.../ticketTake
//Funktion ticket nehmnen
func TicketTakeHandler(response http.ResponseWriter, request *http.Request) {

	if sessionHandler.IsUserLoggedIn(request) {
		idToParse := request.FormValue("TicketID")
		ticketId, err := strconv.Atoi(idToParse)
		sessionHandler.HandleError(err)

		ticket.TakeTicket(ticketId, 0)
		// Zurück zu der Ticketseite
		http.Redirect(response, request, "/dashbord", http.StatusFound)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}

//TODO: Redirect überarbeiten
// localhost:.../ticketSubmit
//Funktion ticket abgeben
func TicketSubmitHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		idToParse := request.FormValue("TicketID")
		ticketId, err := strconv.Atoi(idToParse)
		sessionHandler.HandleError(err)

		ticket.UnhandTicket(ticketId)
		// Zurück zu der Ticketseite
		http.Redirect(response, request, "/dashbord", http.StatusFound)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}

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
		http.Redirect(response, request, "/dashbord", http.StatusFound)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// localhost:.../ticketNewEntry
//Funktion ticket Eintrag hinzufügen
func TicketNewEntryHandler(response http.ResponseWriter, request *http.Request) {

}
