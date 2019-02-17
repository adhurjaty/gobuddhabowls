package models

import (
	"github.com/gobuffalo/uuid"
)

// GenericItem is an interface for any item whose lifecycle includes inventory
type GenericItem interface {
	GetID() uuid.UUID
	GetInventoryItemID() uuid.UUID
	GetName() string
	GetCategory() ItemCategory
	GetCountUnit() string
	GetIndex() int
}

type GenericItems interface {
	ToGenericItems() *[]GenericItem
}
