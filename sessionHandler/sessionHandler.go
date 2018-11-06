package sessionHandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UserAccounts struct {
	Users []User `json:"users"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Read cookie and extract the value (user name)
func GetSessionUser(request *http.Request) (username string) {
	if cookie, err := request.Cookie("sessionUser"); err == nil {
		username = cookie.Value
	}

	return username
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
	sessionUsername := request.FormValue("username")
	sessionPassword := request.FormValue("password")
	redirectTarget := "/"

	if sessionUsername != "" && sessionPassword != "" {
		userData, err := ioutil.ReadFile("./assets/users.json")
		var users UserAccounts

		err = json.Unmarshal(userData, &users)
		if err != nil {
			fmt.Println(err)
		}

		for i := 0; i < len(users.Users); i++ {
			if sessionUsername == users.Users[i].Username && sessionPassword == users.Users[i].Password {
				setSession(sessionUsername, response)
				redirectTarget = "/internal"
			}
		}
	}

	http.Redirect(response, request, redirectTarget, 302)
}

// Stop active session and redirect the user
func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}
