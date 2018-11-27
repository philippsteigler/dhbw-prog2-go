package sessionHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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

var sessionUsers = map[string]User{}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

// Gib den relativen Pfad zum Ressourcen-Verzeichnis zur Laufzeit zurück.
func GetAssetsDir() string {
	path, err := os.Getwd()
	HandleError(err)

	// Fallunterscheidung für Aufruf über main.go oder *_test.go aus Unterverzeichnissen.
	if filepath.Base(path) == "ticketBackend" {
		return "./assets/"
	} else {
		return "../assets/"
	}
}

// Lies users.json und importiere alle Benutzerdaten nach &users.
func loadUserData() *UserAccounts {
	var users UserAccounts

	userData, err := ioutil.ReadFile(GetAssetsDir() + "users.json")
	HandleError(err)

	err = json.Unmarshal(userData, &users)
	HandleError(err)

	return &users
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

	return nil
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
		users := *loadUserData()

		// Überprüfe, ob der Benutzer registriert ist und prüfe dann das Passwort.
		for _, i := range users.Users {
			if i.Username == inputUsername {

				// Berechne den Hashwert des Input-Passworts mit dem Salt-Wert, der zum verifizerten Benutzer gehört.
				// Die Authentifizerung war erfolgreich, wenn dieser Hashwert mit dem Gespeicherten übereinstimmt.
				if GetHash(inputPassword, i.Salt) == i.Password {

					// Die Berechnung eines zufälligen Salt-Wertes wird hier als sessionToken genutzt.
					sessionToken := generateSalt()
					sessionUsers[sessionToken] = i

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
	users := *loadUserData()

	users.Users = append(users.Users, User{
		ID:       len(users.Users),
		Tier:     2,
		Username: r.FormValue("username"),
		Password: hash,
		Salt:     salt,
	})

	usersJson, err := json.Marshal(users)
	HandleError(err)
	err = ioutil.WriteFile(GetAssetsDir()+"users.json", usersJson, 0644)
	HandleError(err)

	http.Redirect(w, r, "/", 302)
}
