package pageHandler

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"de/vorlesung/projekt/crew/ticket"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Mail struct {
	Email   string
	Subject string
	Content string
}

func CreateNewTicket(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newMail Mail
	err := decoder.Decode(&newMail)
	sessionHandler.HandleError(err)
	validateMail(newMail)

	ticketExist, id := refersToExistingTicket(newMail.Subject, newMail.Email)

	if !ticketExist {
		ticket.NewTicket(newMail.Subject, newMail.Email, newMail.Content)
	} else {
		ticket.AppendEntry(id, newMail.Email, newMail.Content)
		ticket.SetTicketToOpen(id)
	}

}

//Eine Mail bezieht sich auf ein existierendes Ticket,
//wenn der Betreff in folgender Form ist: "RE: <Betreff des Tickets auf das sich bezogen wird>"
//"RE:" am Anfang des Betreffs ist verpflichtend (Groß-/Kleinschreibung ist egal)
func refersToExistingTicket(subject, email string) (bool, int) {
	if strings.ToLower(subject[:3]) != "re:" {
		return false, 0
	} else {
		tickets := *ticket.GetAllTickets()

		for _, t := range tickets {
			if strings.ToLower(t.Subject) == strings.ToLower(parseSubject(subject)) && t.Entries[0].Creator == email {
				return true, t.Id
			}
		}
	}

	return false, 0
}

func parseSubject(subject string) string {
	return strings.TrimSpace(subject[3:])
}

func validateMail(mail Mail) error {
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
