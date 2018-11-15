package ticket

import (
	"encoding/json"
	"io/ioutil"
)

type Id struct {
	FreeId int `json:"free id"`
}

func NewId() (freeId int) {
	var Id Id
	filename := "./assets/idGenerator_resource.json"
	encodedId, errRead := ioutil.ReadFile(filename)
	errorCheck(errRead)
	err := json.Unmarshal(encodedId, &Id)
	errorCheck(err)
	freeId = Id.FreeId
	Id.FreeId += 1

	encodedId, errEnc := json.Marshal(Id)
	errorCheck(errEnc)
	errWrite := ioutil.WriteFile(filename, encodedId, 0600)
	errorCheck(errWrite)
	return freeId
}

func Reset() {
	Id := Id{1}
	filename := "./assets/idGenerator_resource.json"
	encodedId, errEnc := json.Marshal(Id)
	errorCheck(errEnc)
	errWrite := ioutil.WriteFile(filename, encodedId, 0600)
	errorCheck(errWrite)
}

/*
func CheckJson(filename string) {

	encodedFile, errRead := ioutil.ReadFile(filename)
	errorCheck(errRead)
	if bytes.Equal(encodedFile, []byte{}) || countTickets() == 0 {

	}
}*/
