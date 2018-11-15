package pageHandler

import (
	"../sessionHandler"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Handler interface {
	ServeHTTP(http.Response,*http.Request)
}

// localhost:.../ticket
// läd den html code für die Textareas und zeigt ihn an
func TicketPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {

		html, _ := ioutil.ReadFile("./assets/html/ticket.html")

		fmt.Fprint(response, string(html))

		//http.ServeFile(response,request,"./assets/html/ticket.html")
		//http.ServeFile(response,request,"./assets/html/style/ticketStyle.css")

	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// localhost:.../saveTicket
// Speichert den Text aus den Textareas in mail, subject, text
func SaveTicketHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		mail := request.FormValue("ticketMail")
		subject := request.FormValue("ticketSubject")
		text := request.FormValue("ticketText")

		// Testdatei für die Eingabe
		inputTest := string("Mail\n" + mail + "\n\nSubject\n" + subject + "\n\nText\n" + text)
		ioutil.WriteFile("./assets/TicketTest", []byte(inputTest), 0600)

		// Zurück zu der Ticketseite
		http.Redirect(response, request, "/ticket", http.StatusFound)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}
