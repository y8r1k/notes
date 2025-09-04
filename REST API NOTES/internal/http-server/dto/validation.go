package dto

import (
	"errors"
)

var (
	ErrInvalidTitle       = errors.New("title must be 1..20 characters")
	ErrInvalidDescription = errors.New("the description cannot be empty")
)

func (r CreateNotePOSTRequest) Validate() error {
	if len(r.Title) == 0 || len(r.Title) > 20 {
		return ErrInvalidTitle
	} else if len(r.Description) == 0 {
		return ErrInvalidDescription
	}
	return nil
}

func (r ChangeNotePUTRequest) Validate() error {
	if len(r.Title) == 0 || len(r.Title) > 20 {
		return ErrInvalidTitle
	} else if len(r.Description) == 0 {
		return ErrInvalidDescription
	}
	return nil
}
