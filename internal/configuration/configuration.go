package configuration

import (
	"app/internal/response"
	"github.com/google/uuid"
)

type UserIdentifierRequest struct {
	UserId string `json:"user_id"`
}

func (usr *UserIdentifierRequest) Validate() []response.ValidationError {
	var vErrs []response.ValidationError

	_, err := uuid.Parse(usr.UserId)
	if err != nil {
		vErrs = append(vErrs, response.ValidationError{Path: "user_id", Message: err.Error()})
	}

	return vErrs
}

type UserRequest struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
}

func (usr *UserRequest) Validate() []response.ValidationError {
	var vErrs []response.ValidationError

	_, err := uuid.Parse(usr.UserId)
	if err != nil {
		vErrs = append(vErrs, response.ValidationError{Path: "user_id", Message: err.Error()})
	}

	if len(usr.Name) < 4 || len(usr.Name) > 100 {
		vErrs = append(vErrs, response.ValidationError{Path: "name", Message: "invalid length exceeded - (4-100)"})
	}

	if usr.Age < 1 {
		vErrs = append(vErrs, response.ValidationError{Path: "age", Message: "invalid min value exceeded - (1)"})
	}

	return vErrs
}
