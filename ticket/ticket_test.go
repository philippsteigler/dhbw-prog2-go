package ticket

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func CreateDefaultEnv() {
	//Alle Tickets aus dem Ordner "tickets" Löschen
	ticketFiles, err := ioutil.ReadDir(sessionHandler.GetAssetsDir() + "tickets")
	sessionHandler.HandleError(err)

	for _, file := range ticketFiles {
		err := os.Remove(sessionHandler.GetAssetsDir() + "tickets/" + file.Name())
		sessionHandler.HandleError(err)
	}
}

func TestParseFilename(t *testing.T) {
	filename := "123.json"
	id := parseFilename(filename)
	assert.Equal(t, 123, id)
}

func TestReadTicket(t *testing.T) {
	CreateDefaultEnv()

	NewTicket("Test", "Bob", "Test")
	storedTicket := *readTicket(1)
	assert.IsType(t, []Ticket{}, storedTicket)
}

func TestWriteTicket(t *testing.T) {
	CreateDefaultEnv()

	ticketDir := sessionHandler.GetAssetsDir() + "/tickets"
	files, err := ioutil.ReadDir(ticketDir)
	sessionHandler.HandleError(err)

	assert.Equal(t, 0, len(files))

	NewTicket("Test", "Bob", "Test")
	files, err = ioutil.ReadDir(ticketDir)
	sessionHandler.HandleError(err)

	assert.Equal(t, 1, len(files))
}

func TestGetTicket(t *testing.T) {
	CreateDefaultEnv()
	NewTicket("Test", "Bob", "TestGetTicket")

	testTicket := *GetTicket(1)
	assert.IsType(t, Ticket{}, testTicket)
	assert.Equal(t, 1, testTicket.Id)
	assert.Equal(t, Open, testTicket.Status)
	assert.Equal(t, 0, testTicket.EditorId)
	assert.Equal(t, "Test", testTicket.Subject)
	assert.Equal(t, "TestGetTicket", testTicket.Entries[0].Content)
	assert.Equal(t, "Bob", testTicket.Entries[0].Creator)
}

func TestNewId(t *testing.T) {
	CreateDefaultEnv()

	id := newId()
	assert.Equal(t, 1, id)

	NewTicket("Test", "Bob", "Test")
	id = newId()
	assert.Equal(t, 2, id)
}

func TestTicketExist(t *testing.T) {
	CreateDefaultEnv()

	exist := ticketExist(1)
	assert.Equal(t, false, exist)

	NewTicket("Test", "Bob", "Test")
	exist = ticketExist(1)
	assert.Equal(t, true, exist)
}

func TestNewTicket(t *testing.T) {
	CreateDefaultEnv()

	NewTicket("Test", "Bob", "Test")
	testTicket := *GetTicket(1)

	assert.NotEmpty(t, testTicket)
}

func TestNewEntry(t *testing.T) {
	entry := NewEntry("Bob", "Test")
	assert.IsType(t, Entry{}, entry)
	assert.Equal(t, "Test", entry.Content)
	assert.Equal(t, "Bob", entry.Creator)
}

func TestAppendEntry(t *testing.T) {
	CreateDefaultEnv()

	NewTicket("Test", "Bob", "Test")

	AppendEntry(1, "Chris", "Test")
	AppendEntry(1, "Petra", "Test")
	AppendEntry(1, "Bob", "Test")

	testTicket := GetTicket(1)
	assert.Equal(t, 4, len(testTicket.Entries))
}

func TestGetTicketsByEditorId(t *testing.T) {
	CreateDefaultEnv()

	NewTicket("Test", "Alice", "Test")
	NewTicket("Test", "Bob", "Test")
	NewTicket("Test", "Chris", "Test")
	TakeTicket(2, 7)

	orderedTickets := *GetTicketsByEditorId(0)
	assert.Equal(t, 2, len(orderedTickets))

	orderedTickets = *GetTicketsByEditorId(7)
	assert.Equal(t, 1, len(orderedTickets))
}

func TestTakeTicket(t *testing.T) {
	CreateDefaultEnv()

	NewTicket("Test", "Bob", "Test")
	TakeTicket(1, 7)
	testTicket := GetTicket(1)
	assert.Equal(t, InProcess, testTicket.Status)
	assert.Equal(t, 7, testTicket.EditorId)
}

func TestGetAllOpenTickets(t *testing.T) {
	CreateDefaultEnv()

	NewTicket("Test", "Alice", "Test")
	NewTicket("Test", "Bob", "Test")
	NewTicket("Test", "Chris", "Test")

	openTickets := *GetAllOpenTickets()
	assert.Equal(t, 3, len(openTickets))
}

func TestUnhandTicket(t *testing.T) {
	CreateDefaultEnv()

	NewTicket("Test", "Tim", "Test")
	TakeTicket(1, 7)

	UnhandTicket(1)
	testTicket := GetTicket(1)

	assert.Equal(t, Open, testTicket.Status)
	assert.Equal(t, 0, testTicket.EditorId)
}

func TestDelegateTicket(t *testing.T) {
	CreateDefaultEnv()

	NewTicket("Test", "Bob", "Test")
	DelegateTicket(1, 4)
	testTicket := GetTicket(1)
	assert.Equal(t, InProcess, testTicket.Status)
	assert.Equal(t, 4, testTicket.EditorId)
}

func TestMergeTickets(t *testing.T) {
	CreateDefaultEnv()
	NewTicket("Test", "Bob", "Test")
	NewTicket("Test", "Bob", "Test")

	MergeTickets(1, 2)
	testTicket := GetTicket(1)
	assert.Equal(t, 2, len(testTicket.Entries))
}

func TestDeleteTicket(t *testing.T) {
	CreateDefaultEnv()

	NewTicket("Test", "Bob", "Test")
	deleteTicket(1)

	ticketDir := sessionHandler.GetAssetsDir() + "/tickets"
	files, err := ioutil.ReadDir(ticketDir)
	sessionHandler.HandleError(err)

	assert.Equal(t, 0, len(files))
}
