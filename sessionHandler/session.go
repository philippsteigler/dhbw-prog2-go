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
	Username string `json:"username"`
	Password string `json:"password"`
}

// Read cookie and extract the value as username
func GetSessionUser(request *http.Request) (username string) {
	if cookie, err := request.Cookie("sessionUser"); err == nil {
		username = cookie.Value
	}

	return username
}

// Deploy cookie to save active user session
func setSession(username string, response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "sessionUser",
		Value: username,
		Path:  "/",
	}

	http.SetCookie(response, cookie)
}

// Delete cookie to end active session
func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "sessionUser",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(response, cookie)
}

// Process user input from login action and start new session
func LoginHandler(response http.ResponseWriter, request *http.Request) {
	sessionUsername := request.FormValue("username")
	sessionPassword := request.FormValue("password")
	redirectTarget := "/"

	// Only authenticate user if input has been submitted
	if sessionUsername != "" && sessionPassword != "" {
		var users UserAccounts

		// Read data for registered users
		userData, err := ioutil.ReadFile("./assets/users.json")
		if err != nil {
			fmt.Print(err)
		}

		err = json.Unmarshal(userData, &users)
		if err != nil {
			fmt.Print(err)
		}

		// Search claimed user and evaluate password from user input
		for i := 0; i < len(users.Users); i++ {
			tmp, err := decryptString(users.Users[i].Password, getKey())
			if err != nil {
				fmt.Print(err)
			}

			if sessionUsername == users.Users[i].Username && sessionPassword == tmp {
				setSession(sessionUsername, response)
				redirectTarget = "/internal"
			}
		}
	}

	http.Redirect(response, request, redirectTarget, 302)
}

// Stop active session and redirect user to front page
func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}
