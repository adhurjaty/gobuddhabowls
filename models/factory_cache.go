package models

import (
	"buddhabowls/helpers"
	"errors"
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

var _categoriesCache *ItemCategories
var _invItemCache *InventoryItems
var _recipesCache *Recipes
var _orderItemsCache *OrderItems
var _vendorItemsCache *VendorItems
var _countInvItemsCache *CountInventoryItems
var _recipeItemsCache *RecipeItems
var _prepItemsCache *PrepItems
var _countPrepItemsCache *CountPrepItems

func resetCache() {
	_categoriesCache = nil
	_invItemCache = nil
	_recipesCache = nil
	_orderItemsCache = nil
	_vendorItemsCache = nil
	_countInvItemsCache = nil
	_recipeItemsCache = nil
	_prepItemsCache = nil
	_countPrepItemsCache = nil
}

func populateCategories(item *InventoryItem, tx *pop.Connection) error {
	if _categoriesCache == nil {
		_categoriesCache = &ItemCategories{}
		if err := tx.All(_categoriesCache); err != nil {
			return err
		}
	}

	for i, category := range *_categoriesCache {
		if item.CategoryID.String() == category.ID.String() {
			item.Category = (*_categoriesCache)[i]
			return nil
		}
	}

	return errors.New("no matching category ID")
}

func populateRecipe(item *PrepItem) error {
	for i, recipe := range *_recipesCache {
		if item.BatchRecipeID.String() == recipe.ID.String() {
			item.BatchRecipe = (*_recipesCache)[i]
			return nil
		}
	}

	fmt.Println(item.BatchRecipeID)
	return errors.New("no matching Recipe ID")
}

func populateInvItemCache(tx *pop.Connection) error {
	if _invItemCache != nil {
		return nil
	}

	_invItemCache = &InventoryItems{}
	if err := tx.All(_invItemCache); err != nil {
		return err
	}

	for i := range *_invItemCache {
		if err := populateCategories(&(*_invItemCache)[i], tx); err != nil {
			return err
		}
	}

	_invItemCache.Sort()

	return nil
}

func populatePrepItemsCache(itemList *PrepItems, tx *pop.Connection) error {
	if _prepItemsCache != nil {
		return nil
	}

	_prepItemsCache = itemList
	if err := tx.All(_prepItemsCache); err != nil {
		return err
	}

	ids := toIDListWithFunc(_prepItemsCache, func(item Model) uuid.UUID {
		prep := item.(*PrepItem)
		return prep.BatchRecipeID
	})
	fmt.Println(ids)
	queryIds := toIntefaceList(ids)

	if _recipesCache == nil {
		_recipesCache = &Recipes{}
	}
	if err := LoadRecipes(_recipesCache, tx.Where("id IN (?)", queryIds...)); err != nil {
		return err
	}

	fmt.Println(_recipesCache)

	for i := range *_prepItemsCache {
		if err := populateRecipe(&(*_prepItemsCache)[i]); err != nil {
			return err
		}
	}

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

func populateCountPrepItemsCache(tx *pop.Connection, ids []string) error {
	if _countPrepItemsCache != nil || len(ids) == 0 {
		return nil
	}

	return initCache(&CountPrepItems{}, tx, ids)
}

func populateRecipeItemsCache(tx *pop.Connection, ids []string) error {
	if _recipeItemsCache != nil || len(ids) == 0 {
		return nil
	}

	return initCache(&RecipeItems{}, tx, ids)
}

func populateRecipesAndItemsCache(tx *pop.Connection, ids []string) error {
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
			if err == nil {
				item.SetBaseItem(baseItem)
			} else if !helpers.IsBlankUUID(item.GetBaseItemID()) {
				ids = append(ids, item.GetBaseItemID().String())
			}
		}
		*_recipeItemsCache = append(*_recipeItemsCache, (*addRecItems)...)

		if len(ids) == 0 {
			break
		}
		if len(ids) == prevLen {
			return fmt.Errorf("could not find recipe IDs: %v", ids)
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
		return populateRecipesAndItemsCache(tx, ids)
	case *CountPrepItems:
		_countPrepItemsCache = initVal.(*CountPrepItems)
		cache = _countPrepItemsCache
		idCol = "prep_item_id"
	default:
		return errors.New("unimplemented type")
	}

	idsInt := toIntefaceList(ids)
	if err := tx.Eager().Where(fmt.Sprintf("%s IN (?)", idCol), idsInt...).
		All(cache); err != nil {
		return err
	}

	err := populateBaseItems(cache)

	return err
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
	case *PrepItem:
		return _prepItemsCache, nil
	case *CountPrepItem:
		return _countPrepItemsCache, nil
	}
	return nil, errors.New("unimplemented type")
}

func toIDList(m Models) []string {
	return toIDListWithFunc(m, func(item Model) uuid.UUID {
		return item.GetID()
	})
}

type getIdFnc func(m Model) uuid.UUID

func toIDListWithFunc(m Models, idGetter getIdFnc) []string {
	lst := m.ToModels()

	ids := make([]string, len(*lst))
	for i, item := range *lst {
		ids[i] = idGetter(item).String()
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
