package ticket

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Status string

type Mail struct {
	Email   string
	Subject string
	Content string
}

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

//Parst aus dem Dateinamen eines Tickets die entsprechende ID raus
func parseFilename(filename string) int {
	filename = strings.TrimSuffix(filename, ".json")
	id, err := strconv.Atoi(filename)
	sessionHandler.HandleError(err)
	return id
}

//Liest das gewünschte Ticket aus der JSON-Datei
//und gibt eine Referenz auf das Ticket zurück
func readTicket(id int) *Ticket {
	var readTicket Ticket
	filename := sessionHandler.GetAssetsDir() + "tickets/" + strconv.Itoa(id) + ".json"
	encodedTicket, errRead := ioutil.ReadFile(filename)
	sessionHandler.HandleError(errRead)
	err := json.Unmarshal(encodedTicket, &readTicket)
	sessionHandler.HandleError(err)
	return &readTicket
}

//Schreibt ein Ticket in seine entsprechende JSON-Datei oder erzeugt diese
func writeTicket(ticket *Ticket) {

	filename := sessionHandler.GetAssetsDir() + "tickets/" + strconv.Itoa((*ticket).Id) + ".json"

	encodedTicket, err := json.Marshal(*ticket)
	sessionHandler.HandleError(err)

	err = ioutil.WriteFile(filename, encodedTicket, 0600)
	sessionHandler.HandleError(err)
}

//Schreibt eine Mail in eine JSON-Datei
func writeMail(mail Mail) {
	filename := sessionHandler.GetAssetsDir() + "mails/" + strconv.Itoa(newId("/mails")) + ".json"

	encodedMail, err := json.Marshal(mail)
	sessionHandler.HandleError(err)

	err = ioutil.WriteFile(filename, encodedMail, 0600)
	sessionHandler.HandleError(err)
}

//Liefert eine Referenz auf das angegebene Ticket
func GetTicket(id int) *Ticket {
	return readTicket(id)
}

//Zur Erzeugung einer TicketID wird die höchste vergebene ID inkrementiert
func newId(path string) int {
	var ids []int

	ticketDir := sessionHandler.GetAssetsDir() + path
	files, err := ioutil.ReadDir(ticketDir)
	sessionHandler.HandleError(err)

	//Falls keine Tickets existieren
	if len(files) == 0 {
		return 1
	}

	//Jede ID aus dem Dateinamen parsen und in ids speichern
	for _, file := range files {

		indexOfFileExtension := strings.Index(file.Name(), ".")
		fileId, err := strconv.Atoi(file.Name()[:indexOfFileExtension])
		sessionHandler.HandleError(err)
		ids = append(ids, fileId)
	}

	//ids sortieren (aufsteigend) und die höchste vergebene ID inkrementieren und zurückgeben
	sort.Ints(ids)

	return ids[len(ids)-1] + 1
}

//A-5:
//Ticketerstellung, Erfassung der Eingabedaten
func NewTicket(subject string, creator string, content string) {
	newTicket := Ticket{Id: newId("/tickets"), Subject: subject, Status: Open, EditorId: 0, Entries: []Entry{NewEntry(creator, content)}}
	writeTicket(&newTicket)
}

//Erstellt einen neuen Eintrag
func NewEntry(creator string, content string) Entry {
	date := time.Now().Local().Format("2006-01-02T15:04:05.0000")
	return Entry{Date: date, Creator: creator, Content: content}
}

//Fügt einen neuen Eintrag einem bestehenden Ticket hinzu
//Wenn mail true ist, wird eine E-Mail erzeugt und als JSON gespeichert
func AppendEntry(id int, creator, content string, mail bool) {
	ticketToAppend := *GetTicket(id)
	ticketToAppend.Entries = append(ticketToAppend.Entries, NewEntry(creator, content))
	writeTicket(&ticketToAppend)
	if mail {
		//Email und Subject aus Ticket laden
		storedTicket := readTicket(id)
		subject := "Rueckmeldung bzgl.:" + storedTicket.Subject
		writeMail(Mail{Email: storedTicket.Entries[0].Creator, Subject: subject, Content: content})
	}
}

//Tickets nach einer bestimmten EditorID filtern und zurückgeben
func GetTicketsByEditorId(editorId int) *[]Ticket {

	var orderedTickets []Ticket

	files, err := ioutil.ReadDir(sessionHandler.GetAssetsDir() + "tickets")
	sessionHandler.HandleError(err)

	for _, file := range files {
		actualTicket := GetTicket(parseFilename(file.Name()))

		if actualTicket.EditorId == editorId {
			orderedTickets = append(orderedTickets, *actualTicket)
		}
	}

	return &orderedTickets
}

//A-8.2:
//Bearbeitung eines Tickets, Ticket nehmen
func TakeTicket(id int, editorId int) {
	ticketToTake := GetTicket(id)
	ticketToTake.Status = InProcess
	ticketToTake.EditorId = editorId
	writeTicket(ticketToTake)
}

//A-8.3
//Bearbeitung eines Tickets, alle offenen Tickets einsehen
func GetAllOpenTickets() *[]Ticket {
	var orderedTickets []Ticket

	files, err := ioutil.ReadDir(sessionHandler.GetAssetsDir() + "tickets")
	sessionHandler.HandleError(err)

	for _, file := range files {
		actualTicket := GetTicket(parseFilename(file.Name()))

		if actualTicket.Status == Open {
			orderedTickets = append(orderedTickets, *actualTicket)
		}
	}

	return &orderedTickets

}

//A-8.4:
//Bearbeitung eines Tickets, Tickets nach Übernahme wieder freigeben
func UnhandTicket(id int) {
	ticketToUnhand := GetTicket(id)
	ticketToUnhand.Status = Open
	ticketToUnhand.EditorId = 0
	writeTicket(ticketToUnhand)
}

//A-8.5:
//Bearbeitung eines Tickets, Ticket jmd anderem zuteilen
func DelegateTicket(id int, editorId int) {
	ticketToDelegate := GetTicket(id)
	ticketToDelegate.Status = InProcess
	ticketToDelegate.EditorId = editorId
	writeTicket(ticketToDelegate)
}

//A-12:
//Zusammenführen von Tickets
//Die ID's der Tickets, welche zusammengeführt werden, werden übergeben.
//Das erste Argument ist das Ticket, in welches geschrieben wird, das Zweite das Ticket welches gelöcht wird.
func MergeTickets(dest int, source int) {
	destTicket := GetTicket(dest)
	sourceTicket := GetTicket(source)

	//Die Einträge des "sourceTickets" werden an die Einträge des "destTickets" gehangen
	for _, entry := range sourceTicket.Entries {
		destTicket.Entries = append(destTicket.Entries, entry)
	}

	writeTicket(destTicket)
	deleteTicket(source)
}

//Löscht die JSON-Datei des angegebenen Tickets
func deleteTicket(id int) {
	filename := sessionHandler.GetAssetsDir() + "tickets/" + strconv.Itoa(id) + ".json"
	err := os.Remove(filename)
	sessionHandler.HandleError(err)
}

//Liefert alle Tickets zurück, die existieren
func GetAllTickets() *[]Ticket {
	var orderedTickets []Ticket

	files, err := ioutil.ReadDir(sessionHandler.GetAssetsDir() + "/tickets")
	sessionHandler.HandleError(err)

	for _, file := range files {
		orderedTickets = append(orderedTickets, *GetTicket(parseFilename(file.Name())))
	}

	return &orderedTickets
}

//Setzt geschlossene Tickets auf offen
func SetTicketToOpenIfClosed(id int) {
	ticketToOpen := GetTicket(id)
	if ticketToOpen.Status == Closed {
		ticketToOpen.Status = Open
		writeTicket(ticketToOpen)
	}
}

//Ticket schließen
func CloseTicket(id int) {
	ticketToClose := GetTicket(id)
	ticketToClose.Status = Closed
	ticketToClose.EditorId = 0
	writeTicket(ticketToClose)
}

//Liefert alle eine Referenz auf alle Mails
func GetAllMails() *[]Mail {
	var orderedMails []Mail
	var mail Mail

	files, err := ioutil.ReadDir(sessionHandler.GetAssetsDir() + "/mails")
	sessionHandler.HandleError(err)

	for _, file := range files {
		filename := sessionHandler.GetAssetsDir() + "mails/" + file.Name()

		encodedTicket, errRead := ioutil.ReadFile(filename)
		sessionHandler.HandleError(errRead)
		err := json.Unmarshal(encodedTicket, &mail)
		sessionHandler.HandleError(err)

		orderedMails = append(orderedMails, mail)
	}

	return &orderedMails
}

//Löscht eine Mail. Dazu wird eine Mail übergeben und mit allen abgeglichen. Wenn die übergebene Mail existiert, wird sie gelöscht
func DeleteMail(mail Mail) {
	var storedMail Mail
	files, err := ioutil.ReadDir(sessionHandler.GetAssetsDir() + "/mails")
	sessionHandler.HandleError(err)
	for _, file := range files {
		filename := sessionHandler.GetAssetsDir() + "mails/" + file.Name()

		encodedTicket, errRead := ioutil.ReadFile(filename)
		sessionHandler.HandleError(errRead)
		err := json.Unmarshal(encodedTicket, &storedMail)
		sessionHandler.HandleError(err)

		if mail.Content == storedMail.Content && mail.Email == storedMail.Email && mail.Subject == storedMail.Subject {
			err := os.Remove(filename)
			sessionHandler.HandleError(err)
		}
	}
}
