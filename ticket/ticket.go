package ticket

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

type Status string

const (
	Open      Status = "offen"
	InProcess Status = "in Bearbeitung"
	Closed    Status = "geschlossen"
)

type Ticket struct {
	Id       int     `json:"id"`
	Subject  string  `json:"subject"`
	Status   Status  `json:"status"`
	EditorId int     `json:"editorId"`
	Entries  []Entry `json:"entries"`
}

type Entry struct {
	Date    string `json:"date"`
	Creator string `json:"creator"`
	Content string `json:"content"`
}

func errorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func NewTicket(subject string, entries []Entry) {
	id := createTicketId()
	ticket := Ticket{Id: id, Subject: subject, Status: Open, EditorId: 0, Entries: entries}
	encodedTicket, errEnc := json.Marshal(ticket)
	errorCheck(errEnc)
	filename := "./assets/tickets/" + strconv.Itoa(id) + ".json"
	fmt.Print(filename)
	err := ioutil.WriteFile(filename, encodedTicket, 0600)
	errorCheck(err)
}

func NewEntry(creator string, content string) Entry {
	date := time.Now().Local().Format("2006-01-02T15:04:05Z07:00")
	return Entry{Date: date, Creator: creator, Content: content}
}

func createTicketId() int {
	numberOfTickets := countTickets()
	if numberOfTickets < 0 {
		errorCheck(errors.New("unexpected error: number of tickets < 0"))
	}
	return numberOfTickets + 1
}

func countTickets() int {
	tickets, err := ioutil.ReadDir("./assets/tickets")
	errorCheck(err)
	return len(tickets)
}
