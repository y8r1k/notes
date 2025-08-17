package domain

import (
	"time"
)

type Note struct {
	ID          int64
	Title       string
	Description string
	Created     time.Time
	Changed     *time.Time
}
