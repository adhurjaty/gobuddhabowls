package models

import (
	"errors"
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

var _invItemCache *InventoryItems
var _recipesCache *Recipes
var _orderItemsCache *OrderItems
var _vendorItemsCache *VendorItems
var _countInvItemsCache *CountInventoryItems
var _recipeItemsCache *RecipeItems

func resetCache() {
	_invItemCache = nil
	_recipesCache = nil
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
	if err := tx.All(_invItemCache); err != nil {
		return err
	}

	_invItemCache.Sort()

	return nil
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
	if _recipeItemsCache != nil || len(ids) == 0 {
		return nil
	}

	return initCache(&RecipeItems{}, tx, ids)
}

func populateRecipesCache(tx *pop.Connection, ids []string) error {
	_recipesCache = &Recipes{}
	_recipeItemsCache = &RecipeItems{}

	prevLen := -1

	for true {
		if prevLen == len(ids) {
		}
		prevLen = len(ids)

		addRecItems := &RecipeItems{}
		if err := tx.Eager().Where("recipe_id IN (?)", toIntefaceList(ids)...).
			All(addRecItems); err != nil {
			return err
		}

		ids = []string{}

		cacheItems := addRecItems.ToCompoundItems()
		for i := range *cacheItems {
			item := (*cacheItems)[i]
			baseItem, err := getCacheItem(item.GetBaseItem(), item.GetBaseItemID())
			if err != nil {
				ids = append(ids, item.GetBaseItemID().String())
			}
			item.SetBaseItem(baseItem)
		}
		*_recipeItemsCache = append(*_recipeItemsCache, (*addRecItems)...)

		if len(ids) == 0 {
			break
		}
		if len(ids) == prevLen {
			return errors.New(fmt.Sprintf("could not find recipe IDs: %v", ids))
		}

		additionalRecipes := &Recipes{}
		if err := tx.Eager().Where("id IN (?)", toIntefaceList(ids)...).
			All(additionalRecipes); err != nil {
			return err
		}

		*_recipesCache = append(*_recipesCache, (*additionalRecipes)...)
	}

	return nil
}

func initCache(initVal CompoundItems, tx *pop.Connection, ids []string) error {
	var cache CompoundItems
	var idCol string

	if err := populateInvItemCache(tx); err != nil {
		return err
	}

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
		return populateRecipesCache(tx, ids)
	default:
		return errors.New("unimplemented type")
	}

	idsInt := toIntefaceList(ids)
	if err := tx.Eager().Where(fmt.Sprintf("%s IN (?)", idCol), idsInt...).
		All(cache); err != nil {
		return err
	}

	return populateBaseItems(cache)
}

func populateBaseItems(cache CompoundItems) error {
	cacheItems := cache.ToCompoundItems()
	for i := range *cacheItems {
		item := (*cacheItems)[i]
		baseItem, err := getCacheItem(item.GetBaseItem(), item.GetBaseItemID())
		if err != nil {
			return err
		}
		item.SetBaseItem(baseItem)
	}

	return nil
}

func getCacheItem(itemProp GenericItem, id uuid.UUID) (GenericItem, error) {
	cache, err := getCacheFromType(itemProp)
	if err != nil {
		return nil, err
	}
	cacheItems := cache.ToGenericItems()
	for _, item := range *cacheItems {
		if item.GetID().String() == id.String() {
			return item, nil
		}
	}

	return nil, errors.New("no item ID matches")
}

func getCacheFromType(itemProp GenericItem) (GenericItems, error) {
	switch itemProp.(type) {
	case *InventoryItem:
		return _invItemCache, nil
	case *Recipe:
		return _recipesCache, nil
	case *OrderItem:
		return _orderItemsCache, nil
	case *VendorItem:
		return _vendorItemsCache, nil
	case *CountInventoryItem:
		return _countInvItemsCache, nil
	case *RecipeItem:
		return _recipeItemsCache, nil
	}
	return nil, errors.New("unimplemented type")
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
