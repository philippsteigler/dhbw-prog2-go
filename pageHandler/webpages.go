package pageHandler

import (
	"../sessionHandler"
	"fmt"
	"io/ioutil"
	"net/http"
)

func IndexPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
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
	if sessionHandler.IsUserLoggedIn(request) {
		file, err := ioutil.ReadFile("./assets/html/internal.html")
		if err != nil {
			fmt.Print(err)
		}
		fmt.Fprintf(response, string(file), sessionHandler.GetSessionUser(request))
	} else {
		http.Redirect(response, request, "/", 302)
	}
}
