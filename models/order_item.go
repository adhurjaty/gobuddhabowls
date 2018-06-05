package models

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type OrderItem struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
	InventoryItemID uuid.UUID     `json:"inventory_item_id" db:"inventory_item_id"`
	InventoryItem   InventoryItem `belongs_to:"inventory_items" db:"-"`
	OrderID         uuid.UUID     `json:"order_id" db:"order_id"`
	Order           PurchaseOrder `belongs_to:"purchase_orders" db:"-"`
	Count           float64       `json:"count,string,omitempty" db:"count"`
	Price           float64       `json:"price,string,omitempty" db:"price"`
}

// String is not required by pop and may be deleted
func (o OrderItem) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// OrderItems is not required by pop and may be deleted
type OrderItems []OrderItem

// String is not required by pop and may be deleted
func (o OrderItems) String() string {
	jo, _ := json.Marshal(o)
	return string(jo)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (o *OrderItem) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (o *OrderItem) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (o *OrderItem) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// GetCategory returns the inventory item category of the vendor item
func (o *OrderItem) GetCategory() InventoryItemCategory {
	return o.InventoryItem.Category
}

// ToCountItems converts the VendorItems to a CountItem slice
func (o *OrderItems) ToCountItems() []CountItem {
	items := make([]CountItem, len(*o))
	for i := 0; i < len(*o); i++ {
		items[i] = &(*o)[i]
	}

	return items
}

// Extension returns the total cost (price * count) of item
func (o *OrderItem) Extension() float64 {
	return o.Price * o.Count
}

// Sort sorts the items based on category then inventory item indices
func (o *OrderItems) Sort() {
	sort.Slice(*o, func(i, j int) bool {
		return (*o)[i].InventoryItem.GetSortValue() < (*o)[j].InventoryItem.GetSortValue()
	})
}

// LoadOrderItems populates the nested objects in a purchase order's items
func LoadOrderItems(tx *pop.Connection, po *PurchaseOrder) error {
	for i := 0; i < len(po.Items); i++ {
		count := po.Items[i].Count
		if err := tx.Eager().Find(&po.Items[i], po.Items[i].ID); err != nil {
			return err
		}
		if err := tx.Eager().Find(&po.Items[i].InventoryItem, po.Items[i].InventoryItemID); err != nil {
			return err
		}
		po.Items[i].Count = count
	}

	po.Items.Sort()

	return nil
}
