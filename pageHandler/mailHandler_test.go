package pageHandler

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	sessionHandler.BackupEnvironment()
	sessionHandler.DemoMode()
}

func teardown() {
	sessionHandler.RestoreEnvironment()
}

func TestCreateNewTicket(t *testing.T) {
	//TODO: Handler Funktion testen
}

func TestValidateMail(t *testing.T) {
	validMail := Mail{Email: "testing@golang.org", Subject: "Unit Test", Content: "Valid Email"}
	invalidMail := Mail{Email: "testing@golang.org"}

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
