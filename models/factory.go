package models

import (
	"errors"
	"github.com/gobuffalo/pop"
)

// Factory describes an abstract factory for creating model objects
type Factory interface {
	CreateModel(m interface{}, tx *pop.Connection, id string) error
	CreateModelSlice(s interface{}, q *pop.Query) error
}

// ModelFactory is the concrete implementation of the Factory interface
type ModelFactory struct{}

// CreateModel loads a full model based on the type m given
func (mf *ModelFactory) CreateModel(m interface{}, tx *pop.Connection, id string) error {
	switch m.(type) {
	case *PurchaseOrder:
		if err := LoadPurchaseOrder(m.(*PurchaseOrder), tx, id); err != nil {
			return err
		}
		return nil
	case *Vendor:
		if err := LoadVendor(m.(*Vendor), tx, id); err != nil {
			return err
		}
		return nil
	}
	return errors.New("unimplemented type")
}

// CreateModelSlice loads a full model slice based on the type s given
func (mf *ModelFactory) CreateModelSlice(s interface{}, q *pop.Query) error {
	switch s.(type) {
	case *PurchaseOrders:
		if err := LoadPurchaseOrders(s.(*PurchaseOrders), q); err != nil {
			return err
		}
		return nil
	case *Vendors:
		if err := LoadVendors(s.(*Vendors), q); err != nil {
			return err
		}
		return nil
	}

	return errors.New("unimplemented type")
}

// LoadPurchaseOrder gets purchase order and sub-components matching the given ID
func LoadPurchaseOrder(po *PurchaseOrder, tx *pop.Connection, id string) error {
	if err := tx.Eager().Find(po, id); err != nil {
		return err
	}

	err := PopulateOrderItems(&po.Items, tx)
	return err
}

// LoadPurchaseOrders gets the purchase orders with the specified query
// including all sub-components
func LoadPurchaseOrders(poList *PurchaseOrders, q *pop.Query) error {
	if err := q.All(poList); err != nil {
		return err
	}

	// I don't love the fact that I need to load the nested models manually
	// TODO: look for a solution to eager loading nested objects
	for _, po := range *poList {
		if err := PopulateOrderItems(&po.Items, q.Connection); err != nil {
			return err
		}
	}

	return nil
}

// LoadOrderItem fully loads the order item
func LoadOrderItem(item *OrderItem, tx *pop.Connection, id string) error {
	if err := tx.Eager().Find(item, id); err != nil {
		return err
	}
	err := tx.Eager().Find(&item.InventoryItem, item.InventoryItemID)

	return err
}

// PopulateOrderItems populates the existing order items slice
func PopulateOrderItems(items *OrderItems, tx *pop.Connection) error {
	for i := 0; i < len(*items); i++ {
		count := (*items)[i].Count
		if err := LoadOrderItem(&(*items)[i], tx, (*items)[i].ID.String()); err != nil {
			return err
		}
		(*items)[i].Count = count
	}

	items.Sort()

	return nil
}

// LoadVendor gets a vendor given ID
func LoadVendor(vendor *Vendor, tx *pop.Connection, id string) error {
	if err := tx.Eager().Find(vendor, id); err != nil {
		return err
	}

	for i := 0; i < len(vendor.Items); i++ {
		if err := tx.Eager().Find(&vendor.Items[i], vendor.Items[i].ID); err != nil {
			return err
		}
		if err := tx.Eager().Find(&vendor.Items[i].InventoryItem, vendor.Items[i].InventoryItemID); err != nil {
			return err
		}
	}

	vendor.Items.Sort()

	return nil
}

// LoadVendors loads vendors and sub-models
func LoadVendors(vendList *Vendors, q *pop.Query) error {
	if err := q.All(vendList); err != nil {
		return err
	}

	// I don't love the fact that I need to load the nested models manually
	// TODO: look for a solution to eager loading nested objects
	for _, v := range *vendList {
		for i := 0; i < len(v.Items); i++ {
			if err := q.Connection.Eager().Find(&v.Items[i], v.Items[i].ID); err != nil {
				return err
			}
			if err := q.Connection.Eager().Find(&v.Items[i].InventoryItem, v.Items[i].InventoryItemID); err != nil {
				return err
			}
		}

		v.Items.Sort()
	}

	return nil
}
