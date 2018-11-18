package ticket

import (
	"encoding/json"
	"io/ioutil"
)

type Id struct {
	FreeId int `json:"free id"`
}

var id Id

//Gibt eine einzigartige, gültige ID zurück
func NewId() int {
	//Hier befindet sich die gültige ID und wird ausgelesen
	filename := "./assets/idGenerator_resource.json"
	encodedId, errRead := ioutil.ReadFile(filename)
	errorCheck(errRead)
	err := json.Unmarshal(encodedId, &id)
	errorCheck(err)

	//In "freeId" wird die ID gespeichert und die Zahl in der Datei um eins erhöht (und zurückgeschrieben)
	freeId := id.FreeId
	id.FreeId += 1

	encodedId, errEnc := json.Marshal(id)
	errorCheck(errEnc)
	errWrite := ioutil.WriteFile(filename, encodedId, 0600)
	errorCheck(errWrite)
	return freeId
}

//Wird aktuell für die Tests benötigt, da für das durchführen der meisten Tests ein default Stand
//in der Datei "idGenerator_resource.json" benötigt wird
func Reset() {
	id = Id{1}
	filename := "./assets/idGenerator_resource.json"
	encodedId, errEnc := json.Marshal(id)
	errorCheck(errEnc)
	errWrite := ioutil.WriteFile(filename, encodedId, 0600)
	errorCheck(errWrite)
}

//TODO Funktion zur Überprüfung ob resource json file valide ist
/*
func CheckJson(filename string) {

	encodedFile, errRead := ioutil.ReadFile(filename)
	errorCheck(errRead)
	if bytes.Equal(encodedFile, []byte{}) || countTickets() == 0 {

	}
}*/
