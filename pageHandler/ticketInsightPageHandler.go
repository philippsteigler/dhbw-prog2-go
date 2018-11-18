package pageHandler

import (
	"../sessionHandler"
	"../ticket"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

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

func TicketTakeHandler(response http.ResponseWriter, request *http.Request) {

}

func TicketSubmitHandler(response http.ResponseWriter, request *http.Request) {

}

func TicketDelegateHandler(response http.ResponseWriter, request *http.Request) {

}

func TicketNewEntryHandler(response http.ResponseWriter, request *http.Request) {

}
