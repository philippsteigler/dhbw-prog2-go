package ticket

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func setup() {
	sessionHandler.BackupEnvironment()
	sessionHandler.DemoMode()
}

func teardown() {
	sessionHandler.RestoreEnvironment()
}

func TestParseFilename(t *testing.T) {
	filename := "123.json"
	id := parseFilename(filename)
	assert.Equal(t, 123, id)
}

func TestReadTicket(t *testing.T) {
	setup()
	defer teardown()

	storedTicket := *readTicket(1)
	assert.IsType(t, []Ticket{}, storedTicket)
}

func TestWriteTicket(t *testing.T) {
	setup()
	defer teardown()

	ticketDir := sessionHandler.GetAssetsDir() + "/tickets"
	files, err := ioutil.ReadDir(ticketDir)
	sessionHandler.HandleError(err)

	assert.Equal(t, 2, len(files))

	NewTicket("Test", "Bob", "Test")
	files, err = ioutil.ReadDir(ticketDir)
	sessionHandler.HandleError(err)

	assert.Equal(t, 3, len(files))
}

func TestWriteMail(t *testing.T) {
	setup()
	defer teardown()

	mailsDir := sessionHandler.GetAssetsDir() + "/mails"
	files, err := ioutil.ReadDir(mailsDir)
	sessionHandler.HandleError(err)

	assert.Equal(t, 3, len(files))

	mail := Mail{Email: "testing@home.ru", Subject: "Unit Test", Content: "Test Test Test"}
	writeMail(mail)
	files, err = ioutil.ReadDir(mailsDir)
	sessionHandler.HandleError(err)

	assert.Equal(t, 4, len(files))
	DeleteMail(mail)
}

func TestGetTicket(t *testing.T) {
	setup()
	defer teardown()

	testTicket := *GetTicket(1)
	assert.IsType(t, Ticket{}, testTicket)
	assert.Equal(t, 1, testTicket.Id)
	assert.Equal(t, Open, testTicket.Status)
	assert.Equal(t, 0, testTicket.EditorId)
	assert.Equal(t, "Test 1", testTicket.Subject)
	assert.Equal(t, "Test", testTicket.Entries[0].Content)
	assert.Equal(t, "bob@dhbw.de", testTicket.Entries[0].Creator)
}

func TestNewId(t *testing.T) {
	setup()
	defer teardown()

	id := newId("/tickets")
	assert.Equal(t, 3, id)

	NewTicket("Test", "Bob", "Test")
	id = newId("/tickets")
	assert.Equal(t, 4, id)
}

func TestTicketExist(t *testing.T) {
	setup()
	defer teardown()

	exist := ticketExist(1)
	assert.Equal(t, true, exist)

	exist = ticketExist(3)
	assert.Equal(t, false, exist)
}

func TestNewTicket(t *testing.T) {
	setup()
	defer teardown()

	NewTicket("Test", "Bob", "Test")
	testTicket := *GetTicket(3)

	assert.NotEmpty(t, testTicket)
}

func TestNewEntry(t *testing.T) {
	entry := NewEntry("Bob", "Test")
	assert.IsType(t, Entry{}, entry)
	assert.Equal(t, "Test", entry.Content)
	assert.Equal(t, "Bob", entry.Creator)
}

func TestAppendEntry(t *testing.T) {
	setup()
	defer teardown()

	AppendEntry(1, "Chris", "Test", false)
	AppendEntry(1, "Petra", "Test", false)
	AppendEntry(1, "Bob", "Test", false)

	testTicket := GetTicket(1)
	assert.Equal(t, 4, len(testTicket.Entries))
}

func TestGetTicketsByEditorId(t *testing.T) {
	setup()
	defer teardown()

	TakeTicket(2, 7)

	orderedTickets := *GetTicketsByEditorId(0)
	assert.Equal(t, 1, len(orderedTickets))

	orderedTickets = *GetTicketsByEditorId(7)
	assert.Equal(t, 1, len(orderedTickets))
}

func TestTakeTicket(t *testing.T) {
	setup()
	defer teardown()

	TakeTicket(1, 7)
	testTicket := GetTicket(1)
	assert.Equal(t, InProcess, testTicket.Status)
	assert.Equal(t, 7, testTicket.EditorId)
}

func TestGetAllOpenTickets(t *testing.T) {
	setup()
	defer teardown()

	openTickets := *GetAllOpenTickets()
	assert.Equal(t, 2, len(openTickets))
}

func TestUnhandTicket(t *testing.T) {
	setup()
	defer teardown()

	TakeTicket(1, 7)

	UnhandTicket(1)
	testTicket := GetTicket(1)

	assert.Equal(t, Open, testTicket.Status)
	assert.Equal(t, 0, testTicket.EditorId)
}

func TestDelegateTicket(t *testing.T) {
	setup()
	defer teardown()

	DelegateTicket(1, 4)
	testTicket := GetTicket(1)
	assert.Equal(t, InProcess, testTicket.Status)
	assert.Equal(t, 4, testTicket.EditorId)
}

func TestMergeTickets(t *testing.T) {
	setup()
	defer teardown()

	MergeTickets(1, 2)
	testTicket := GetTicket(1)
	assert.Equal(t, 2, len(testTicket.Entries))
}

func TestDeleteTicket(t *testing.T) {
	setup()
	defer teardown()

	deleteTicket(1)

	ticketDir := sessionHandler.GetAssetsDir() + "/tickets"
	files, err := ioutil.ReadDir(ticketDir)
	sessionHandler.HandleError(err)

	assert.Equal(t, 1, len(files))
}

func TestGetTicketHistory(t *testing.T) {
	setup()
	defer teardown()

	NewTicket("Test", "Bob", "Test")
	TakeTicket(3, 2)
	UnhandTicket(3)
	DelegateTicket(3, 4)

	ticketHistory := *GetTicketHistory(3)

	assert.Equal(t, 4, len(ticketHistory))
	assert.Equal(t, 0, ticketHistory[0].EditorId)
	assert.Equal(t, 2, ticketHistory[1].EditorId)
	assert.Equal(t, 4, ticketHistory[3].EditorId)
}

func TestGetAllTickets(t *testing.T) {
	setup()
	defer teardown()

	assert.Equal(t, []Ticket{*GetTicket(1), *GetTicket(2)}, *GetAllTickets())
}

func TestGetAllMailsAndDeleteMail(t *testing.T) {
	setup()
	defer teardown()

	mail1 := Mail{Email: "testing@home.ru", Subject: "Unit Test 1", Content: "Test"}
	mail2 := Mail{Email: "testing@work.com", Subject: "Unit Test 2", Content: "Test Test"}
	mail3 := Mail{Email: "testing@dhbw.de", Subject: "Unit Test 3", Content: "Test Test Test"}

	writeMail(mail1)
	writeMail(mail2)
	writeMail(mail3)

	mails := *GetAllMails()

	assert.IsType(t, []Mail{}, mails)
	assert.Equal(t, 6, len(mails))
	assert.Equal(t, "testing@home.ru", mails[3].Email)
	assert.Equal(t, "Unit Test 2", mails[4].Subject)
	assert.Equal(t, "Test Test Test", mails[5].Content)
}
