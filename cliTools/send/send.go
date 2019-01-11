package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"de/vorlesung/projekt/crew/sessionHandler"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Matrikelnummern:
//
// 3333958
// 3880065
// 8701350

// A-6.6
// Es soll ein einfaches CLI-Tool zur Abgabe von Nachrichten an den Server geben
func main() {
	response, httpPostErr := sendPostRequest(parseUserInput())

	// Überprüfung auf Fehler
	if httpPostErr != nil {
		fmt.Printf("The HTTP POST request to: https://localhost:4443/ticket failed (status code: %s) with error %s\n", response.Status, httpPostErr)
	} else {
		fmt.Printf("The HTTP POST request finished with status code: %s.\n", response.Status)
	}
}

// Eingabe über das CLI wird geparsed und auf vollständigkeit überprüft
func parseUserInput() []string {
	email := flag.String("email", "", "Email of creator. (Required)")
	subject := flag.String("subject", "", "Subject of ticket. (Required)")
	content := flag.String("content", "", "Content of ticket. (Required)")

	flag.Parse()

	//Für jedes Flag muss es eine Eingabe geben, sonst wird der Request nicht gesendet
	//Dem User wird dann die Hilfe angezeigt
	if *email == "" || *subject == "" || *content == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return []string{*email, *subject, *content}
}

// Post Request zur Abgabe einer Mail
func sendPostRequest(userInput []string) (response *http.Response, httpPostErr error) {
	// Beim senden des Requests an den Server trat folgender Fehler auf: x509: certificate signed by unknown authority
	// Dieser Fehler trat auf, da es sich beim Zertifikat des Servers um ein self-signed Zertifikat handelt
	// Entsprechend muss es dem CertPool hinzugefügt werden
	// Credits: https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/

	// Pfad des self-signed Zertifikats
	selfSignedCert := sessionHandler.GetAssetsDir() + "certificates/cert.pem"

	// SystemCertPool laden bzw. erzeugen
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// self-signed Zertifikat einlesen
	certs, err := ioutil.ReadFile(selfSignedCert)
	sessionHandler.HandleError(err)

	// self-signed Zertifikat dem SystemCertPool hinzufügen
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}

	// dem neuen SystemCertPool vertrauen
	config := &tls.Config{RootCAs: rootCAs}
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: tr}

	// Eigentliche Funktion zum senden des POST-Requests

	// Erzeugen eines JSON-Objekts
	input := map[string]string{"email": userInput[0], "subject": userInput[1], "content": userInput[2]}
	jsonInput, err := json.Marshal(input)
	sessionHandler.HandleError(err)

	// POST-Request an den Server
	req, _ := http.NewRequest(http.MethodPost, "https://localhost:4443/ticket", bytes.NewBuffer(jsonInput))
	response, httpPostErr = client.Do(req)

	return
}
