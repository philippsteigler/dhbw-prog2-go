package sessionHandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

var users UserAccounts

// Read data for registered users
func refreshUserData() {
	userData, err := ioutil.ReadFile("./assets/users.json")
	if err != nil {
		fmt.Print(err)
	}

	err = json.Unmarshal(userData, &users)
	if err != nil {
		fmt.Print(err)
	}
}

// Deploy cookie to save user ID and user name as session identifier
func setSession(id int, username string, response http.ResponseWriter) {
	cookieUserID := &http.Cookie{
		Name:  "sessionUserID",
		Value: strconv.Itoa(id),
		Path:  "/",
	}
	cookieUserName := &http.Cookie{
		Name:  "sessionUserName",
		Value: username,
		Path:  "/",
	}

	http.SetCookie(response, cookieUserID)
	http.SetCookie(response, cookieUserName)
}

// Delete cookie to end active session
func clearSession(response http.ResponseWriter) {
	cookieUserID := &http.Cookie{
		Name:   "sessionUserID",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	cookieUserName := &http.Cookie{
		Name:   "sessionUserName",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(response, cookieUserID)
	http.SetCookie(response, cookieUserName)
}

// Read cookie and extract the value as username
func GetSessionUserName(request *http.Request) string {
	if cookie, err := request.Cookie("sessionUserName"); err == nil {
		return cookie.Value
	}

	return ""
}

// Read cookie and extract the value as user ID
func GetSessionUserID(request *http.Request) int {
	if cookie, err := request.Cookie("sessionUserID"); err == nil {
		i, err := strconv.Atoi(cookie.Value)
		if err != nil {
			fmt.Print(err)
		}
		return i
	}
	return 0
}

func IsUserLoggedIn(request *http.Request) bool {
	if GetSessionUserID(request) != 0 && GetSessionUserName(request) != "" {
		return true
	}
	return false
}

// Process user input from login action and start new session
func LoginHandler(response http.ResponseWriter, request *http.Request) {
	inputUsername := request.FormValue("username")
	inputPassword := request.FormValue("password")
	redirectTarget := "/"

	// Only authenticate user if input has been submitted
	if inputUsername != "" && inputPassword != "" {
		refreshUserData()

		// Search claimed user and evaluate password from user input
		for i := 0; i < len(users.Users); i++ {
			if users.Users[i].Username == inputUsername {
				if GetHash(inputPassword, users.Users[i].Salt) == users.Users[i].Password {
					setSession(users.Users[i].ID, users.Users[i].Username, response)
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

// Generate user credentials and append user to users.json
func RegistrationHandler(response http.ResponseWriter, request *http.Request) {
	hash, salt := HashString(request.FormValue("password"))
	refreshUserData()

	users.Users = append(users.Users, User{ID: len(users.Users), Username: request.FormValue("username"), Password: hash, Salt: salt})

	usersJson, _ := json.Marshal(users)
	err := ioutil.WriteFile("./assets/users.json", usersJson, 0644)
	if err != nil {
		fmt.Print(err)
	}

	http.Redirect(response, request, "/", 302)
}
