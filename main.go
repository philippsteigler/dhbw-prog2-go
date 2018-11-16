package main

import (
	"./pageHandler"
	"./sessionHandler"
	"./ticket"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Status string

const (
	Open      Status = "offen"
	InProcess Status = "in Bearbeitung"
	Closed    Status = "geschlossen"
)

type Ticket struct {
	Id       int     `json:"id"`
	Subject  string  `json:"subject"`
	Status   Status  `json:"status"`
	EditorId int     `json:"editorId"`
	Entries  []Entry `json:"entries"`
}

type Entry struct {
	Date    string `json:"date"`
	Creator string `json:"creator"`
	Content string `json:"content"`
}

func errorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		http.Redirect(response, request, "/internal", 302)
	} else {
		http.ServeFile(response, request, "./assets/html/index.html")
	}
}

func internalTickets(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {

		tEntry := Entry{Date: time.Now().Local().Format("2006-01-02T15:04:05Z07:00"),
			Creator: "jamaijarno",
			Content: "auf die druffe du sack"}

		tTicket := Ticket{Id: 1,
			Subject:  "Test",
			Status:   Open,
			EditorId: 420,
			Entries:  []Entry{tEntry}}

		internal, err := template.ParseFiles("./pageHandler/internalTemplate.tmpl")
		if err != nil {
			fmt.Println(err)
		}

		internal.ExecuteTemplate(response, "internal", tTicket)

		internal.Execute(response, tTicket)

	} else {
		http.ServeFile(response, request, "./assets/html/index.html")
	}
}

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {

		var templateFiles []string
		templateFiles = append(templateFiles, "./pageHandler/internalTicketsTemplate.tmpl")
		templateFiles = append(templateFiles, "./pageHandler/internalTicketsListTemplate.tmpl")

		templates, err := template.ParseFiles(templateFiles...)
		if err != nil {
			fmt.Println(err)
		}

		templates.ExecuteTemplate(response, "outer", nil)

		for i := 2; i <= len(ticket.GetOpenTickets())+1; i++ {
			//tmp2 := templates.Lookup("internalTicketsListTemplate.tmpl")
			//tmp2.ExecuteTemplate(response, "inner", ticket.ReadTicket(i))
			templates.ExecuteTemplate(response, "inner", ticket.ReadTicket(i))
		}

		templates.Execute(response, nil)

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

func createDirIfNotExist(folder string) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, 0755)
		errorCheck(err)
	}
}

// Golang webserver example:
// https://github.com/jimmahoney/golang-webserver/blob/master/webserver.go
func main() {
	createDirIfNotExist("./assets/tickets")
	port := 8000
	portstring := strconv.Itoa(port)

	mux := http.NewServeMux()
	// Webpages
	mux.Handle("/", http.HandlerFunc(indexPageHandler))
	mux.Handle("/internal", http.HandlerFunc(internalPageHandler))
	mux.Handle("/internal/ticket/information", http.HandlerFunc(internalTickets))
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
