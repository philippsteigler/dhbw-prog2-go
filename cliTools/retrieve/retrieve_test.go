package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Matrikelnummern:
//
// 3333958
// 3880065
// 8701350

func TestSendGetRequest(t *testing.T) {

	response, err := sendGetRequest()

	assert.NoError(t, err)
	assert.Equal(t, "200 OK", response.Status)
}
