package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"ticketBackend/pageHandler"
	"ticketBackend/sessionHandler"
)

func resetData() {
	err := syscall.Rmdir(sessionHandler.GetAssetsDir() + "tickets")
	sessionHandler.HandleError(err)
}

// TODO: Umwandeln in init() --> Funktion soll flags empfangen und verarbeiten
func createDirIfNotExist(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, 0755)
		sessionHandler.HandleError(err)
	}
}

// A-3.1:
// Die Web-Seite soll nur per HTTPS erreichbar sein.
//
// Die Startfunktion initialisiert die Anwendung und startet anschließend den Server.
// Eingehende HTTP-Anfragen auf Webseiten werden hier geroutet.
//
//
// A-10.2:
// Der Port muss sich über ein Flag festlegen lassen.
//
// Der Anwender kann beim Starten über die Kommandozeile mit -port=X den Port des Servers bestimmen.
// Der Default-Wert ist Port :8000.
//
//
// A-11.3:
// Die Anwendung soll zwar HTTPS und die entsprechenden erforderlichen Zertifikate unterstützen,
// es kann jedoch davon ausgegangen werden, dass geeignete Zertifikate gestellt werden.
//
// Self-signed Zertifikate sind default vorhanden und unter ./assets/certificates gespeichert.
func main() {
	port := flag.Int("port", 8000, "Port for webserver.")
	reset := flag.Bool("reset", false, "Delete all ticket and user data.")
	//demo := flag.Bool("demo", false, "Install example data for tickets and users.")
	flag.Parse()

	if *reset == true {
		resetData()
	}

	createDirIfNotExist("./assets/tickets")
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

	log.Print("Listening on port " + strconv.Itoa(*port) + " ... ")
	err := http.ListenAndServeTLS(
		":"+strconv.Itoa(*port),
		sessionHandler.GetAssetsDir()+"certificates/server.crt",
		sessionHandler.GetAssetsDir()+"certificates/server.key",
		mux,
	)

	sessionHandler.HandleError(err)
}
