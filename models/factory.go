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
	// maybe split this according to compound vs. item type
	resetCache()
	switch m.(type) {
	case *PurchaseOrder:
		return LoadPurchaseOrder(m.(*PurchaseOrder), tx, id)
	case *Vendor:
		return LoadVendor(m.(*Vendor), tx, id)
	case *Recipe:
		err := populateInvItemCache(tx)
		if err != nil {
			return err
		}
		return LoadRecipe(m.(*Recipe), tx, id)
	case *Inventory:
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
	resetCache()
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

// PopulateOrderItems populates the existing order items slice
func PopulateOrderItems(pos *PurchaseOrders, tx *pop.Connection) error {
	ids := toIDList(pos)

	if err := populateOrderItemsCache(tx, ids); err != nil {
		return err
	}
	if _orderItemsCache == nil {
		return nil
	}

	return setModelItemsFromCache(pos)
}

// LoadVendor gets a vendor given ID
func LoadVendor(vendor *Vendor, tx *pop.Connection, id string) error {
	if err := tx.Eager().Find(vendor, id); err != nil {
		return err
	}

	return PopulateVendorItems(&Vendors{*vendor}, tx)
}

// LoadVendors loads vendors and sub-models
func LoadVendors(vendList *Vendors, q *pop.Query) error {
	if err := q.All(vendList); err != nil {
		return err
	}

	return PopulateVendorItems(vendList, q.Connection)
}

func PopulateVendorItems(vendors *Vendors, tx *pop.Connection) error {
	ids := toIDList(vendors)

	if err := populateVendorItemsCache(tx, ids); err != nil {
		return err
	}
	if _vendorItemsCache == nil {
		return nil
	}

	return setModelItemsFromCache(vendors)
}

func LoadInventory(inventory *Inventory, tx *pop.Connection, id string) error {
	if err := tx.Eager().Find(inventory, id); err != nil {
		return err
	}

	return PopulateCountInvItems(&Inventories{*inventory}, tx)
}

func LoadInventories(invList *Inventories, q *pop.Query) error {
	if err := q.All(invList); err != nil {
		return err
	}

	return PopulateCountInvItems(invList, q.Connection)
}

func PopulateCountInvItems(inventories *Inventories, tx *pop.Connection) error {
	ids := toIDList(inventories)

	if err := populateCountInvItemsCache(tx, ids); err != nil {
		return err
	}
	if _countInvItemsCache == nil {
		return nil
	}

	return setModelItemsFromCache(inventories)
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

	return PopulateRecipeItems(&Recipes{*recipe}, tx)
}

func LoadRecipes(recipes *Recipes, q *pop.Query) error {
	if err := q.All(recipes); err != nil {
		return err
	}

	return PopulateRecipeItems(recipes, q.Connection)
}

func PopulateRecipeItems(recipes *Recipes, tx *pop.Connection) error {
	ids := toIDList(recipes)

	if err := populateRecipeItemsCache(tx, ids); err != nil {
		return err
	}
	if _recipeItemsCache == nil {
		return nil
	}

	return setModelItemsFromCache(recipes)
}

func LoadRecipeItem(item *RecipeItem, tx *pop.Connection, id string) error {
	err := tx.Eager().Find(item, id)
	if err != nil {
		return err
	}

	if item.InventoryItemID.Valid {
		err = tx.Eager().Find(&item.InventoryItem, item.InventoryItemID.UUID)
	} else if item.BatchRecipeID.Valid {
		err = tx.Eager().Find(&item.BatchRecipe, item.BatchRecipeID.UUID)
	}

	return err
}

func setModelItemsFromCache(models CompoundModels) error {
	modelList := models.ToCompoundModels()
	for _, m := range *modelList {
		items := m.GetItems().ToCompoundItems()
		for i := 0; i < len(*items); i++ {
			genItem := (*items)[i].(GenericItem)
			cacheItem, err := getCacheItem(genItem, (*items)[i].GetID())
			if err != nil {
				return err
			}
			(*items)[i] = cacheItem.(CompoundItem)
		}
		m.SetItems(items)
		m.GetItems().Sort()
	}

	return nil
}
