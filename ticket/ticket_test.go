package ticket

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCountTickets(t *testing.T) {
	assert.Equal(t, 0, countTickets())
}

func TestNewTicket(t *testing.T) {
	NewTicket("Test", "Bob", "Test")
	NewTicket("Test", "Bob", "Test")

	ticket1 := *readTicket(1)
	ticket2 := *readTicket(2)

	assert.NotEqual(t, ticket1, ticket2)

	DeleteTicket(1)
	DeleteTicket(2)
	Reset()
}

func TestDeleteTicket(t *testing.T) {
	NewTicket("Test", "Bob", "Test")
	DeleteTicket(1)
	assert.Equal(t, 0, countTickets())
	Reset()
}

func TestNewEntry(t *testing.T) {
	entry1 := *NewEntry("Bob", "Test")
	entry2 := Entry{Date: time.Now().Local().Format("2006-01-02T15:04:05Z07:00"), Creator: "Bob", Content: "Test"}
	assert.Equal(t, entry1, entry2)
}

func TestReadTicket(t *testing.T) {
	NewTicket("Test", "Bob", "Test")
	readTicket(1)
	assert.IsType(t, Ticket{}, ticket)
	DeleteTicket(1)
	Reset()
}

func TestWriteTicket(t *testing.T) {
	NewTicket("Test", "Bob", "Test")
	writeTicket(&ticket)
	assert.Equal(t, 1, countTickets())
	DeleteTicket(1)
	Reset()
}

func TestAppendEntry(t *testing.T) {
	NewTicket("Test", "Bob", "Test")
	readTicket(1)

	AppendEntry(1, "Chris", "Test")
	AppendEntry(1, "Petra", "Test")
	AppendEntry(1, "Bob", "Test")

	readTicket(1)
	assert.Equal(t, 4, len(ticket.Entries))
	DeleteTicket(1)
	Reset()
}

func TestTakeTicket(t *testing.T) {
	NewTicket("Test", "Bob", "Test")
	TakeTicket(1, 7)
	readTicket(1)
	assert.Equal(t, InProcess, ticket.Status)
	assert.Equal(t, 7, ticket.EditorId)
	DeleteTicket(1)
	Reset()
}

func TestGetTickets(t *testing.T) {
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

	DeleteTicket(1)
	DeleteTicket(2)
	DeleteTicket(3)
	Reset()
}

func TestParseFilename(t *testing.T) {
	filename := "123.json"
	id := parseFilename(filename)
	assert.Equal(t, 123, id)
}

func TestUnhandTicket(t *testing.T) {
	NewTicket("Test", "Bob", "Test")
	TakeTicket(1, 7)

	UnhandTicket(1)
	assert.Equal(t, Open, ticket.Status)
	assert.Equal(t, 0, ticket.EditorId)
	DeleteTicket(1)
	Reset()
}

func TestDelegateTicket(t *testing.T) {
	NewTicket("Test", "Bob", "Test")
	DelegateTicket(1, 4)
	assert.Equal(t, InProcess, ticket.Status)
	assert.Equal(t, 4, ticket.EditorId)
	DeleteTicket(1)
	Reset()
}

func TestMergeTickets(t *testing.T) {
	NewTicket("Test", "Bob", "Test")
	NewTicket("Test", "Bob", "Test")

	MergeTickets(1, 2)
	readTicket(1)
	assert.Equal(t, 2, len(ticket.Entries))
	DeleteTicket(1)
	Reset()
}

func TestGetTicket(t *testing.T) {
	NewTicket("Test", "Bob", "Test")
	ticket = *GetTicket(1)
	assert.IsType(t, Ticket{}, ticket)
	DeleteTicket(1)
	Reset()
}
