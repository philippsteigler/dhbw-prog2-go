package pageHandler

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewTicketViewPageHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://localhost:8000/", nil)

	NewTicketViewPageHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Page didn't return %v, %v", http.StatusOK, response.Code)
	}
}

func TestDashboardViewPageHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://localhost:8000/dashboard", nil)

	DashboardViewPageHandler(response, request)

	if response.Code != http.StatusFound {
		t.Errorf("Page didn't return %v, %v", http.StatusFound, response.Code)
	}
}

func TestTicketAppendEntry(t *testing.T) {
	setup()
	defer teardown()

	//Erzeugen einer Testmail
	mail := map[string]string{"email": "test@home.com", "subject": "CreateNewTicket Test", "content": "Ein weiterer Test."}
	jsonMail, err := json.Marshal(mail)
	sessionHandler.HandleError(err)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "https://localhost:8000/ticket", bytes.NewBuffer(jsonMail))
	CreateNewTicket(response, request)

	assert.Equal(t, http.StatusOK, response.Code)

	newTicket := ticket.GetTicket(3)

	assert.Equal(t, "test@home.com", newTicket.Entries[0].Creator)
	assert.Equal(t, "CreateNewTicket Test", newTicket.Subject)
	assert.Equal(t, "Ein weiterer Test.", newTicket.Entries[0].Content)
}