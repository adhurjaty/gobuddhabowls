package models

import (
	"errors"
	"github.com/gobuffalo/pop"
)

var _invItemCache *InventoryItems
var _orderItemsCache *OrderItems

// Factory describes an abstract factory for creating model objects
type Factory interface {
	CreateModel(m interface{}, tx *pop.Connection, id string) error
	CreateModelSlice(s interface{}, q *pop.Query) error
}

// ModelFactory is the concrete implementation of the Factory interface
type ModelFactory struct{}

// CreateModel loads a full model based on the type m given
func (mf *ModelFactory) CreateModel(m interface{}, tx *pop.Connection, id string) error {
	// maybe split this according to compound vs. item type

	switch m.(type) {
	case *PurchaseOrder:
		return LoadPurchaseOrder(m.(*PurchaseOrder), tx, id)
	case *Vendor:
		err := populateInvItemCache(tx)
		if err != nil {
			return err
		}
		return LoadVendor(m.(*Vendor), tx, id)
	case *Recipe:
		err := populateInvItemCache(tx)
		if err != nil {
			return err
		}
		return LoadRecipe(m.(*Recipe), tx, id)
	case *Inventory:
		err := populateInvItemCache(tx)
		if err != nil {
			return err
		}
		return LoadInventory(m.(*Inventory), tx, id)
	case *RecipeItem:
		return LoadRecipeItem(m.(*RecipeItem), tx, id)
	case *InventoryItem:
		return LoadInventoryItem(m.(*InventoryItem), tx, id)
	case *VendorItem:
		return LoadVendorItem(m.(*VendorItem), tx, id)
	}

	return errors.New("unimplemented type")
}

// CreateModelSlice loads a full model slice based on the type s given
func (mf *ModelFactory) CreateModelSlice(s interface{}, q *pop.Query) error {
	_, ok := s.(*InventoryItems)
	if !ok {
		err := populateInvItemCache(q.Connection)
		if err != nil {
			return err
		}
	}

	switch s.(type) {
	case *PurchaseOrders:
		return LoadPurchaseOrders(s.(*PurchaseOrders), q)
	case *Vendors:
		return LoadVendors(s.(*Vendors), q)
	case *InventoryItems:
		return LoadInventoryItems(s.(*InventoryItems), q)
	case *Inventories:
		return LoadInventories(s.(*Inventories), q)
	case *Recipes:
		return LoadRecipes(s.(*Recipes), q)
	}

	return errors.New("unimplemented type")
}

// LoadPurchaseOrder gets purchase order and sub-components matching the given ID
func LoadPurchaseOrder(po *PurchaseOrder, tx *pop.Connection, id string) error {
	if err := tx.Eager().Find(po, id); err != nil {
		return err
	}

	return PopulateOrderItems(&PurchaseOrders{*po}, tx)
}

// LoadPurchaseOrders gets the purchase orders with the specified query
// including all sub-components
func LoadPurchaseOrders(poList *PurchaseOrders, q *pop.Query) error {
	if err := q.All(poList); err != nil {
		return err
	}

	if err := PopulateOrderItems(poList, q.Connection); err != nil {
		return err
	}

	return nil
}

// LoadOrderItem fully loads the order item
func LoadOrderItem(item *OrderItem, tx *pop.Connection, id string) (*OrderItem, error) {
	return getOrderItem(item)
}

// PopulateOrderItems populates the existing order items slice
func PopulateOrderItems(pos *PurchaseOrders, tx *pop.Connection) error {
	ids := make([]string, len(*pos))
	for i := range *pos {
		ids[i] = (*pos)[i].ID.String()
	}

	err := populateOrderItemsCache(tx, ids)
	if err != nil {
		return err
	}
	if _orderItemsCache == nil {
		return nil
	}

	for _, po := range *pos {
		for i := 0; i < len(po.Items); i++ {
			count := po.Items[i].Count
			orderItem, err := getOrderItem(&po.Items[i])
			if err != nil {
				return err
			}
			po.Items[i] = *orderItem
			po.Items[i].Count = count
		}
		po.Items.Sort()
	}

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
		if err := getInventoryItem(&vendor.Items[i].InventoryItem, vendor.Items[i].InventoryItemID); err != nil {
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
			if err := getInventoryItem(&v.Items[i].InventoryItem, v.Items[i].InventoryItemID); err != nil {
				return err
			}
		}

		v.Items.Sort()
	}

	return nil
}

func LoadInventory(inventory *Inventory, tx *pop.Connection, id string) error {
	if err := tx.Eager().Find(inventory, id); err != nil {
		return err
	}

	for i := 0; i < len(inventory.Items); i++ {
		if err := tx.Eager().Find(&inventory.Items[i], inventory.Items[i].ID); err != nil {
			return err
		}
		if err := getInventoryItem(&inventory.Items[i].InventoryItem, inventory.Items[i].InventoryItemID); err != nil {
			return err
		}
	}

	inventory.Items.Sort()

	return nil
}

func LoadInventories(invList *Inventories, q *pop.Query) error {
	if err := q.All(invList); err != nil {
		return err
	}

	// I don't love the fact that I need to load the nested models manually
	// TODO: look for a solution to eager loading nested objects
	for _, inv := range *invList {
		for i := 0; i < len(inv.Items); i++ {
			if err := q.Connection.Eager().Find(&inv.Items[i], inv.Items[i].ID); err != nil {
				return err
			}
			if err := getInventoryItem(&inv.Items[i].InventoryItem, inv.Items[i].InventoryItemID); err != nil {
				return err
			}
		}

		inv.Items.Sort()
	}

	return nil
}

func LoadInventoryItem(item *InventoryItem, tx *pop.Connection, id string) error {
	return tx.Eager().Find(item, id)
}

func LoadInventoryItems(itemList *InventoryItems, q *pop.Query) error {
	if err := q.All(itemList); err != nil {
		return err
	}

	itemList.Sort()

	return nil
}

func LoadRecipe(recipe *Recipe, tx *pop.Connection, id string) error {
	if err := tx.Eager().Find(recipe, id); err != nil {
		return err
	}

	return PopulateRecipeItems(&recipe.Items, tx)
}

func LoadRecipes(recipes *Recipes, q *pop.Query) error {
	if err := q.All(recipes); err != nil {
		return err
	}

	for _, recipe := range *recipes {
		if err := PopulateRecipeItems(&recipe.Items, q.Connection); err != nil {
			return err
		}
	}

	return nil
}

func PopulateRecipeItems(items *RecipeItems, tx *pop.Connection) error {
	for i, _ := range *items {
		item := &(*items)[i]
		count := item.Count
		if err := LoadRecipeItem(item, tx, item.ID.String()); err != nil {
			return err
		}

		(*items)[i].Count = count
	}

	items.Sort()

	return nil
}

func LoadRecipeItem(item *RecipeItem, tx *pop.Connection, id string) error {
	err := tx.Eager().Find(item, id)
	if err != nil {
		return err
	}
	if item.InventoryItemID.Valid {
		// err = tx.Eager().Find(&item.InventoryItem, item.InventoryItemID.UUID)
		err = getInventoryItem(&item.InventoryItem, item.InventoryItemID.UUID)
	} else if item.BatchRecipeID.Valid {
		err = tx.Eager().Find(&item.BatchRecipe, item.BatchRecipeID.UUID)
	}

	return err
}
