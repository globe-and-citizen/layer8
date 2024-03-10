package entities

import (
	utilities "github.com/globe-and-citizen/layer8-utils"
)

type AbstractUser struct {
	Username string `json:"username" validate:"nonzero"`
	Email    string `json:"email" validate:"nonzero"`
	Fname    string `json:"first_name" validate:"nonzero"`
	Lname    string `json:"last_name" validate:"nonzero"`
}

type User struct {
	ID               int64        `json:"id"`
	Password         string       `json:"password" validate:"nonzero"`
	PsedonymizedData AbstractUser `json:"pseudonymized_data"`
	AbstractUser
}

// Validate validates the user struct by checking if the required fields are
// present and if the email is valid.
func (u *User) Validate() error {
	// Check if the required fields are present
	err := utilities.ValidateRequiredFields(u)
	if err != nil {
		return err
	}

	// Check if the email is valid
	err = utilities.ValidateEmail(u.Email)
	if err != nil {
		return err
	}

	return nil
}
