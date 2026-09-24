package auth

import (
	"github.com/sushanthach12/ecom-go/internal/constants"
	"github.com/sushanthach12/ecom-go/internal/helpers"
)

type registerPayloadDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto *registerPayloadDto) Validate() error {
	if helpers.CheckIfStringEmpty(dto.Name) {
		return &constants.ValidationError{
			Field:   "name",
			Message: "must not be empty",
		}
	}

	if helpers.CheckIfStringEmpty(dto.Email) {
		return &constants.ValidationError{
			Field:   "email",
			Message: "must not be empty",
		}
	}

	if !helpers.CheckIfValidEmail(dto.Email) {
		return &constants.ValidationError{
			Field:   "email",
			Message: "must be a valid email address",
		}
	}

	if helpers.CheckIfStringEmpty(dto.Password) {
		return &constants.ValidationError{
			Field:   "password",
			Message: "must not be empty",
		}
	}

	if !helpers.CheckIfStrongPassword(dto.Password) {
		return &constants.ValidationError{
			Field:   "password",
			Message: "must be at least 8 characters and include an uppercase letter, a number, and a special character",
		}
	}

	return nil
}

type loginPayloadDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto *loginPayloadDto) Validate() error {
	if helpers.CheckIfStringEmpty(dto.Email) {
		return &constants.ValidationError{
			Field:   "email",
			Message: "must not be empty",
		}
	}

	if !helpers.CheckIfValidEmail(dto.Email) {
		return &constants.ValidationError{
			Field:   "email",
			Message: "must be a valid email address",
		}
	}

	if helpers.CheckIfStringEmpty(dto.Password) {
		return &constants.ValidationError{
			Field:   "password",
			Message: "must not be empty",
		}
	}

	if !helpers.CheckIfStrongPassword(dto.Password) {
		return &constants.ValidationError{
			Field:   "password",
			Message: "must be at least 8 characters and include an uppercase letter, a number, and a special character",
		}
	}

	return nil
}
