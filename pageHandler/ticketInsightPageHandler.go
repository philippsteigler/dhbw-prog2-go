package pageHandler

import (
	"../sessionHandler"
	"../ticket"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

//TODO: Nur

// A-8.3
// Bearbeiter sollen alle Tickets einsehen k¨onnen, welche noch kein Bearbeiter
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

// A-8.2
// Bearbeiter sollen ein Ticket übernehmen können.
//
// https://localhost:8000/ticketTake
// Ticket übernehmen, Web Interaction
func TicketTakeHandler(response http.ResponseWriter, request *http.Request) {

}

// localhost:.../ticketSubmit
//Funktion ticket abgeben

// A-8.4
// Bearbeiter sollen Tickets nach der Übernahme auch freigeben können, so das
// diese eine anderer Bearbeiter übernehmen kann.
//
// https://localhost:8000/ticketSubmit
// Ticket abgeben, Web Interaction
func TicketSubmitHandler(response http.ResponseWriter, request *http.Request) {

}

// A-8.5
// Ein Bearbeiter soll ein Ticket einem anderen Bearbeiter zuteilen können
//
// https://localhost:8000/ticketDelegate
// Ticket delegieren, Web Interaction
func TicketDelegateHandler(response http.ResponseWriter, request *http.Request) {

}

// localhost:.../ticketNewEntry
//Funktion ticket Eintrag hinzufügen
func TicketNewEntryHandler(response http.ResponseWriter, request *http.Request) {

}
