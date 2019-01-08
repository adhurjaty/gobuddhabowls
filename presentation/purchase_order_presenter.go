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
	return logic.InsertPurchaseOrder(purchaseOrder, p.tx)
}

func (p *Presenter) UpdatePurchaseOrder(poAPI *PurchaseOrderAPI) (*validate.Errors, error) {
	purchaseOrder, err := ConvertToModelPurchaseOrder(poAPI)
	if err != nil {
		return nil, err
	}
	return logic.UpdatePurchaseOrder(purchaseOrder, p.tx)
}

func (p *Presenter) DestroyPurchaseOrder(poAPI *PurchaseOrderAPI) error {
	purchaseOrder, err := ConvertToModelPurchaseOrder(poAPI)
	if err != nil {
		return err
	}
	return logic.DeletePurchaseOrder(purchaseOrder, p.tx)
}
