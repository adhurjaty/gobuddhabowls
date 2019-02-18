package models

import (
	"errors"
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

var _invItemCache *InventoryItems
var _orderItemsCache *OrderItems
var _vendorItemsCache *VendorItems
var _countInvItemsCache *CountInventoryItems
var _recipeItemsCache *RecipeItems

func resetCache() {
	_invItemCache = nil
	_orderItemsCache = nil
	_vendorItemsCache = nil
	_countInvItemsCache = nil
	_recipeItemsCache = nil
}

func populateInvItemCache(tx *pop.Connection) error {
	if _invItemCache != nil {
		return nil
	}

	_invItemCache = &InventoryItems{}
	return LoadInventoryItems(_invItemCache, tx.Eager().Q())
}

func populateOrderItemsCache(tx *pop.Connection, ids []string) error {
	if _orderItemsCache != nil || len(ids) == 0 {
		return nil
	}

	return initCache(&OrderItems{}, tx, ids)
}

func populateVendorItemsCache(tx *pop.Connection, ids []string) error {
	if _vendorItemsCache != nil || len(ids) == 0 {
		return nil
	}

	return initCache(&VendorItems{}, tx, ids)
}

func populateCountInvItemsCache(tx *pop.Connection, ids []string) error {
	if _countInvItemsCache != nil || len(ids) == 0 {
		return nil
	}

	return initCache(&CountInventoryItems{}, tx, ids)
}

func populateRecipeItemsCache(tx *pop.Connection, ids []string) error {

	return nil
}

func initCache(initVal CompoundItems, tx *pop.Connection, ids []string) error {
	var cache CompoundItems
	var idCol string

	switch initVal.(type) {
	case *OrderItems:
		_orderItemsCache = initVal.(*OrderItems)
		cache = _orderItemsCache
		idCol = "order_id"
	case *VendorItems:
		_vendorItemsCache = initVal.(*VendorItems)
		cache = _vendorItemsCache
		idCol = "vendor_id"
	case *CountInventoryItems:
		_countInvItemsCache = initVal.(*CountInventoryItems)
		cache = _countInvItemsCache
		idCol = "inventory_id"
	case *RecipeItems:
		_recipeItemsCache = initVal.(*RecipeItems)
		cache = _recipeItemsCache
		idCol = "recipe_id"
	default:
		return errors.New("unimplemented type")
	}

	if err := populateInvItemCache(tx); err != nil {
		return err
	}

	idsInt := toIntefaceList(ids)
	if err := tx.Eager().Where(fmt.Sprintf("%s IN (?)", idCol), idsInt...).
		All(cache); err != nil {
		return err
	}

	fmt.Println(cache)
	fmt.Println(_countInvItemsCache)

	cacheItems := cache.ToCompoundItems()
	for i := range *cacheItems {
		item := (*cacheItems)[i]
		invItem, err := getCacheItem(item.GetBaseItem(), item.GetBaseItemID())
		if err != nil {
			return err
		}
		item.SetBaseItem(invItem)
	}

	return nil
}

func getCacheItem(itemProp GenericItem, id uuid.UUID) (GenericItem, error) {
	var cache GenericItems
	switch itemProp.(type) {
	case *InventoryItem:
		cache = _invItemCache
	case *OrderItem:
		cache = _orderItemsCache
	case *VendorItem:
		cache = _vendorItemsCache
	case *CountInventoryItem:
		cache = _countInvItemsCache
	case *RecipeItem:
		cache = _recipeItemsCache
	default:
		return nil, errors.New("unimplemented type")
	}

	cacheItems := cache.ToGenericItems()
	for _, item := range *cacheItems {
		if item.GetID().String() == id.String() {
			return item, nil
		}
	}

	return nil, errors.New("no item ID matches")
}

func toIDList(m Models) []string {
	lst := m.ToModels()

	ids := make([]string, len(*lst))
	for i, item := range *lst {
		ids[i] = item.GetID().String()
	}

	return ids
}

func toIntefaceList(lst []string) []interface{} {
	out := make([]interface{}, len(lst))
	for i := range lst {
		out[i] = lst[i]
	}

	return out
}
