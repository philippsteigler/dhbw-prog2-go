package ticket

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type Status string

const (
	Open      Status = "offen"
	InProcess Status = "in Bearbeitung"
	Closed    Status = "geschlossen"
)

var ticket Ticket

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

func NewTicket(subject string, creator string, content string) {
	id := NewId()
	ticket = Ticket{Id: id, Subject: subject, Status: Open, EditorId: 0, Entries: []Entry{NewEntry(creator, content)}}
	writeTicket(ticket)
}

//TODO pointer übergeben
func NewEntry(creator string, content string) Entry {
	date := time.Now().Local().Format("2006-01-02T15:04:05Z07:00")
	return Entry{Date: date, Creator: creator, Content: content}
}

func countTickets() int {
	tickets, err := ioutil.ReadDir("./assets/tickets")
	errorCheck(err)
	return len(tickets)
}

//Erw. Funktion: Für xxx Json wird eingelesen
func readTicket(id int) Ticket {
	filename := "./assets/tickets/" + strconv.Itoa(id) + ".json"
	encodedTicket, errRead := ioutil.ReadFile(filename)
	errorCheck(errRead)
	err := json.Unmarshal(encodedTicket, &ticket)
	errorCheck(err)
	return ticket
}

func writeTicket(ticket Ticket) {
	encodedTicket, errEnc := json.Marshal(ticket)
	errorCheck(errEnc)
	filename := "./assets/tickets/" + strconv.Itoa(ticket.Id) + ".json"
	err := ioutil.WriteFile(filename, encodedTicket, 0600)
	errorCheck(err)
}

//Funktion: NeN Hinzufügen von Einträgen
func AppendEntry(id int, creator string, content string) {
	ticket = readTicket(id)
	ticket.Entries = append(ticket.Entries, NewEntry(creator, content))
	writeTicket(ticket)
}

//Funktion: 8.2 Bearbeiter soll ein Ticket übernehmen können
func TakeTicket(id int, userId int) {
	ticket = readTicket(id)
	ticket.Status = InProcess
	ticket.EditorId = userId
	writeTicket(ticket)
}

//TODO
//Funktion: 8.3 Bearbeiter soll alle offenen Tickets sehen können
func GetOpenTickets() (openTickets []int) {
	files, err := ioutil.ReadDir("./assets/tickets")
	errorCheck(err)
	var id int
	for _, file := range files {
		id = parseFilename(file.Name())
		ticket = readTicket(id)
		if ticket.Status == Open {
			openTickets = append(openTickets, ticket.Id)
		}
	}
	return openTickets
}

//Erw. Funktion: für 8.3 zum Auslesen der Ticketnummer aus dem Dateinamen
func parseFilename(filename string) int {
	filename = strings.TrimSuffix(filename, ".json")
	id, err := strconv.Atoi(filename)
	errorCheck(err)
	return id
}

//Funktion: 8.4 Tickets nach Übernahme freigeben
func UnhandTicket(id int) {
	ticket = readTicket(id)
	ticket.Status = Open
	ticket.EditorId = 0
	writeTicket(ticket)
}

//Funktion: 8.5 Ticket jmd anderem zuteilen
func delegateTicket(id int, editorId int) {
	ticket = readTicket(id)
	ticket.Status = InProcess
	ticket.EditorId = editorId
	writeTicket(ticket)
}

//Funktion: 12 Zusammenführen von Tickets
func mergeTickets(dest int, source int) {
	ticket = readTicket(dest)
	sourceTicket := readTicket(source)
	for _, entry := range sourceTicket.Entries {
		ticket.Entries = append(ticket.Entries, entry)
	}
	writeTicket(ticket)
	DeleteTicket(source)
}

//Erw. Funktion: Für 12 Löschen von Tickets
func DeleteTicket(id int) {
	filename := "./assets/tickets/" + strconv.Itoa(id) + ".json"
	err := os.Remove(filename)
	errorCheck(err)
}
