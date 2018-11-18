package pageHandler

import (
	"../sessionHandler"
	"../ticket"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
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

// localhost:.../ticketTake
//Funktion ticket nehmnen
func TicketTakeHandler(response http.ResponseWriter, request *http.Request) {

}

// localhost:.../ticketSubmit
//Funktion ticket abgeben
func TicketSubmitHandler(response http.ResponseWriter, request *http.Request) {

}

// localhost:.../ticketDelegate
//Funktion ticket Delegieren
func TicketDelegateHandler(response http.ResponseWriter, request *http.Request) {

}

// localhost:.../ticketNewEntry
//Funktion ticket Eintrag hinzufügen
func TicketNewEntryHandler(response http.ResponseWriter, request *http.Request) {

}
