package models

import (
	"github.com/gobuffalo/uuid"
)

type Model interface {
	GetID() uuid.UUID
}

type Models interface {
	ToModels() *[]Model
}

// GenericItem is an interface for any item whose lifecycle includes inventory
type GenericItem interface {
	Model
	GetName() string
	GetCategory() ItemCategory
	GetCategoryID() uuid.UUID
	SetCategory(ItemCategory)
	GetCountUnit() string
	GetIndex() int
}

type GenericItems interface {
	ToGenericItems() *[]GenericItem
}

type CompoundItem interface {
	GenericItem
	GetBaseItemID() uuid.UUID
	GetBaseItem() GenericItem
	SetBaseItem(GenericItem)
}

type CompoundItems interface {
	ToCompoundItems() *[]CompoundItem
	Sort()
}

type CompoundModel interface {
	Model
	GetItems() CompoundItems
	SetItems(*[]CompoundItem)
}

type CompoundModels interface {
	ToCompoundModels() *[]CompoundModel
}
