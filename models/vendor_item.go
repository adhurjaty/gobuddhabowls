package models

import (
	"encoding/json"
	"sort"
	"time"

	"database/sql"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	// "github.com/gobuffalo/validate/validators"
)

type VendorItem struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	InventoryItemID uuid.UUID      `json:"inventory_item_id" db:"inventory_item_id"`
	InventoryItem   InventoryItem  `belongs_to:"inventory_item" db:"-"`
	VendorID        uuid.UUID      `json:"vendor_id" db:"vendor_id"`
	Vendor          Vendor         `belongs_to:"vendors" db:"-"`
	PurchasedUnit   sql.NullString `json:"purchased_unit" db:"purchased_unit"`
	Conversion      float64        `json:"conversion" db:"conversion"`
	Price           float64        `json:"price" db:"price"`
}

// String is not required by pop and may be deleted
func (v VendorItem) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// VendorItems is not required by pop and may be deleted
type VendorItems []VendorItem

// String is not required by pop and may be deleted
func (v VendorItems) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
// func (v *VendorItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
// 	return validate.Validate(
// 		&validators.StringIsPresent{Field: v.PurchasedUnit, Name: "PurchasedUnit"},
// 	), nil
// }

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (v *VendorItem) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (v *VendorItem) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (v *VendorItem) GetID() uuid.UUID {
	return v.ID
}

func (v *VendorItem) GetBaseItemID() uuid.UUID {
	return v.InventoryItemID
}

func (v *VendorItem) GetBaseItem() GenericItem {
	return &v.InventoryItem
}

func (v *VendorItem) SetBaseItem(item GenericItem) {
	v.InventoryItem = *item.(*InventoryItem)
}

func (v *VendorItem) GetName() string {
	return v.InventoryItem.Name
}

// GetCategory returns the inventory item category of the vendor item
func (v *VendorItem) GetCategory() ItemCategory {
	return v.InventoryItem.Category
}

func (v *VendorItem) GetCountUnit() string {
	return v.InventoryItem.CountUnit
}

func (v *VendorItem) GetIndex() int {
	return v.InventoryItem.Index
}

// ToOrderItem converts a vendor item to an order item (0 count and no order ID)
func (v *VendorItem) ToOrderItem() *OrderItem {
	return &OrderItem{
		InventoryItem:   v.InventoryItem,
		InventoryItemID: v.InventoryItemID,
		Count:           0,
		Price:           v.Price,
	}
}

// ToGenericItems converts the VendorItems to a CountItem slice
func (v VendorItems) ToGenericItems() *[]GenericItem {
	items := make([]GenericItem, len(v))
	for i := 0; i < len(v); i++ {
		items[i] = &v[i]
	}

	return &items
}

func (v *VendorItems) ToCompoundItems() *[]CompoundItem {
	items := make([]CompoundItem, len(*v))
	for i := 0; i < len(*v); i++ {
		items[i] = &(*v)[i]
	}

	return &items
}

// Sort sorts the items based on category then inventory item indices
func (v *VendorItems) Sort() {
	sort.Slice(*v, func(i, j int) bool {
		return (*v)[i].InventoryItem.GetSortValue() < (*v)[j].InventoryItem.GetSortValue()
	})
}

// ToOrderItems converts list vendor items to order items
func (v *VendorItems) ToOrderItems() *OrderItems {
	items := OrderItems{}
	for _, vi := range *v {
		oItem := vi.ToOrderItem()
		items = append(items, *oItem)
	}

	return &items
}

func LoadVendorItem(vendorItem *VendorItem, tx *pop.Connection, id string) error {
	if err := tx.Find(vendorItem, id); err != nil {
		return err
	}

	return tx.Eager().Find(&vendorItem.InventoryItem,
		vendorItem.InventoryItemID)
}
