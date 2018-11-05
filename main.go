package main

import (
	"fmt"
	"io"
	"net/http"
	"ticketBackend/sessionHandler"
)

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := sessionHandler.GetUserName(request)

	if userName != "" {
		http.Redirect(response, request, "/internal", 302)
	} else {
		fmt.Fprintf(response, indexPage)
	}
}

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := sessionHandler.GetUserName(request)

	if userName != "" {
		fmt.Fprintf(response, internalPage, userName)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	server := http.Server{
		Addr:    ":8000",
		Handler: &myHandler{},
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = indexPageHandler
	mux["/internal"] = internalPageHandler
	mux["/login"] = sessionHandler.LoginHandler
	mux["/logout"] = sessionHandler.LogoutHandler

	server.ListenAndServe()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	io.WriteString(w, "Not found: "+r.URL.String())
}
