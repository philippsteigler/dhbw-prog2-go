package pageHandler

import (
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
