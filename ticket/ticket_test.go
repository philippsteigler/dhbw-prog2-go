package ticket

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

// Matrikelnummern:
//
// 3333958
// 3880065
// 8701350

func setup() {
	sessionHandler.BackupEnvironment()
	sessionHandler.DemoMode()
	time.Sleep(100 * time.Millisecond)
}

func teardown() {
	sessionHandler.RestoreEnvironment()
	time.Sleep(100 * time.Millisecond)
}

func TestParseFilename(t *testing.T) {
	filename := "123.json"
	id := parseFilename(filename)
	assert.Equal(t, 123, id)
}

func TestReadTicket(t *testing.T) {
	setup()

	storedTicket := *readTicket(1)
	assert.IsType(t, Ticket{}, storedTicket)

	teardown()
}

func TestWriteTicket(t *testing.T) {
	setup()

	ticketDir := sessionHandler.GetAssetsDir() + "/tickets"
	files, err := ioutil.ReadDir(ticketDir)
	sessionHandler.HandleError(err)

	assert.Equal(t, 2, len(files))

	NewTicket("Test", "Bob", "Test")
	files, err = ioutil.ReadDir(ticketDir)
	sessionHandler.HandleError(err)

	assert.Equal(t, 3, len(files))

	teardown()
}

func TestWriteMail(t *testing.T) {
	setup()

	mailsDir := sessionHandler.GetAssetsDir() + "/mails"
	files, err := ioutil.ReadDir(mailsDir)
	sessionHandler.HandleError(err)

	assert.Equal(t, 3, len(files))

	mail := Mail{Email: "testing@home.ru", Subject: "Unit Test", Content: "Test Test Test"}
	writeMail(mail)
	files, err = ioutil.ReadDir(mailsDir)
	sessionHandler.HandleError(err)

	assert.Equal(t, 4, len(files))

	teardown()
}

func TestGetTicket(t *testing.T) {
	setup()

	testTicket := *GetTicket(1)
	assert.IsType(t, Ticket{}, testTicket)
	assert.Equal(t, 1, testTicket.Id)
	assert.Equal(t, Open, testTicket.Status)
	assert.Equal(t, 0, testTicket.EditorId)
	assert.Equal(t, "Test 1", testTicket.Subject)
	assert.Equal(t, "Test", testTicket.Entries[0].Content)
	assert.Equal(t, "bob@dhbw.de", testTicket.Entries[0].Creator)

	teardown()
}

func TestNewId(t *testing.T) {
	setup()

	id := newId("/tickets")
	assert.Equal(t, 3, id)

	NewTicket("Test", "Bob", "Test")
	id = newId("/tickets")
	assert.Equal(t, 4, id)

	teardown()
}

func TestNewTicket(t *testing.T) {
	setup()

	NewTicket("Test", "Bob", "Test")
	testTicket := *GetTicket(3)

	assert.NotEmpty(t, testTicket)

	teardown()
}

func TestNewEntry(t *testing.T) {
	entry := NewEntry("Bob", "Test")
	assert.IsType(t, Entry{}, entry)
	assert.Equal(t, "Test", entry.Content)
	assert.Equal(t, "Bob", entry.Creator)
}

func TestAppendEntry(t *testing.T) {
	setup()

	AppendEntry(1, "Chris", "Test", false)
	AppendEntry(1, "Petra", "Test", false)
	AppendEntry(1, "Bob", "Test", false)

	testTicket := GetTicket(1)
	assert.Equal(t, 4, len(testTicket.Entries))

	teardown()
}

func TestGetTicketsByEditorId(t *testing.T) {
	setup()

	TakeTicket(2, 7)

	orderedTickets := *GetTicketsByEditorId(0)
	assert.Equal(t, 1, len(orderedTickets))

	orderedTickets = *GetTicketsByEditorId(7)
	assert.Equal(t, 1, len(orderedTickets))

	teardown()
}

func TestTakeTicket(t *testing.T) {
	setup()

	TakeTicket(1, 7)
	testTicket := GetTicket(1)
	assert.Equal(t, InProcess, testTicket.Status)
	assert.Equal(t, 7, testTicket.EditorId)

	teardown()
}

func TestGetAllOpenTickets(t *testing.T) {
	setup()

	openTickets := *GetAllOpenTickets()
	assert.Equal(t, 2, len(openTickets))

	teardown()
}

func TestUnhandTicket(t *testing.T) {
	setup()

	TakeTicket(1, 7)

	UnhandTicket(1)
	testTicket := GetTicket(1)

	assert.Equal(t, Open, testTicket.Status)
	assert.Equal(t, 0, testTicket.EditorId)

	teardown()
}

func TestDelegateTicket(t *testing.T) {
	setup()

	DelegateTicket(1, 4)
	testTicket := GetTicket(1)
	assert.Equal(t, InProcess, testTicket.Status)
	assert.Equal(t, 4, testTicket.EditorId)

	teardown()
}

func TestMergeTickets(t *testing.T) {
	setup()

	MergeTickets(1, 2)
	testTicket := GetTicket(1)
	assert.Equal(t, 2, len(testTicket.Entries))

	teardown()
}

func TestDeleteTicket(t *testing.T) {
	setup()

	deleteTicket(1)

	ticketDir := sessionHandler.GetAssetsDir() + "/tickets"
	files, err := ioutil.ReadDir(ticketDir)
	sessionHandler.HandleError(err)

	assert.Equal(t, 1, len(files))

	teardown()
}

func TestGetAllTickets(t *testing.T) {
	setup()

	assert.Equal(t, []Ticket{*GetTicket(1), *GetTicket(2)}, *GetAllTickets())

	teardown()
}

func TestGetAllMails(t *testing.T) {
	setup()

	mails := *GetAllMails()

	assert.IsType(t, []Mail{}, mails)
	assert.Equal(t, 3, len(mails))
	assert.Equal(t, "testing@home.ru", mails[0].Email)
	assert.Equal(t, "Unit Test 2", mails[1].Subject)
	assert.Equal(t, "Test Test Test", mails[2].Content)

	teardown()
}

func TestDeleteMail(t *testing.T) {
	setup()

	mails := *GetAllMails()
	assert.Equal(t, 3, len(mails))

	mail := Mail{Email: "testing@home.ru", Subject: "Unit Test 1", Content: "Test"}
	DeleteMail(mail)

	mails = *GetAllMails()
	assert.Equal(t, 2, len(mails))

	teardown()
}

func TestSetTicketToOpenIfClosed(t *testing.T) {
	setup()

	assert.Equal(t, Open, (*GetTicket(1)).Status)
	SetTicketToOpenIfClosed(1)
	assert.Equal(t, Open, (*GetTicket(1)).Status)

	CloseTicket(1)
	SetTicketToOpenIfClosed(1)
	assert.Equal(t, Open, (*GetTicket(1)).Status)

	teardown()
}

func TestCloseTicket(t *testing.T) {
	setup()

	assert.Equal(t, Open, (*GetTicket(1)).Status)

	CloseTicket(1)
	assert.Equal(t, Closed, (*GetTicket(1)).Status)

	teardown()
}
