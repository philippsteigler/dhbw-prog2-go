package sessionHandler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Matrikelnummern:
//
// 3333958
// 3880065
// 8701350

func setup() {
	BackupEnvironment()
	DemoMode()
	time.Sleep(100 * time.Millisecond)
}

func teardown() {
	RestoreEnvironment()
	time.Sleep(100 * time.Millisecond)
}

func TestGetAssetsDir(t *testing.T) {
	path := GetAssetsDir()
	assert.Equal(t, "../assets/", path, "Wrong path to folder 'assets'.")
}

func TestLoadUserData(t *testing.T) {
	setup()

	var users UserAccounts
	assert.Empty(t, &users, "User rollback should be empty at first.")

	users = *LoadUserData()
	assert.NotEmpty(t, &users, "User rollback should be available after reading from storage.")

	teardown()
}

func TestLoginHandler(t *testing.T) {
	setup()

	// Teste Eingaben für autorisierten Benutzer.
	req1, _ := http.NewRequest(http.MethodPost, "/login", nil)
	req1.ParseForm()
	req1.Form.Add("username", "admin")
	req1.Form.Add("password", "test123")
	res1 := httptest.NewRecorder()

	LoginHandler(res1, req1)
	path1 := string(res1.Header().Get("Location"))
	assert.Equal(t, "/dashboard", path1, "Redirect target should be '/dashboard' after successful login.")

	// Teste falsche Nutzerdaten.
	req1, _ = http.NewRequest(http.MethodPost, "/login", nil)
	req1.ParseForm()
	req1.Form.Add("username", "randumUser")
	req1.Form.Add("password", "IdOnTkNoW")
	res1 = httptest.NewRecorder()

	LoginHandler(res1, req1)
	path1 = string(res1.Header().Get("Location"))
	assert.Equal(t, "/loginView", path1, "Redirect should be '/loginView' after failed login.")

	teardown()
}
