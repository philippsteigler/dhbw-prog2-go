package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"ticketBackend/pageHandler"
	"ticketBackend/sessionHandler"
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

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		http.Redirect(response, request, "/internal", 302)
	} else {
		http.ServeFile(response, request, "./assets/html/index.html")
	}
}

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		http.ServeFile(response, request, "./assets/html/internal.html")
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func ticketPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		http.ServeFile(response, request, "./assets/html/ticket.html")
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// A-3.1:
// Die Web-Seite soll nur per HTTPS erreichbar sein.
//
// Die Startfunktion initialisiert die Anwendung und startet anschließend den Server.
// Eingehende HTTP-Anfragen auf Webseiten werden hier geroutet.
// https://github.com/jimmahoney/golang-webserver/blob/master/webserver.go
func main() {
	createDirIfNotExist("./assets/tickets")
	port := 8000
	portString := strconv.Itoa(port)
	mux := http.NewServeMux()

	// Webseiten routen
	mux.Handle("/", http.HandlerFunc(indexPageHandler))
	mux.Handle("/internal", http.HandlerFunc(internalPageHandler))
	mux.Handle("/ticket", http.HandlerFunc(ticketPageHandler))

	// Interaktionen verarbeiten (z.B. Buttons, ...)
	mux.Handle("/login", http.HandlerFunc(sessionHandler.LoginHandler))
	mux.Handle("/logout", http.HandlerFunc(sessionHandler.LogoutHandler))
	mux.Handle("/saveTicket", http.HandlerFunc(pageHandler.SaveTicketHandler))

	log.Print("Listening on port " + portString + " ... ")
	err := http.ListenAndServeTLS(":"+portString, "./assets/certificates/server.crt", "./assets/certificates/server.key", mux)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
