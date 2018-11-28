package main

import (
	"de/vorlesung/projekt/crew/pageHandler"
	"de/vorlesung/projekt/crew/sessionHandler"
	"flag"
	"log"
	"net/http"
	"strconv"
)

// A-3.1:
// Die Web-Seite soll nur per HTTPS erreichbar sein.
//
// Die Startfunktion initialisiert die Anwendung und startet anschließend den Server.
// Eingehende HTTP-Anfragen auf Webseiten werden hier geroutet.
//
//
// A-10.1:
// Die Konfiguration soll komplett über Startparameter erfolgen.
//
// Der Anwender kann beim Starten über die Kommandozeile folgende Flags optional setzen:
//  -port=x     int     Port für den Webserver
//  -default=x  bool    True: Löscht alle Tickets und Nutzerdaten.
//  -demo=x     bool    True: Setzt den Webserver zurück und installiert Testdaten
// Das Flag -default überschreibt dabei das Flag -demo.
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
// Self-signed Zertifikate sind demo vorhanden und unter ./assets/certificates gespeichert.
func main() {
	port := flag.Int("port", 8000, "Port for webserver.")
	reset := flag.Bool("reset", false, "Delete all ticket and user rollback.")
	demo := flag.Bool("demo", false, "Install example rollback for tickets and users.")
	flag.Parse()

	sessionHandler.CheckEnvironment()

	if *demo == true {
		sessionHandler.DemoMode()
	} else if *reset == true {
		sessionHandler.ResetData()
	}

	mux := http.NewServeMux()

	// Webpages
	mux.Handle("/", http.HandlerFunc(pageHandler.NewTicketViewPageHandler))
	mux.Handle("/loginView", http.HandlerFunc(pageHandler.LoginPageHandler))
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
		sessionHandler.GetAssetsDir()+"certificates/cert.pem",
		sessionHandler.GetAssetsDir()+"certificates/key.pem",
		mux,
	)

	sessionHandler.HandleError(err)
}
