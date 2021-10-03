package main

import (
	"net/mail"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailDisplayName1(t *testing.T) {
	email := "\"Display Name\" <user@domain.com>"
	parsed, err := mail.ParseAddress(email)
	assert.Nil(t, err)
	assert.Equal(t, "\"Display Name\" <user@domain.com>", parsed.String())
	assert.Equal(t, "Display Name", parsed.Name)
	assert.Equal(t, "user@domain.com", parsed.Address)
}

func TestEmailDisplayName2(t *testing.T) {
	email := "Display Name <user@domain.com>"
	parsed, err := mail.ParseAddress(email)
	assert.Nil(t, err)
	assert.Equal(t, "\"Display Name\" <user@domain.com>", parsed.String())
	assert.Equal(t, "Display Name", parsed.Name)
	assert.Equal(t, "user@domain.com", parsed.Address)
}

func TestEmailDisplayName3(t *testing.T) {
	email := "user@domain.com"
	parsed, err := mail.ParseAddress(email)
	assert.Nil(t, err)
	assert.Equal(t, "<user@domain.com>", parsed.String())
	assert.Equal(t, "", parsed.Name)
	assert.Equal(t, "user@domain.com", parsed.Address)
}
