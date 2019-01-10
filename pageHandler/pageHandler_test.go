package pageHandler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDashboardViewInit(t *testing.T) {
	assert.Empty(t, dashboardViewTemplates, "dashboardViewTemplate should be empty")
	DashboardViewInit()
	assert.NotEmpty(t, dashboardViewTemplates, "dashboardViewTemplates should not be empty")
}

func TestNewTicketViewInit(t *testing.T) {
	assert.Empty(t, newTicketViewTemplates, "newTickerViewTemplates should be empty")
	NewTicketViewInit()
	assert.NotEmpty(t, newTicketViewTemplates, "newTickerViewTemplates should not be empty")
}

func TestTicketInsightInit(t *testing.T) {
	assert.Empty(t, ticketInsightTemplates, "ticketInsightTemplates should be empty")
	TicketInsightInit()
	assert.NotEmpty(t, ticketInsightTemplates, "ticketInsightTemplates should not be empty")
}

func TestTicketsViewInit(t *testing.T) {
	assert.Empty(t, ticketsViewTemplates, "ticketViewTemplates should be empty")
	TicketsViewInit()
	assert.NotEmpty(t, ticketsViewTemplates, "ticketViewTemplates should not be empty")
}

func TestDashboardViewPageHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://localhost:8000/dashboard", nil)

	DashboardViewPageHandler(response, request)

	if response.Code != 302 {
		t.Errorf("Page didn't redirect %v, %v", 302, response.Code)
	}
}

func TestNewTicketViewPageHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://localhost:8000/", nil)

	NewTicketViewInit()
	NewTicketViewPageHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Page didn't return %v, %v", http.StatusOK, response.Code)
	}
}

func TestTicketInsightPageHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://localhost:8000/ticketInsightView", nil)

	TicketInsightPageHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Page didn't return %v, %v", http.StatusOK, response.Code)
	}
}

func TestTicketsViewPageHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://localhost:8000/ticketsView", nil)

	TicketsViewPageHandler(response, request)

	if response.Code != 302 {
		t.Errorf("Page didn't redirect %v, %v", 302, response.Code)
	}
}

/*
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
*/
