package models

import (
	"errors"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

var _invItemCache *InventoryItems
var _orderItemsCache *OrderItems
var _vendorItemsCache *VendorItems

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
