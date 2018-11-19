package main

import (
	"./pageHandler"
	"./sessionHandler"
	"log"
	"net/http"
	"os"
	"strconv"
)

func errorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

// TODO: Umwandeln in init() --> Funktion soll flags empfangen und verarbeiten
func createDirIfNotExist(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, 0755)
		errorCheck(err)
	}
}

// A-3.1:
// Die Web-Seite soll nur per HTTPS erreichbar sein.
//
// Die Startfunktion initialisiert die Anwendung und startet anschließend den Server.
// Eingehende HTTP-Anfragen auf Webseiten werden hier geroutet.
//
//
// A-11.3:
// Die Anwendung soll zwar HTTPS und die entsprechenden erforderlichen Zertifikate unterstützen,
// es kann jedoch davon ausgegangen werden, dass geeignete Zertifikate gestellt werden.
//
// Self-signed Zertifikate sind default vorhanden und unter ./assets/certificates gespeichert.
func main() {
	createDirIfNotExist("./assets/tickets")
	port := 8000
	portString := strconv.Itoa(port)
	mux := http.NewServeMux()

	// Webpages
	mux.Handle("/", http.HandlerFunc(pageHandler.IndexPageHandler))
	mux.Handle("/ticketsView", http.HandlerFunc(pageHandler.TicketsViewPageHandler))
	mux.Handle("/ticketInsightView", http.HandlerFunc(pageHandler.TicketInsightPageHandler))
	mux.Handle("/newTicketView", http.HandlerFunc(pageHandler.NewTicketViewPageHandler))
	mux.Handle("/dashboard", http.HandlerFunc(pageHandler.DashboardViewPageHandler))

	// Interactions
	mux.Handle("/login", http.HandlerFunc(sessionHandler.LoginHandler))
	mux.Handle("/logout", http.HandlerFunc(sessionHandler.LogoutHandler))
	mux.Handle("/ticketSafe", http.HandlerFunc(pageHandler.TicketSafeHandler))
	mux.Handle("/ticketTake", http.HandlerFunc(pageHandler.TicketTakeHandler))
	mux.Handle("/ticketSubmit", http.HandlerFunc(pageHandler.TicketSubmitHandler))
	mux.Handle("/ticketDelegate", http.HandlerFunc(pageHandler.TicketDelegateHandler))
	mux.Handle("/ticketNewEntry", http.HandlerFunc(pageHandler.TicketNewEntryHandler))

	log.Print("Listening on port " + portString + " ... ")
	err := http.ListenAndServeTLS(
		":"+portString,
		"./assets/certificates/server.crt",
		"./assets/certificates/server.key",
		mux,
	)

	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
