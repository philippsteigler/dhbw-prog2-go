package pageHandler

import (
	"../sessionHandler"
	"net/http"
)

func IndexPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		http.Redirect(response, request, "/ticket", 302) //Umbennenen in / ... die Benutzerseite
	} else {
		http.ServeFile(response, request, "./assets/html/index.html")
	}
}
