package ticket

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type Status string

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

const (
	Open      Status = "offen"
	InProcess Status = "in Bearbeitung"
	Closed    Status = "geschlossen"
)

var ticket Ticket
var entry Entry
var orderedTickets []Ticket

//Error Handler
func errorCheck(err error) {
	if err != nil {
		panic(err)
	}
}

//Zählt, wie viele Tickets sich im Ordner "../assets/tickets" befinden
func countTickets() int {
	tickets, err := ioutil.ReadDir("../assets/tickets")
	errorCheck(err)
	return len(tickets)
}

//A-5:
//Ticketerstellung, Erfassung der Eingabedaten
//Ticket wird in die globale Variable "ticket" geladen und anschließend in einer .json Datei gespeichert
func NewTicket(subject string, creator string, content string) {
	id := NewId()
	ticket = Ticket{Id: id, Subject: subject, Status: Open, EditorId: 0, Entries: []Entry{*NewEntry(creator, content)}}
	writeTicket(&ticket)
}

//Erstellt einen neuen Eintrag und gibt die Referenz auf ihn zurück
func NewEntry(creator string, content string) *Entry {
	date := time.Now().Local().Format("2006-01-02T15:04:05Z07:00")
	entry = Entry{Date: date, Creator: creator, Content: content}
	return &entry
}

//Liest das Ticket mit der ID "id" aus der entsprechenden .json Datei
//und gibt die Referenz auf das entsprechende Ticket zurück
func readTicket(id int) *Ticket {
	filename := "../assets/tickets/" + strconv.Itoa(id) + ".json"
	encodedTicket, errRead := ioutil.ReadFile(filename)
	errorCheck(errRead)
	err := json.Unmarshal(encodedTicket, &ticket)
	errorCheck(err)
	return &ticket
}

//Schreibt das Ticket "ticket" in die entsprechende .json Datei oder erzeugt diese
func writeTicket(ticket *Ticket) {
	encodedTicket, errEnc := json.Marshal(ticket)
	errorCheck(errEnc)
	filename := "../assets/tickets/" + strconv.Itoa(ticket.Id) + ".json"
	err := ioutil.WriteFile(filename, encodedTicket, 0600)
	errorCheck(err)
}

//Funktion: NeN Hinzufügen von Einträgen
func AppendEntry(id int, creator string, content string) {
	readTicket(id)
	ticket.Entries = append(ticket.Entries, *NewEntry(creator, content))
	writeTicket(&ticket)
}

//A-8.2:
//Bearbeitung eines Tickets, Ticket nehmen
//Das entsprechende Ticket wird gelesen, die Werte Status und EditorId überschrieben
//und das Ticket in die Datei zurückgeschrieben
func TakeTicket(id int, editorId int) {
	readTicket(id)
	ticket.Status = InProcess
	ticket.EditorId = editorId
	writeTicket(&ticket)
}

//TODO Input auf nicht erlaubte Strings überprüfen
//A-8.3
//Bearbeitung eines Tickets, Alle (offenen) Tickets einsehen
//Übergben werden können null bis zwei Status und eine Referenz auf die Tickets mit entsprechendem Status
//werden zurückgeliefert.
//Bei null Argumenten wird auf alle Tickets referenziert
func GetTickets(status ...Status) *[]Ticket {

	//Überprüfung ob die Eingabe gültig ist
	if len(status) > 2 {
		errorCheck(errors.New("invalid input by getTickets"))
	}

	orderedTickets = []Ticket{}
	files, err := ioutil.ReadDir("../assets/tickets")
	errorCheck(err)
	var id int

	//Es wird durch alle gelesenen Tickets durchiteriert
	for _, file := range files {
		id = parseFilename(file.Name())
		readTicket(id)

		//Überprüfung ob alle oder nur bestimmte Tickets zurückgeliefert werden
		if len(status) > 0 {

			//Für jeden Status der übergeben wird
			for _, state := range status {
				if ticket.Status == state {
					orderedTickets = append(orderedTickets, ticket)
				}
			}
		} else {
			orderedTickets = append(orderedTickets, ticket)
		}
	}

	return &orderedTickets
}

//Parst aus dem Dateinamen eines Tickets die entsprechende ID raus
func parseFilename(filename string) int {
	filename = strings.TrimSuffix(filename, ".json")
	id, err := strconv.Atoi(filename)
	errorCheck(err)
	return id
}

//A-8.4:
//Bearbeitung eines Tickets, Tickets nach Übernahme wieder freigeben
//Das entsprechende Ticket wird gelesen, die Werte Status und EditorId überschrieben
//und das Ticket in die Datei zurückgeschrieben
func UnhandTicket(id int) {
	readTicket(id)
	ticket.Status = Open
	ticket.EditorId = 0
	writeTicket(&ticket)
}

//A-8.5:
//Bearbeitung eines Tickets, Ticket jmd anderem zuteilen
//Das entsprechende Ticket wird gelesen, die Werte Status und EditorId überschrieben
//und das Ticket in die Datei zurückgeschrieben
func DelegateTicket(id int, editorId int) {
	readTicket(id)
	ticket.Status = InProcess
	ticket.EditorId = editorId
	writeTicket(&ticket)
}

//A-12:
//Zusammenführen von Tickets
//Die ID's der Tickets, welche zusammengeführt werden, werden übergeben.
//Das erste Argument ist das Ticket, in welches geschrieben wird, das Zweite das Ticket welches gelöcht wird.
func MergeTickets(dest int, source int) {
	destTicket := *readTicket(dest)
	sourceTicket := *readTicket(source)

	//Die Einträge des "sourceTickets" werden an die Einträge des "destTickets" gehangen
	for _, entry := range sourceTicket.Entries {
		destTicket.Entries = append(destTicket.Entries, entry)
	}
	writeTicket(&destTicket)

	DeleteTicket(source)
}

//Löscht die .json Datei des angegebenen Tickets
func DeleteTicket(id int) {
	filename := "../assets/tickets/" + strconv.Itoa(id) + ".json"
	err := os.Remove(filename)
	errorCheck(err)
}

//Liefert eine Referenz auf das angegebene Ticket
func GetTicket(id int) *Ticket {
	readTicket(id)
	return &ticket
}
