package ticketPage

import (
	"../sessionHandler"
	"fmt"
	"io/ioutil"
	"net/http"
)

//localhost:.../ticket
//läd den html code für die Textareas und zeigt ihn an
func TicketPageHandler(response http.ResponseWriter, request *http.Request) {
	username := sessionHandler.GetSessionUser(request)

	if username != "" {
		//Einlesen der .html Datei
		file, err := ioutil.ReadFile("./assets/page/ticket.html")
		if err != nil {
			fmt.Print(err)
		}

		//Anzeigen der .html
		//Format Print Formatiert dach dem responseWriter => html
		fmt.Fprintf(response, string(file))

	} else {
		http.Redirect(response, request, "/", 302)
	}
}

//localhost:.../saveTicket
//Speichert den Text aus den Textareas in mail, subject, text
func SaveTicketHandler(response http.ResponseWriter, request *http.Request) {
	username := sessionHandler.GetSessionUser(request)

	if username != "" {
		mail := request.FormValue("ticketMail")
		subject := request.FormValue("ticketSubject")
		text := request.FormValue("ticketText")

		//Testdatei für die Eingabe
		inputTest := string("Mail\n" + mail + "\n\nSubject\n" + subject + "\n\nText\n" + text)
		ioutil.WriteFile("Testdatei für das Ticket", []byte(inputTest), 0600)

		//Zurück zu der Ticketseite
		http.Redirect(response, request, "/ticket", http.StatusFound)

	} else {
		http.Redirect(response, request, "/", 302)
	}
}
