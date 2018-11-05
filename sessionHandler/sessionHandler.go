package sessionHandler

import "net/http"

func GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("sessionHandler"); err == nil {
		userName = cookie.Value
	}

	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "sessionHandler",
		Value: userName,
		Path:  "/",
	}

	http.SetCookie(response, cookie)
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "sessionHandler",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(response, cookie)
}

func LoginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	redirectTarget := "/"

	if name != "" && pass != "" {
		setSession(name, response)
		redirectTarget = "/internal"
	}

	http.Redirect(response, request, redirectTarget, 302)
}

func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}
