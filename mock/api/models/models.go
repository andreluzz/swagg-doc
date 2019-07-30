package models

import (
	"time"

	"github.com/andreluzz/swagg-doc/mock/shared"
)

// User swagg-doc:model
// name swagg-doc:attribute:string
// id swagg-doc:attribute:ignore_write
// created_date swagg-doc:attribute:ignore_write
// Defines a user model
type User struct {
	ID          int                `json:"id" pk:"true"`
	Code        string             `json:"code" updatable:"false" validate:"required"`
	Name        shared.Translation `json:"name" validate:"required"`
	Addr        Address            `json:"address"`
	Resource    shared.Resource    `json:"resource"`
	CreatedDate time.Time          `json:"created_date"`
}

// Address swagg-doc:model
type Address struct {
	Street string `json:"street" validate:"required"`
	Number int    `json:"number"`
	City   string `json:"city"`
}
