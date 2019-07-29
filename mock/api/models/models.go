package models

import (
	"time"

	"github.com/andreluzz/swagg-doc/mock/shared"
)

// User swagg-doc:model
// Defines a user model
type User struct {
	ID          int             `json:"id"`
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Resource    shared.Resource `json:"resource"`
	CreatedDate time.Time       `json:"created_date"`
}
