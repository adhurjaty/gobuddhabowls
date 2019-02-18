package models

import (
	"errors"
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

type recipeBaseItems []GenericItem

func (r *recipeBaseItems) ToGenericItems() *[]GenericItem {
	model := make([]GenericItem, len(*r))
	for idx, item := range *r {
		model[idx] = item
	}

	return &model
}

var _invItemCache *InventoryItems
var _recipeBaseItemCache *recipeBaseItems
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
	_recipeBaseItemCache = nil
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

func populateRecipeBaseItemsCache(tx *pop.Connection) error {
	_recipeBaseItemCache = &recipeBaseItems{}
	*_recipeBaseItemCache = append(*_recipeBaseItemCache,
		*_invItemCache.ToGenericItems()...)

	batchRecipes, err := getFullBatchRecipes(tx)
	if err != nil {
		return err
	}

	*_recipeBaseItemCache = append(*_recipeBaseItemCache,
		*batchRecipes.ToGenericItems()...)

	return nil
}

func populateRecipeItemsCache(tx *pop.Connection, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	if _recipeItemsCache == nil {
		_recipeItemsCache = &RecipeItems{}
	}

	return initCache(_recipeItemsCache, tx, ids)
}

func getFullBatchRecipes(tx *pop.Connection) (*Recipes, error) {
	batchRecipes := &Recipes{}
	if err := tx.Eager().Where("is_batch = true").All(batchRecipes); err != nil {
		return nil, err
	}

	ids := toIDList(batchRecipes)
	_recipeItemsCache = &RecipeItems{}
	if err := tx.Eager().Where("recipe_id IN (?)", toIntefaceList(ids)...).
		All(_recipeItemsCache); err != nil {
		return nil, err
	}

	err := setModelItemsFromCache(batchRecipes)

	return batchRecipes, err
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
		_recipeItemsCache = initVal.(*RecipeItems)
		cache = _recipeItemsCache
		idCol = "recipe_id"
		if err := populateRecipeBaseItemsCache(tx); err != nil {
			return err
		}
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

func populateBaseItems(cache CompoundItems) error {
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
