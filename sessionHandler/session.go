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
	Salt     string `json:"salt"`
}

// Read cookie and extract the value as username
func GetSessionUser(request *http.Request) (username string) {
	if cookie, err := request.Cookie("sessionUser"); err == nil {
		username = cookie.Value
		return
	} else {
		username = ""
		return
	}
}

func IsUserLoggedIn(request *http.Request) bool {
	if GetSessionUser(request) != "" {
		return true
	}

	return false
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
	inputUsername := request.FormValue("username")
	inputPassword := request.FormValue("password")
	redirectTarget := "/"

	// Only authenticate user if input has been submitted
	if inputUsername != "" && inputPassword != "" {
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
			if users.Users[i].Username == inputUsername {
				if GetHash(inputPassword, users.Users[i].Salt) == users.Users[i].Password {
					setSession(inputUsername, response)
					redirectTarget = "/internal"
				}
			}
		}
	}

	http.Redirect(response, request, redirectTarget, 302)
}

// Stop active session and redirect user to front html
func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}
