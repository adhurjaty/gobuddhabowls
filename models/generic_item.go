package models

import (
	"github.com/gobuffalo/uuid"
)

type Model interface {
	GetID() uuid.UUID
}

// GenericItem is an interface for any item whose lifecycle includes inventory
type GenericItem interface {
	Model
	GetInventoryItemID() uuid.UUID
	GetBaseItem() GenericItem
	SetBaseItem(GenericItem)
	GetName() string
	GetCategory() ItemCategory
	GetCountUnit() string
	GetIndex() int
}

type GenericItems interface {
	ToGenericItems() *[]GenericItem
}
