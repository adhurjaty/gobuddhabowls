package presentation

import (
	"buddhabowls/logic"
	"fmt"
	"github.com/gobuffalo/validate"
	"time"
)

var _ = fmt.Println

// GetPurchaseOrders gets the purchase orders from the given date interval
func (p *Presenter) GetPurchaseOrders(startTime time.Time, endTime time.Time) (*PurchaseOrdersAPI, error) {
	purchaseOrders, err := logic.GetPurchaseOrders(startTime, endTime, p.tx)
	if err != nil {
		return nil, err
	}

	apiPurchaseOrders := NewPurchaseOrdersAPI(purchaseOrders)

	return &apiPurchaseOrders, nil
}

// GetPurchaseOrder retrieves a purchase order by ID to display to the user
func (p *Presenter) GetPurchaseOrder(id string) (*PurchaseOrderAPI, error) {
	purchaseOrder, err := logic.GetPurchaseOrder(id, p.tx)
	if err != nil {
		return nil, err
	}

	apiPO := NewPurchaseOrderAPI(purchaseOrder)
	return &apiPO, nil
}

func (p *Presenter) InsertPurchaseOrder(poAPI *PurchaseOrderAPI) (*validate.Errors, error) {
	purchaseOrder, err := ConvertToModelPurchaseOrder(poAPI)
	if err != nil {
		return nil, err
	}

	if poAPI.ReceivedDate.Valid {
		verrs, err := p.updateVendorItemsFromPO(poAPI)
		if verrs.HasAny() || err != nil {
			return verrs, err
		}
	}

	return logic.InsertPurchaseOrder(purchaseOrder, p.tx)
}

func (p *Presenter) UpdatePurchaseOrder(poAPI *PurchaseOrderAPI) (*validate.Errors, error) {
	purchaseOrder, err := ConvertToModelPurchaseOrder(poAPI)
	if err != nil {
		return validate.NewErrors(), err
	}

	oldPO, err := p.GetPurchaseOrder(poAPI.ID)
	if err != nil {
		return validate.NewErrors(), err
	}

	if poAPI.ReceivedDate.Valid {
		verrs, err := p.updateVendorItemsFromPO(poAPI)
		if verrs.HasAny() || err != nil {
			return verrs, err
		}
	}

	verrs, err := logic.UpdatePurchaseOrder(purchaseOrder, p.tx)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	if oldPO.ReceivedDate.Valid && !poAPI.ReceivedDate.Valid {
		err = p.restoreVendorPrevPrices(oldPO)
		if err != nil {
			return verrs, err
		}
	}

	return verrs, nil
}

func (p *Presenter) updateVendorItemsFromPO(po *PurchaseOrderAPI) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	for _, item := range po.Items {
		verrs, err := p.updateVendorItemOnReceive(&item, po)
		if verrs.HasAny() || err != nil {
			return verrs, err
		}
	}

	return verrs, nil
}

func (p *Presenter) updateVendorItemOnReceive(item *ItemAPI, po *PurchaseOrderAPI) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	latestOrder, err := logic.GetLatestOrder(item.InventoryItemID,
		po.Vendor.ID, p.tx)
	if err != nil {
		// if there is no matching order, allow updating vendor item
		return verrs, nil
	}

	if po.ReceivedDate.Time.Unix() >= latestOrder.ReceivedDate.Time.Unix() {
		vendorItem, err := logic.GetVendorItemByInvItem(item.InventoryItemID,
			po.Vendor.ID, p.tx)
		if err != nil {
			return verrs, err
		}

		vendorItem.Price = item.Price
		verrs, err = logic.UpdateVendorItem(vendorItem, p.tx)
		if verrs.HasAny() || err != nil {
			return verrs, err
		}
	}

	return verrs, nil
}

func (p *Presenter) DestroyPurchaseOrder(poAPI *PurchaseOrderAPI) error {
	purchaseOrder, err := ConvertToModelPurchaseOrder(poAPI)
	if err != nil {
		return err
	}

	err = logic.DeletePurchaseOrder(purchaseOrder, p.tx)
	if err != nil || !poAPI.ReceivedDate.Valid {
		return err
	}

	return p.restoreVendorPrevPrices(poAPI)
}

func (p *Presenter) restoreVendorPrevPrices(po *PurchaseOrderAPI) error {
	for _, item := range po.Items {
		err := p.restoreVendorItemPrice(&item, po)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Presenter) restoreVendorItemPrice(item *ItemAPI,
	po *PurchaseOrderAPI) error {

	latestOrder, err := logic.GetLatestOrder(item.InventoryItemID,
		po.Vendor.ID, p.tx)
	if err != nil {
		// don't alter the price if there's no previous order for
		// the item
		return nil
	}

	if latestOrder.ReceivedDate.Time.Unix() < po.ReceivedDate.Time.Unix() {
		vendorItem, err := logic.GetVendorItemByInvItem(item.InventoryItemID,
			po.Vendor.ID, p.tx)
		if err != nil {
			return err
		}

		latestOrderItem, err := logic.GetItemFromOrder(latestOrder.ID.String(),
			item.InventoryItemID, p.tx)
		vendorItem.Price = latestOrderItem.Price
		_, err = logic.UpdateVendorItem(vendorItem, p.tx)
		if err != nil {
			return err
		}
	}

	return nil
}
