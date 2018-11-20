package ticket

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"ticketBackend/sessionHandler"
)

func CreateDefaultEnv() {
	//Alle Tickets aus dem Ordner "tickets" Löschen
	ticketFiles, err := ioutil.ReadDir(sessionHandler.GetAssetsDir() + "tickets")
	sessionHandler.HandleError(err)

	for _, file := range ticketFiles {
		err := os.Remove(sessionHandler.GetAssetsDir() + "tickets/" + file.Name())
		sessionHandler.HandleError(err)
	}

	//Setzt die Datei ticketId_resource.json in default Zustand
	id = Id{1}
	filename := sessionHandler.GetAssetsDir() + "ticketId_resource.json"
	encodedId, errEnc := json.Marshal(id)
	sessionHandler.HandleError(errEnc)
	errWrite := ioutil.WriteFile(filename, encodedId, 0600)
	sessionHandler.HandleError(errWrite)
}

func TestCountTickets(t *testing.T) {
	assert.Equal(t, 0, countTickets())
}

func TestNewTicket(t *testing.T) {
	CreateDefaultEnv()
	NewTicket("Test", "Bob", "Test")
	NewTicket("Test", "Bob", "Test")

	ticket1 := *readTicket(1)
	ticket2 := *readTicket(2)

	assert.NotEqual(t, ticket1, ticket2)
	CreateDefaultEnv()
}

func TestDeleteTicket(t *testing.T) {
	CreateDefaultEnv()
	NewTicket("Test", "Bob", "Test")
	DeleteTicket(1)
	assert.Equal(t, 0, countTickets())
	CreateDefaultEnv()
}

func TestNewEntry(t *testing.T) {
	entry := *NewEntry("Bob", "Test")
	assert.IsType(t, Entry{}, entry)
}

func TestReadTicket(t *testing.T) {
	CreateDefaultEnv()
	NewTicket("Test", "Bob", "Test")
	readTicket(1)
	assert.IsType(t, Ticket{}, ticket)
	CreateDefaultEnv()
}

func TestWriteTicket(t *testing.T) {
	CreateDefaultEnv()
	NewTicket("Test", "Bob", "Test")
	writeTicket(&ticket)
	assert.Equal(t, 1, countTickets())
	CreateDefaultEnv()
}

func TestAppendEntry(t *testing.T) {
	CreateDefaultEnv()
	NewTicket("Test", "Bob", "Test")
	readTicket(1)

	AppendEntry(1, "Chris", "Test")
	AppendEntry(1, "Petra", "Test")
	AppendEntry(1, "Bob", "Test")

	readTicket(1)
	assert.Equal(t, 4, len(ticket.Entries))
	CreateDefaultEnv()
}

func TestTakeTicket(t *testing.T) {
	CreateDefaultEnv()
	NewTicket("Test", "Bob", "Test")
	TakeTicket(1, 7)
	readTicket(1)
	assert.Equal(t, InProcess, ticket.Status)
	assert.Equal(t, 7, ticket.EditorId)
	CreateDefaultEnv()
}

func TestGetTickets(t *testing.T) {
	CreateDefaultEnv()
	NewTicket("Test", "Bob", "Test")
	NewTicket("Test", "Bob", "Test")
	NewTicket("Test", "Bob", "Test")

	TakeTicket(2, 7)
	orderedTickets = *GetTickets(Open)
	assert.Equal(t, 2, len(orderedTickets))

	orderedTickets = *GetTickets(InProcess, Closed)
	assert.Equal(t, 1, len(orderedTickets))

	orderedTickets = *GetTickets()
	assert.Equal(t, 3, len(orderedTickets))
	CreateDefaultEnv()
}

func TestParseFilename(t *testing.T) {
	filename := "123.json"
	id := parseFilename(filename)
	assert.Equal(t, 123, id)
}

func TestUnhandTicket(t *testing.T) {
	CreateDefaultEnv()
	NewTicket("Test", "Bob", "Test")
	TakeTicket(1, 7)

	UnhandTicket(1)
	assert.Equal(t, Open, ticket.Status)
	assert.Equal(t, 0, ticket.EditorId)
	CreateDefaultEnv()
}

func TestDelegateTicket(t *testing.T) {
	CreateDefaultEnv()
	NewTicket("Test", "Bob", "Test")
	DelegateTicket(1, 4)
	assert.Equal(t, InProcess, ticket.Status)
	assert.Equal(t, 4, ticket.EditorId)
	CreateDefaultEnv()
}

func TestMergeTickets(t *testing.T) {
	CreateDefaultEnv()
	NewTicket("Test", "Bob", "Test")
	NewTicket("Test", "Bob", "Test")

	MergeTickets(1, 2)
	readTicket(1)
	assert.Equal(t, 2, len(ticket.Entries))
	CreateDefaultEnv()
}

func TestNewId(t *testing.T) {
	CreateDefaultEnv()
	newId := newId()
	assert.Equal(t, 1, newId)
	assert.Equal(t, 2, id.FreeId)
	CreateDefaultEnv()
}
