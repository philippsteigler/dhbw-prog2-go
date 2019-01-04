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

func main() {
	sendPostRequest(parseUserInput())
}

func parseUserInput() []string {
	email := flag.String("email", "", "Email of creator. (Required)")
	subject := flag.String("subject", "", "Subject of ticket. (Required)")
	content := flag.String("content", "", "Content of ticket. (Required)")

	flag.Parse()

	if *email == "" || *subject == "" || *content == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return []string{*email, *subject, *content}
}

func sendPostRequest(userInput []string) {
	//Beim senden des Requests an den Server trat folgender Fehler auf: x509: certificate signed by unknown authority
	//Dieser Fehler trat auf, da es sich beim Zertifikat des Servers um ein self-signed Zertifikat handelt
	//Entsprechend muss es dem CertPool hinzugefügt werden
	//Credits: https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/

	selfSignedCert := sessionHandler.GetAssetsDir() + "certificates/cert.pem"

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	certs, err := ioutil.ReadFile(selfSignedCert)
	sessionHandler.HandleError(err)

	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}

	config := &tls.Config{RootCAs: rootCAs}
	tr := &http.Transport{TLSClientConfig: config}
	client := &http.Client{Transport: tr}

	//Eigentliche Funktion zum senden des Requests

	input := map[string]string{"email": userInput[0], "subject": userInput[1], "content": userInput[2]}
	jsonInput, err := json.Marshal(input)
	sessionHandler.HandleError(err)

	req, _ := http.NewRequest(http.MethodPost, "https://localhost:8000/ticket", bytes.NewBuffer(jsonInput))
	response, httpPostErr := client.Do(req)
	sessionHandler.HandleError(httpPostErr)
	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(data))
}
