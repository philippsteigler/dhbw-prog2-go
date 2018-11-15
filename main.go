package main

import (
	"./pageHandler"
	"./sessionHandler"
	"fmt"
	"io"
	"net/http"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	fmt.Println("[Server]: STARTING...")
	server := http.Server{
		Addr:    ":8000",
		Handler: &myHandler{},
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = pageHandler.IndexPageHandler
	mux["/internal"] = pageHandler.InternalPageHandler
	mux["/login"] = sessionHandler.LoginHandler
	mux["/logout"] = sessionHandler.LogoutHandler
	mux["/ticket"] = pageHandler.TicketPageHandler
	mux["/saveTicket"] = pageHandler.SaveTicketHandler

	fmt.Println("[Server]: Listening on https://localhost:8000/")
	err := server.ListenAndServeTLS("./assets/certificates/server.crt", "./assets/certificates/server.key")
	if err != nil {
		fmt.Print(err)
	}
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	io.WriteString(w, "Not found: "+r.URL.String())
}
