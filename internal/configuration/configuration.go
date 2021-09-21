package configuration

import (
	errs "app/internal/errors"
	"github.com/google/uuid"
)

type UserIdentifierRequest struct {
	UserId string `json:"user_id"`
}

func (usr *UserIdentifierRequest) Validate() []errs.ValidationError {
	var vErrs []errs.ValidationError

	_, err := uuid.Parse(usr.UserId)
	if err != nil {
		vErrs = append(vErrs, errs.ValidationError{Path: "user_id", Message: err.Error()})
	}

	return vErrs
}

type UserRequest struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
}

func (usr *UserRequest) Validate() []errs.ValidationError {
	var vErrs []errs.ValidationError

	_, err := uuid.Parse(usr.UserId)
	if err != nil {
		vErrs = append(vErrs, errs.ValidationError{Path: "user_id", Message: err.Error()})
	}

	if len(usr.Name) < 4 || len(usr.Name) > 100 {
		vErrs = append(vErrs, errs.ValidationError{Path: "name", Message: "invalid length exceeded - (4-100)"})
	}

	if usr.Age < 1 {
		vErrs = append(vErrs, errs.ValidationError{Path: "age", Message: "invalid min value exceeded - (1)"})
	}

	return vErrs
}
