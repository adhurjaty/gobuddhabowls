package presentation

import (
	"buddhabowls/models"
	"encoding/json"
	"github.com/gobuffalo/uuid"
)

// ItemAPI is an object for serving order and vendor items to UI
type ItemAPI struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	Category        CategoryAPI `json:"Category"`
	Index           int         `json:"index"`
	InventoryItemID string      `json:"inventory_item_id,omitempty"`
	Count           float64     `json:"count,omitempty"`
	Price           float64     `json:"price"`
	PurchasedUnit   string      `json:"purchased_unit,omitempty"`
	Conversion      float64     `json:"conversion,omitempty"`
}

type ItemsAPI []ItemAPI

func (item ItemAPI) String() string {
	jo, _ := json.Marshal(item)
	return string(jo)
}

func (items ItemsAPI) String() string {
	jo, _ := json.Marshal(items)
	return string(jo)
}

// NewItemAPI converts an order/vendor/inventory item to an api item
func NewItemAPI(item models.GenericItem) ItemAPI {
	itemAPI := ItemAPI{
		ID:              item.GetID().String(),
		InventoryItemID: item.GetID().String(),
		Name:            item.GetName(),
		Category:        NewCategoryAPI(item.GetCategory()),
		Index:           item.GetIndex(),
	}

	switch item.(type) {
	case models.OrderItem:
		orderItem, _ := item.(models.OrderItem)
		itemAPI.InventoryItemID = orderItem.InventoryItemID.String()
		itemAPI.Count = orderItem.Count
		itemAPI.Price = orderItem.Price

	case models.VendorItem:
		vendorItem, _ := item.(models.VendorItem)
		itemAPI.InventoryItemID = vendorItem.InventoryItemID.String()
		itemAPI.Price = vendorItem.Price
		itemAPI.PurchasedUnit = vendorItem.PurchasedUnit.String
		itemAPI.Conversion = vendorItem.Conversion
	}

	return itemAPI
}

// NewItemsAPI converts an order/vendor/inventory item slice to an api item slice
func NewItemsAPI(modelItems interface{}) ItemsAPI {
	var apis []ItemAPI
	// TODO: gotta be a better way to do this
	switch modelItems.(type) {
	case models.OrderItems:
		modelSlice := modelItems.(models.OrderItems)
		apis = make([]ItemAPI, len(modelSlice))
		for i, modelItem := range modelSlice {
			apis[i] = NewItemAPI(modelItem)
		}
	case models.VendorItems:
		modelSlice := modelItems.(models.VendorItems)
		apis = make([]ItemAPI, len(modelSlice))
		for i, modelItem := range modelSlice {
			apis[i] = NewItemAPI(modelItem)
		}
	case models.InventoryItems:
		modelSlice := modelItems.(models.InventoryItems)
		apis = make([]ItemAPI, len(modelSlice))
		for i, modelItem := range modelSlice {
			apis[i] = NewItemAPI(modelItem)
		}
	}

	return apis
}

func ConvertToModelOrderItem(item ItemAPI, orderID uuid.UUID) (*models.OrderItem, error) {
	id := uuid.UUID{}
	if len(item.ID) > 0 {
		var err error
		id, err = uuid.FromString(item.ID)
		if err != nil {
			return nil, err
		}
	}
	invID, err := uuid.FromString(item.InventoryItemID)
	if err != nil {
		return nil, err
	}
	return &models.OrderItem{
		ID:              id,
		InventoryItemID: invID,
		Count:           item.Count,
		Price:           item.Price,
		OrderID:         orderID,
	}, nil
}

func ConvertToModelOrderItems(items ItemsAPI, orderID uuid.UUID) (*models.OrderItems, error) {
	modelItems := models.OrderItems{}
	for _, item := range items {
		modelItem, err := ConvertToModelOrderItem(item, orderID)
		if err != nil {
			return nil, err
		}
		modelItems = append(modelItems, *modelItem)
	}
	return &modelItems, nil
}

func ConvertToModelVendorItem(item ItemAPI, vendorID uuid.UUID) (*models.VendorItem, error) {
	id := uuid.UUID{}
	if len(item.ID) > 0 {
		var err error
		id, err = uuid.FromString(item.ID)
		if err != nil {
			return nil, err
		}
	}

	invID, err := uuid.FromString(item.InventoryItemID)
	if err != nil {
		return nil, err
	}
	return &models.VendorItem{
		ID:              id,
		InventoryItemID: invID,
		Price:           item.Price,
		PurchasedUnit:   models.StringToNullString(item.PurchasedUnit),
		Conversion:      item.Conversion,
		VendorID:        vendorID,
	}, nil
}

func ConvertToModelVendorItems(items ItemsAPI, vendorID uuid.UUID) (*models.VendorItems, error) {
	modelItems := models.VendorItems{}
	for _, item := range items {
		modelItem, err := ConvertToModelVendorItem(item, vendorID)
		if err != nil {
			return nil, err
		}
		modelItems = append(modelItems, *modelItem)
	}
	return &modelItems, nil
}

// AddVendorInfo adds the vendorItem-specific data to the item
func AddVendorInfo(items ItemsAPI, vendorItems ItemsAPI) ItemsAPI {
	outItems := []ItemAPI{}
	for _, item := range items {
		for _, vendorItem := range vendorItems {
			if item.InventoryItemID == vendorItem.InventoryItemID {
				item.PurchasedUnit = vendorItem.PurchasedUnit
				item.Conversion = vendorItem.Conversion
				break
			}
		}
		outItems = append(outItems, item)
	}

	return outItems
}

// SelectValue returns the ID for select input tags
func (item ItemAPI) SelectValue() interface{} {
	return item.ID
}

// SelectLabel returs the name for select input tags
func (item ItemAPI) SelectLabel() string {
	if item.ID == "" {
		return "- Select an item -"
	}
	return item.Name
}
