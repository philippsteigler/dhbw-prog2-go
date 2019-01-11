package pageHandler

import (
	"bytes"
	"de/vorlesung/projekt/crew/sessionHandler"
	"de/vorlesung/projekt/crew/ticket"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Matrikelnummern:
//
// 3333958
// 3880065
// 8701350

func setup() {
	sessionHandler.BackupEnvironment()
	sessionHandler.DemoMode()
}

func teardown() {
	sessionHandler.RestoreEnvironment()
}

func TestCreateNewTicket(t *testing.T) {
	setup()
	defer teardown()

	//Erzeugen einer Testmail
	mail := map[string]string{"email": "test@home.com", "subject": "CreateNewTicket Test", "content": "Ein weiterer Test."}
	jsonMail, err := json.Marshal(mail)
	sessionHandler.HandleError(err)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "https://localhost:8000/ticket", bytes.NewBuffer(jsonMail))
	CreateNewTicket(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	newTicket := ticket.GetTicket(3)

	assert.Equal(t, "test@home.com", newTicket.Entries[0].Creator)
	assert.Equal(t, "CreateNewTicket Test", newTicket.Subject)
	assert.Equal(t, "Ein weiterer Test.", newTicket.Entries[0].Content)
}

func TestValidateMail(t *testing.T) {
	validMail := ticket.Mail{Email: "testing@golang.org", Subject: "Unit Test", Content: "Valid Email"}
	invalidMail := ticket.Mail{Email: "testing@golang.org"}

	assert.NoError(t, validateMail(validMail))
	assert.Error(t, validateMail(invalidMail))
}

func TestParseSubject(t *testing.T) {
	testSubject1 := parseSubject("RE: Probleme mit Lotus Notes")
	testSubject2 := parseSubject("RE:Probleme mit IE 11")
	testSubject3 := parseSubject("RE:   Probleme mit Whitespaces  ")

	assert.Equal(t, "Probleme mit Lotus Notes", testSubject1)
	assert.Equal(t, "Probleme mit IE 11", testSubject2)
	assert.Equal(t, "Probleme mit Whitespaces", testSubject3)
}

func TestRefersToExistingTicket(t *testing.T) {
	setup()
	defer teardown()

	testSubjectOk1 := "RE: Test 1"
	testSubjectOk2 := "RE: test 1"
	testSubjectOk3 := "Re: Test 2"
	testSubjectNotOk := "RE: Test 3"
	emailOk1 := "bob@dhbw.de"
	emailOk2 := "alice@dhbw.de"
	emailNotOk := "chris@work.com"

	ok, id := refersToExistingTicket(testSubjectOk1, emailOk1)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, id)

	ok, id = refersToExistingTicket(testSubjectOk2, emailOk1)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, id)

	ok, id = refersToExistingTicket(testSubjectOk3, emailOk2)
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, id)

	notOk, _ := refersToExistingTicket(testSubjectNotOk, emailOk1)
	assert.Equal(t, false, notOk)

	notOk, _ = refersToExistingTicket(testSubjectOk1, emailNotOk)
	assert.Equal(t, false, notOk)
}

func TestMailsRetrieveMails(t *testing.T) {
	setup()
	defer teardown()

	expectedResponse := "[{\"Email\":\"testing@home.ru\",\"Subject\":\"Unit Test 1\",\"Content\":\"Test\"}," +
		"{\"Email\":\"testing@work.com\",\"Subject\":\"Unit Test 2\",\"Content\":\"Test Test\"}," +
		"{\"Email\":\"testing@dhbw.de\",\"Subject\":\"Unit Test 3\",\"Content\":\"Test Test Test\"}]"

	response := httptest.NewRecorder()
	getRequest := httptest.NewRequest(http.MethodGet, "https://localhost:8000/mail", nil)
	Mails(response, getRequest)

	assert.Equal(t, http.StatusOK, response.Code)

	getResponse, _ := ioutil.ReadAll(response.Body)

	assert.Equal(t, expectedResponse, string(getResponse))
}

func TestMailsSentMails(t *testing.T) {
	setup()
	defer teardown()

	sentMails := "[{\"Email\":\"testing@home.ru\",\"Subject\":\"Unit Test 1\",\"Content\":\"Test\"}," +
		"{\"Email\":\"testing@work.com\",\"Subject\":\"Unit Test 2\",\"Content\":\"Test Test\"}," +
		"{\"Email\":\"testing@dhbw.de\",\"Subject\":\"Unit Test 3\",\"Content\":\"Test Test Test\"}]"

	response := httptest.NewRecorder()
	getRequest := httptest.NewRequest(http.MethodPost, "https://localhost:8000/mail", bytes.NewBufferString(sentMails))
	Mails(response, getRequest)

	assert.Equal(t, http.StatusOK, response.Code)

	assert.Equal(t, 0, len(*ticket.GetAllMails()))
}
