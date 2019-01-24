package models

import (
	"github.com/gobuffalo/uuid"
)

type Category interface {
	GetID() uuid.UUID
	GetName() string
	GetBackground() string
	GetIndex() int
}

type Categories []Category
