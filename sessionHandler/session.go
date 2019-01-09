package sessionHandler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type UserAccounts struct {
	Users []User `json:"users"`
}

type User struct {
	ID       int    `json:"id"`
	Tier     int    `json:"tier"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

type IDList struct {
	IDs []int
}

var users UserAccounts
var sessionUsers = map[string]User{}

// Lies users.json und importiere alle Benutzerdaten nach &users.
func LoadUserData() *UserAccounts {
	userData, err := ioutil.ReadFile(GetAssetsDir() + "users.json")
	HandleError(err)

	err = json.Unmarshal(userData, &users)
	HandleError(err)

	return &users
}

func GetUsername(id int) string {
	for _, v := range users.Users {
		if v.ID == id {
			return v.Username
		}
	}

	log.Print("User with ID '" + strconv.Itoa(id) + "' does not exist.")
	return ""
}

// A-3.2:
// Der Zugang für die Bearbeiter soll durch Benutzernamen und Passwort geschützt sein.
//
// Gib den Benutzer zurück, der gerade eingeloggt ist.
func GetSessionUser(r *http.Request) *User {
	// Lies zunächst den Session-Cookie des Benutzers ein, um dessen Session-Token zu erhalten.
	if cookie, err := r.Cookie("session_token"); err == nil {
		sessionToken := cookie.Value

		// Anschließend wird überprüft, ob dieser Token in der sessionUsers-Map enthalten ist.
		// Existiert der Eintrag, so wird der zugehörige Benutzer zurückgegeben.
		user, ok := sessionUsers[sessionToken]
		if ok {
			return &user
		}
	}

	log.Print("No user for current session.")
	return nil
}

// Gib alle UserIDs von Editoren zurück, außer vom derzeit eingeloggten User.
func GetAllOtherUserIDs(r *http.Request) IDList {
	list := IDList{
		IDs: []int{},
	}

	for _, v := range users.Users {
		if v.ID != GetSessionUser(r).ID {
			list.IDs = append(list.IDs, v.ID)
		}
	}

	return list
}

// A-3.2:
// Der Zugang für die Bearbeiter soll durch Benutzernamen und Passwort geschützt sein.
//
// Überprüfe anhand des Session-Cookies, ob ein Benutzer eingeloggt ist.
// Benutzer eingeloggt = true; Benutzer ist nicht eingeloggt = false.
func IsUserLoggedIn(r *http.Request) bool {
	if cookie, err := r.Cookie("session_token"); err == nil {
		sessionToken := cookie.Value
		_, ok := sessionUsers[sessionToken]
		return ok
	}

	return false
}

// A-3.2:
// Der Zugang für die Bearbeiter soll durch Benutzernamen und Passwort geschützt sein.
//
// Authentifizierung des Benutzers.
// Die Eingaben des Nutzers werden mit gespeicherten Credentials abgeglichen.
// Folgende Testuser sind im Demo-Modus vorhanden:
//  -> User: admin  Passwort: test123   Rang: Admin
//  -> User: joker  Passwort: mosbach18 Rang: Editor
//  -> User: rolf   Passwort: gopher    Rang: Editor
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	inputUsername := r.FormValue("username")
	inputPassword := r.FormValue("password")
	redirectTarget := "/loginView"

	if inputUsername != "" && inputPassword != "" {
		// Lade die aktuellen Daten für registrierte Nutzer.
		// Überprüfe, ob der Benutzer registriert ist und prüfe dann das Passwort.
		for _, v := range users.Users {
			if v.Username == inputUsername {

				// Berechne den Hashwert des Input-Passworts mit dem Salt-Wert, der zum verifizerten Benutzer gehört.
				// Die Authentifizerung war erfolgreich, wenn dieser Hashwert mit dem Gespeicherten übereinstimmt.
				if GetHash(inputPassword, v.Salt) == v.Password {

					// Die Berechnung eines zufälligen Salt-Wertes wird hier als sessionToken genutzt.
					sessionToken := generateSalt()
					sessionUsers[sessionToken] = v

					http.SetCookie(w, &http.Cookie{
						Name:  "session_token",
						Value: sessionToken,
						Path:  "/",
					})

					redirectTarget = "/dashboard"
				}
			}
		}
	}

	http.Redirect(w, r, redirectTarget, 302)
}

// A-3.2:
// Der Zugang für die Bearbeiter soll durch Benutzernamen und Passwort geschützt sein.
//
// Beenden einer Sitzung und ausloggen des Benutzers durch löschen der Session-Cookies.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	for i := range sessionUsers {
		delete(sessionUsers, i)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/", 302)
}

// Registrieren und speichern eines neuen Benutzers in users.json.
// Dabei wird der Hashwert des Passworts mit einem persönlichen Salt-Wert verschleiert.
// Der Salt-Wert wird für spätere Abgleiche beider Hashwerte benötigt und folglich gespeichert.
func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	hash, salt := HashString(r.FormValue("password"))

	users.Users = append(users.Users, User{
		ID:       len(users.Users),
		Tier:     2,
		Username: r.FormValue("username"),
		Password: hash,
		Salt:     salt,
	})

	usersJson, err := json.Marshal(users)
	HandleError(err)
	err = ioutil.WriteFile(GetAssetsDir()+"users.json", usersJson, 0744)
	HandleError(err)

	http.Redirect(w, r, "/", 302)
}
