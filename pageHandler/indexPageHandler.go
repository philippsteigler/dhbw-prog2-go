package pageHandler

import (
	"../sessionHandler"
	"net/http"
)

//anzeigen der Index Seite
func IndexPageHandler(response http.ResponseWriter, request *http.Request) {
	if sessionHandler.IsUserLoggedIn(request) {
		http.Redirect(response, request, "/dashboard", 302) //Umbennenen in / ... die Benutzerseite
	} else {
		http.ServeFile(response, request, "./assets/html/index.html")
	}
}
