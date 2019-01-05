package pageHandler

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"de/vorlesung/projekt/crew/ticket"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

//Hilfsfunktion zum überprüfen der Mails
func validateMail(mail ticket.Mail) error {
	if mail.Email == "" {
		return fmt.Errorf("Email is missing.")
	}
	if mail.Subject == "" {
		return fmt.Errorf("Subject is missing.")
	}
	if mail.Content == "" {
		return fmt.Errorf("Content is missing.")
	}
	return nil
}

//A-6:
//Email Empfang über eine REST-API
func CreateNewTicket(w http.ResponseWriter, r *http.Request) {

	//Nur POST Requests werden bearbeitet
	if r.Method == http.MethodPost {
		var newMail ticket.Mail

		//Request wird decodiert
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&newMail)
		sessionHandler.HandleError(err)

		//Mail wird überprüft
		validateMail(newMail)

		//Überprüfung, ob sich die Mail auf ein existierendes Ticket bezieht
		ticketExist, id := refersToExistingTicket(newMail.Subject, newMail.Email)

		if !ticketExist {
			//wenn nicht, wird ein neues Ticket erzeugt
			ticket.NewTicket(newMail.Subject, newMail.Email, newMail.Content)
		} else {
			//sonst wird ein neuer Eintrag an das bestehende Ticket angefügt
			ticket.AppendEntry(id, newMail.Email, newMail.Content, false)
			ticket.SetTicketToOpen(id)
		}
	}
}

//Eine Mail bezieht sich auf ein existierendes Ticket,
//wenn der Betreff in folgender Form ist: "RE: <Betreff des Tickets auf das sich bezogen wird>"
//sowie die Email des Senders mit der Email des Erstellers des Tickets übereinstimmt
//"RE:" am Anfang des Betreffs ist verpflichtend (Groß-/Kleinschreibung ist egal)
func refersToExistingTicket(subject, email string) (bool, int) {

	//Groß- und Kleinschreibung wird ignoriert
	if strings.ToLower(subject[:3]) != "re:" {
		return false, 0
	} else {
		//Alle Tickets werden geladen
		tickets := *ticket.GetAllTickets()

		//Betreff und Emailadresse werden mit jedem Ticket abgeglichen
		for _, t := range tickets {
			if strings.ToLower(t.Subject) == strings.ToLower(parseSubject(subject)) && t.Entries[0].Creator == email {
				return true, t.Id
			}
		}
	}
	return false, 0
}

//"RE:" sowie führende und abschließende Leerzeichen werden aus dem Betreff entfernt entfernt
func parseSubject(subject string) string {
	return strings.TrimSpace(subject[3:])
}

//A-7:
//Email Versand über eine REST-API

func Mails(w http.ResponseWriter, r *http.Request) {

	//Überprüfung ob ein GET oder POST Request auf die .../mails URL vorliegt
	if r.Method == http.MethodGet {
		retrieveMails(w)

	} else if r.Method == http.MethodPost {
		sentMails(r)
	}
}

//A-7.1: Funktion zum Abfrufen von Mails
func retrieveMails(w http.ResponseWriter) {

	//Alle zu versendenden Mails abrufen
	mails := *ticket.GetAllMails()

	//Als JSON codieren
	jsonMails, err := json.Marshal(mails)
	sessionHandler.HandleError(err)

	//Mails als JSON zurückgeben
	w.Write(jsonMails)
}

//A-7.2: Funktion zum Mitteilen, welche Mails versand wurden
func sentMails(r *http.Request) {

	//Gesendete Mails werden als JSON mit einem Array aus Mails erwartet
	var mails []ticket.Mail

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mails)
	sessionHandler.HandleError(err)

	//Jede Mail wird geprüft und anschließend gelöscht
	for _, mail := range mails {
		validateMail(mail)
		ticket.DeleteMail(mail)
	}
}
