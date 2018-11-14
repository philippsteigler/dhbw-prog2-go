package pageHandler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"ticketBackend/sessionHandler"
)

func IndexPageHandler(response http.ResponseWriter, request *http.Request) {
	username := sessionHandler.GetSessionUser(request)

	if username != "" {
		http.Redirect(response, request, "/internal", 302)
	} else {
		file, err := ioutil.ReadFile("./assets/html/index.html")
		if err != nil {
			fmt.Print(err)
		}

		fmt.Fprintf(response, string(file))
	}
}

func InternalPageHandler(response http.ResponseWriter, request *http.Request) {
	username := sessionHandler.GetSessionUser(request)

	if username != "" {
		file, err := ioutil.ReadFile("./assets/html/internal.html")
		if err != nil {
			fmt.Print(err)
		}
		fmt.Fprintf(response, string(file), username)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}
