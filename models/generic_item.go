package models

import (
	"github.com/gobuffalo/uuid"
)

// GenericItem is an interface for any item whose lifecycle includes inventory
type GenericItem interface {
	GetID() uuid.UUID
	GetName() string
	GetCategory() InventoryItemCategory
	GetIndex() int
}

type GenericItems []GenericItem
