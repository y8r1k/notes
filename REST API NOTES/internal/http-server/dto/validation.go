package dto

import (
	"errors"
)

var (
	ErrInvalidTitle = errors.New("title must be 1..20 characters")
)

func (r CreateNotePOSTRequest) Validate() error {
	if len(r.Title) == 0 || len(r.Title) > 20 {
		return ErrInvalidTitle
	}
	return nil
}

func (r ChangeNotePUTRequest) Validate() error {
	if len(r.Title) == 0 || len(r.Title) > 20 {
		return ErrInvalidTitle
	}
	return nil
}
