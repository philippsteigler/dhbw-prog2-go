package sessionHandler

import (
	"testing"
)

func TestEncryptString(t *testing.T) {
	data := []struct {
		input string
		key   string
	}{
		{"TEST", "123"},
		{"code", "review"},
		{"Ticket", ""},
		{"", "System"},
		{"Long input with more than 16 characters", "BETA"},
	}

	for _, d := range data {
		_, err := encryptString(d.input, d.key)
		if err != nil {
			t.Errorf("Unable to encrypt '%v' with key '%v': %v", d.input, d.key, err)
			continue
		}
	}
}

func TestDecryptString(t *testing.T) {
	data := []struct {
		input string
		key   string
	}{
		{"TEST", "123"},
		{"code", "review"},
		{"Ticket", ""},
		{"", "System"},
		{"Long input with more than 16 characters", "BETA"},
	}

	for _, d := range data {
		enc, err := encryptString(d.input, d.key)
		if err != nil {
			t.Errorf("Unable to encrypt '%v' with key '%v': %v", d.input, d.key, err)
			continue
		}
		dec, err := decryptString(enc, d.key)
		if err != nil {
			t.Errorf("Unable to decrypt '%v' with key '%v': %v", enc, d.key, err)
			continue
		}
		if dec != d.input {
			t.Errorf("Decrypt Key %v\n  Input: %v\n  Expected: %v\n  Actual: %v", d.key, enc, d.input, enc)
		}
	}
}
