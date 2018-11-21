package sessionHandler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAssetsDir(t *testing.T) {
	path := GetAssetsDir()
	assert.Equal(t, "../assets/", path, "Wrong path to folder 'assets'.")
}

func TestLoadUserData(t *testing.T) {
	var users UserAccounts
	assert.Empty(t, &users, "User rollback should be empty at first.")

	users = *loadUserData()
	assert.NotEmpty(t, &users, "User rollback should be available after reading from storage.")
}

// TODO: Tests für Session-Handling

func TestLoginHandler(t *testing.T) {
	// Teste Eingaben für autorisierten Benutzer.
	req1, _ := http.NewRequest(http.MethodPost, "/login", nil)
	req1.ParseForm()
	req1.Form.Add("username", "admin")
	req1.Form.Add("password", "test123")
	res1 := httptest.NewRecorder()

	LoginHandler(res1, req1)
	path1 := string(res1.Header().Get("Location"))
	assert.Equal(t, "/internal", path1, "Redirect target should be '/internal' after successful login.")

	// Teste falsche Nutzerdaten.
	req1, _ = http.NewRequest(http.MethodPost, "/login", nil)
	req1.ParseForm()
	req1.Form.Add("username", "randumUser")
	req1.Form.Add("password", "IdOnTkNoW")
	res1 = httptest.NewRecorder()

	LoginHandler(res1, req1)
	path1 = string(res1.Header().Get("Location"))
	assert.Equal(t, "/", path1, "Redirect should be '/' after successful login.")
}
