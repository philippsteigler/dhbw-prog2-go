package main

import (
	"./sessionHandler"
	"log"
	"net/http"
	"strconv"
	"ticketBackend/pageHandler"
	sessionHandler2 "ticketBackend/sessionHandler"
)

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler2.IsUserLoggedIn(request) {
		http.Redirect(response, request, "/internal", 302)
	} else {
		http.ServeFile(response, request, "./assets/html/index.html")
	}
}

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler2.IsUserLoggedIn(request) {
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

// Golang webserver example:
// https://github.com/jimmahoney/golang-webserver/blob/master/webserver.go
func main() {
	port := 8000
	portstring := strconv.Itoa(port)

	mux := http.NewServeMux()
	// Webpages
	mux.Handle("/", http.HandlerFunc(indexPageHandler))
	mux.Handle("/internal", http.HandlerFunc(internalPageHandler))
	mux.Handle("/ticket", http.HandlerFunc(ticketPageHandler))

	// Interactions
	mux.Handle("/login", http.HandlerFunc(sessionHandler.LoginHandler))
	mux.Handle("/logout", http.HandlerFunc(sessionHandler.LogoutHandler))
	mux.Handle("/saveTicket", http.HandlerFunc(pageHandler.SaveTicketHandler))

	log.Print("Listening on port " + portstring + " ... ")
	err := http.ListenAndServeTLS(":"+portstring, "./assets/certificates/server.crt", "./assets/certificates/server.key", mux)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
