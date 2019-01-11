package sessionHandler

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

// Matrikelnummern:
//
// 3333958
// 3880065
// 8701350

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// A-3.4:
// Es soll "salting" eingesetzt werden.
//
// Generiere einen 16-Byte langen String aus zufälligen Zeichen als Salt-Wert.
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func generateSalt() string {
	salt := make([]byte, 16)

	for i := range salt {
		salt[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(salt)
}

// A-3.3:
// Die Passwörter dürfen nicht im Klartext gespeichert werden.
//
// Berechne den Hashwert eines Passworts unter Verwendung von salting.
func HashString(text string) (string, string) {
	salt := generateSalt()

	// Der Salt-Wert wird an das Passwort im Klartext angehängt.
	// Dieser wird später als Schlüssel zur Berechnung des passenden Hashwerts benötigt.
	tmp := []string{text, salt}
	saltedString := strings.Join(tmp, "")

	bytes := sha256.Sum256([]byte(saltedString))
	hashedString := base64.URLEncoding.EncodeToString(bytes[:])

	return hashedString, salt
}

// A-3.3:
// Die Passwörter dürfen nicht im Klartext gespeichert werden.
//
// Berechne den Hashwert eines Strings unter Verwendung eines spezifischen Salt-Wertes.
// Diese Funktion wird beim Abgleich mit gespeicherten Passwort-Hashes verwendet.
func GetHash(text, salt string) string {
	tmp := []string{text, salt}
	saltedString := strings.Join(tmp, "")

	bytes := sha256.Sum256([]byte(saltedString))
	hashedString := base64.URLEncoding.EncodeToString(bytes[:])

	return hashedString
}
