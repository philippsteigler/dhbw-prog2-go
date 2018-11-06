package sessionHandler

import "net/http"

type User struct {
	username string `json:"username"`
	password string `json:"password"`
}

// Read cookie and extract the value (user name)
func GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("sessionUser"); err == nil {
		userName = cookie.Value
	}

	return userName
}

// Deploy cookie to save the active user session
func setSession(username string, response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "sessionUser",
		Value: username,
		Path:  "/",
	}

	http.SetCookie(response, cookie)
}

// Delete cookie to end an active session
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
	// TODO: save credentials to .json file
	sessionUser := User{request.FormValue("username"), request.FormValue("password")}
	redirectTarget := "/"

	if sessionUser.username != "" && sessionUser.password != "" {
		setSession(sessionUser.username, response)
		redirectTarget = "/internal"
	}

	http.Redirect(response, request, redirectTarget, 302)
}

// Stop active session and redirect the user
func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}
