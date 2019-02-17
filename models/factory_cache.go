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

func resetCache() {
	_invItemCache = nil
	_orderItemsCache = nil
	_vendorItemsCache = nil
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

	if err := populateInvItemCache(tx); err != nil {
		return err
	}

	_orderItemsCache = &OrderItems{}
	idsInt := toIntefaceList(ids)
	err := tx.Eager().Where("order_id IN (?)", idsInt...).
		All(_orderItemsCache)
	if err != nil {
		return err
	}

	for i := range *_orderItemsCache {
		item := &(*_orderItemsCache)[i]
		err = getInventoryItem(&item.InventoryItem, item.InventoryItemID)
		if err != nil {
			return err
		}
	}

	return nil
}

func populateVendorItemsCache(tx *pop.Connection, ids []string) error {
	if _vendorItemsCache != nil || len(ids) == 0 {
		return nil
	}

	if err := populateInvItemCache(tx); err != nil {
		return err
	}

	_vendorItemsCache = &VendorItems{}
	idsInt := toIntefaceList(ids)
	if err := tx.Eager().Where("vendor_id IN (?)", idsInt...).
		All(_vendorItemsCache); err != nil {
		return err
	}

	for i := range *_vendorItemsCache {
		item := &(*_vendorItemsCache)[i]
		if err := getInventoryItem(&item.InventoryItem,
			item.InventoryItemID); err != nil {
			return err
		}
	}

	return nil
}

func initCache(initVal *GenericItems, tx *pop.Connection, ids []string) error {
	var cache *GenericItems
	var idCol string

	switch (*initVal).(type) {
	case OrderItems:
		*cache = *_orderItemsCache
		idCol = "order_id"
	case VendorItems:
		*cache = *_vendorItemsCache
		idCol = "vendor_id"
	default:
		return errors.New("unimplemented type")
	}

	*cache = *initVal

	if err := populateInvItemCache(tx); err != nil {
		return err
	}

	idsInt := toIntefaceList(ids)
	if err := tx.Eager().Where(fmt.Sprintf("%s IN (?)", idCol), idsInt...).
		All(cache); err != nil {
		return err
	}

	cacheItems := (*cache).ToGenericItems()
	for i := range *cacheItems {
		item := &(*cacheItems)[i]
		invItem := (*item).GetBaseItem()
		if err := getBaseItem(&invItem,
			(*item).GetInventoryItemID()); err != nil {
			return err
		}
	}

	return nil
}

func getBaseItem(item *GenericItem, id uuid.UUID) error {
	switch (*item).(type) {
	case InventoryItem:
		invItem := (*item).(InventoryItem)
		return getInventoryItem(&invItem, id)
	case Recipe:
		return errors.New("recipes not implemented")
	}

	return errors.New("unimplemented type")
}

func getInventoryItem(invItemProp *InventoryItem, id uuid.UUID) error {
	for _, item := range *_invItemCache {
		if item.ID.String() == id.String() {
			*invItemProp = item
			return nil
		}
	}

	return errors.New("no inventory item ID matches")
}

func getOrderItem(orderItem *OrderItem) error {
	for _, item := range *_orderItemsCache {
		if item.ID.String() == orderItem.ID.String() {
			*orderItem = item
			return nil
		}
	}

	return errors.New("No matching order item")
}

func getVendorItem(vItem *VendorItem) error {
	for _, item := range *_vendorItemsCache {
		if item.ID.String() == vItem.ID.String() {
			*vItem = item
			return nil
		}
	}

	return errors.New("No matching vendor item")
}

func toIntefaceList(lst []string) []interface{} {
	out := make([]interface{}, len(lst))
	for i := range lst {
		out[i] = lst[i]
	}

	return out
}
