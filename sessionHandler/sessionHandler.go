package sessionHandler

import "net/http"

// Read the cookie and extract the value (user name)
func GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("sessionUser"); err == nil {
		userName = cookie.Value
	}

	return userName
}

// Deploy a cookie to save the active user session
func setSession(userName string, response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "sessionUser",
		Value: userName,
		Path:  "/",
	}

	http.SetCookie(response, cookie)
}

// Delete the cookie to end an active session
func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "sessionUser",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(response, cookie)
}

// Process user input from login action and start a new session
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

// Stop the active session and redirect the user
func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}
