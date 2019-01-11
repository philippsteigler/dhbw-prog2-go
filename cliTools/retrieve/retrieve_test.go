package main

import (
	"de/vorlesung/projekt/crew/pageHandler"
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

func TestSendGetRequest(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(pageHandler.Mails))
	defer ts.Close()

	client := ts.Client()

	response, err := client.Get(ts.URL)

	assert.NoError(t, err)
	assert.Equal(t, "200 OK", response.Status)
}
