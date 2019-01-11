package pageHandler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Matrikelnummern:
//
// 3333958
// 3880065
// 8701350

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
	request := httptest.NewRequest("GET", "https://localhost:4443/dashboard", nil)

	DashboardViewPageHandler(response, request)

	if response.Code != 302 {
		t.Errorf("Page didn't redirect %v, %v", 302, response.Code)
	}
}

func TestLoginPageHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://localhost:4443/loginView", nil)

	LoginPageHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Page didn't return %v, %v", http.StatusOK, response.Code)
	}
}

func TestNewTicketViewPageHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://localhost:4443/", nil)

	NewTicketViewInit()
	NewTicketViewPageHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Page didn't return %v, %v", http.StatusOK, response.Code)
	}
}

func TestTicketInsightPageHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://localhost:4443/ticketInsightView", nil)

	TicketInsightPageHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Page didn't return %v, %v", http.StatusOK, response.Code)
	}
}

func TestTicketsViewPageHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "https://localhost:4443/ticketsView", nil)

	TicketsViewPageHandler(response, request)

	if response.Code != 302 {
		t.Errorf("Page didn't redirect %v, %v", 302, response.Code)
	}
}
