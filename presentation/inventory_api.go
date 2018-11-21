package presentation

import (
	"buddhabowls/models"
	"encoding/json"
	"github.com/gobuffalo/uuid"
	"time"
)

type InventoryAPI struct {
	ID    uuid.UUID `json:"id"`
	Date  time.Time `json:"time"`
	Items ItemsAPI  `json:"Items"`
}

type InventoriesAPI []InventoryAPI

func (inv InventoryAPI) String() string {
	jo, _ := json.Marshal(inv)
	return string(jo)
}

func (inv InventoriesAPI) String() string {
	jo, _ := json.Marshal(inv)
	return string(jo)
}

func NewInventoryAPI(inventory *models.Inventory) InventoryAPI {
	return InventoryAPI{
		ID:    inventory.ID,
		Date:  inventory.Date,
		Items: NewItemsAPI(inventory.Items),
	}
}

func NewInventoriesAPI(inventories *models.Inventories) InventoriesAPI {
	apis := make([]InventoryAPI, len(*inventories))
	for i, inventory := range *inventories {
		apis[i] = NewInventoryAPI(&inventory)
	}

	return apis
}
