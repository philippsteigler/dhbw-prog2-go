package main

import (
	"de/vorlesung/projekt/crew/sessionHandler"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Matrikelnummern:
//
// 3333958
// 3880065
// 8701350

func setup() {
	sessionHandler.BackupEnvironment()
	sessionHandler.DemoMode()
}

func teardown() {
	sessionHandler.RestoreEnvironment()
}

func TestSendPostRequest(t *testing.T) {
	setup()
	defer teardown()

	response, err := sendPostRequest([]string{"mail@test.com", "subject", "content"})

	assert.NoError(t, err)
	assert.Equal(t, "200 OK", response.Status)
}
