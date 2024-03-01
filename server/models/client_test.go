package models

import (
	"testing"
)

func TestCreateClient(t *testing.T) {
	client := CreateClient("1234", "5678#1234", "Tester Client", "my_redirection_uri")
	if client.Name != "Tester Client" {
		t.Errorf("Failed to create client struct: wrong name.")
	}
	if client.ID != "1234" {
		t.Errorf("Failed to create client struct: wrong id.")
	}
	if client.Secret != "5678#1234" {
		t.Errorf("Failed to create client struct: wrong secret.")
	}
	if client.RedirectURI != "my_redirection_uri" {
		t.Errorf("Failed to create client struct: wrong redirection URL")
	}
}
