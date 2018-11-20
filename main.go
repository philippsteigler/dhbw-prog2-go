package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"ticketBackend/pageHandler"
	"ticketBackend/sessionHandler"
)

// Erstelle eine Verzeichnis für Tickets, sofern dieses nicht existiert.
func checkEnvironment() {
	if _, err := os.Stat(sessionHandler.GetAssetsDir() + "tickets"); os.IsNotExist(err) {
		err = os.Mkdir(sessionHandler.GetAssetsDir()+"tickets", 0755)
		sessionHandler.HandleError(err)
	}
}

// Kopiere eine Datei an eine andere Stelle.
func copyFile(src string, dst string) {
	data, err := ioutil.ReadFile(src)
	sessionHandler.HandleError(err)

	err = ioutil.WriteFile(dst, data, 0644)
	sessionHandler.HandleError(err)
}

//Setz den Webserver zurück, indem alle Tickets und Nutzerdaten gelöscht werden.
func resetData() {
	err := os.RemoveAll(sessionHandler.GetAssetsDir() + "tickets")
	sessionHandler.HandleError(err)

	err = os.RemoveAll(sessionHandler.GetAssetsDir() + "users.json")
	sessionHandler.HandleError(err)
}

// Setze den Server zurück und installiere Testdaten.
func demoMode() {
	resetData()

	// Kopiere alle Tickets aus dem Demo-Ordner in den Zielordner.
	src := strings.Join([]string{sessionHandler.GetAssetsDir(), "demo/tickets"}, "")
	files, err := ioutil.ReadDir(src)
	sessionHandler.HandleError(err)

	for _, file := range files {
		srcFile := strings.Join([]string{sessionHandler.GetAssetsDir(), "demo/tickets/", file.Name()}, "")
		dstFile := strings.Join([]string{sessionHandler.GetAssetsDir(), "tickets/", file.Name()}, "")
		copyFile(srcFile, dstFile)
	}

	// Kopiere die Nutzerdaten aus dem Demo-Ordner in den Zielordner.
	srcFile := strings.Join([]string{sessionHandler.GetAssetsDir(), "demo/users.json"}, "")
	copyFile(srcFile, sessionHandler.GetAssetsDir())
}

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
//  -reset=x    bool    True: Löscht alle Tickets und Nutzerdaten.
//  -demo=x     bool    True: Setzt den Webserver zurück und installiert Testdaten
// Das Flag -reset überschreibt dabei das Flag -demo.
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
	reset := flag.Bool("reset", false, "Delete all ticket and user data.")
	demo := flag.Bool("demo", false, "Install example data for tickets and users.")
	flag.Parse()

	if *demo == true {
		demoMode()
	} else if *reset == true {
		resetData()
	}

	checkEnvironment()
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
