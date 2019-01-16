package presentation

import (
	"buddhabowls/logic"
	"github.com/gobuffalo/validate"
	"time"
)

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

	verrs, err := p.updateVendorItemsFromPO(poAPI)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	return logic.InsertPurchaseOrder(purchaseOrder, p.tx)
}

func (p *Presenter) UpdatePurchaseOrder(poAPI *PurchaseOrderAPI) (*validate.Errors, error) {
	purchaseOrder, err := ConvertToModelPurchaseOrder(poAPI)
	if err != nil {
		return nil, err
	}

	verrs, err := p.updateVendorItemsFromPO(poAPI)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	return logic.UpdatePurchaseOrder(purchaseOrder, p.tx)
}

func (p *Presenter) updateVendorItemsFromPO(po *PurchaseOrderAPI) (*validate.Errors, error) {
	verrs := validate.NewErrors()
	for _, item := range po.Items {
		latestOrder, err := logic.GetLatestOrder(item.InventoryItemID,
			po.Vendor.ID, p.tx)
		if err != nil {
			return verrs, err
		}

		if po.OrderDate.Time.Unix() > latestOrder.OrderDate.Time.Unix() {
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
	}

	return verrs, nil
}

func (p *Presenter) DestroyPurchaseOrder(poAPI *PurchaseOrderAPI) error {
	purchaseOrder, err := ConvertToModelPurchaseOrder(poAPI)
	if err != nil {
		return err
	}

	err = logic.DeletePurchaseOrder(purchaseOrder, p.tx)
	if err != nil {
		return err
	}

	return p.restoreVendorPrevPrices(poAPI)
}

func (p *Presenter) restoreVendorPrevPrices(po *PurchaseOrderAPI) error {
	for _, item := range po.Items {
		latestOrderItem, err := logic.GetLatestOrderItem(item.InventoryItemID,
			po.Vendor.ID, p.tx)
		if err != nil {
			return err
		}

		vendorItem, err := logic.GetVendorItemByInvItem(item.InventoryItemID,
			po.Vendor.ID, p.tx)
		if err != nil {
			return err
		}

		vendorItem.Price = latestOrderItem.Price
		_, err = logic.UpdateVendorItem(vendorItem, p.tx)
		if err != nil {
			return err
		}
	}

	return nil
}
