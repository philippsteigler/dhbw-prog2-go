package main

import (
	"crypto/tls"
	"crypto/x509"
	"de/vorlesung/projekt/crew/sessionHandler"
	"de/vorlesung/projekt/crew/ticket"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Matrikelnummern:
//
// 3333958
// 3880065
// 8701350

// A-7.4
// Es soll ein einfaches CLI-Tool zur Ausgabe noch zu versendender Mails geben
//
// GET Request zur Ausgabe aller ausstehenden Mails
func main() {
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

	// GET Request an die entsprechende API
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:8000/mails", nil)
	response, httpGetErr := client.Do(req)

	// Überprüfung auf Fehler. Bei erfolgreichem Request werden die Mails ausgegeben (falls welche existieren)
	if httpGetErr != nil {

		fmt.Printf("The HTTP GET request to: https://localhost:8000/mails failed with error %s\n", httpGetErr)
	} else {

		resp, _ := ioutil.ReadAll(response.Body)
		if string(resp) == "null" {
			fmt.Println("No mails found.")
		} else {
			printMails(resp)
		}
	}
}

// Funktion zur leserlichen Ausgabe der Mails
func printMails(responseBody []byte) {
	var mails []ticket.Mail
	err := json.Unmarshal(responseBody, &mails)
	sessionHandler.HandleError(err)
	count := 1
	for _, mail := range mails {
		fmt.Printf("Mail Nr.: %v\n", count)
		fmt.Printf("Recipient: %s\n", mail.Email)
		fmt.Printf("Subject: %s\n", mail.Subject)
		fmt.Printf("Content: %s\n", mail.Content)
		fmt.Printf("\n")
		count++
	}
}
