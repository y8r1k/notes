package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrInvalidURLID = errors.New("invalid ID in url")
)

func parseID(r *http.Request) (int64, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "notes" {
		return 0, fmt.Errorf("%w: %s", ErrInvalidURLID, r.URL.Path)
	}

	idStr := parts[2]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidURLID, idStr)
	}

	return id, nil
}
