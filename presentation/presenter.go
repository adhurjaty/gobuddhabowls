package presentation

import (
	"buddhabowls/logic"
	"buddhabowls/models"
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"time"
)

var _ = fmt.Printf

// periodSelectorContext context for the _period_selector.html template
type periodSelectorContext struct {
	PeriodSelector logic.PeriodSelector
	SelectedYear   int
	SelectedPeriod logic.Period
	SelectedWeek   logic.Week
	Years          []int
}

// Presenter readies values for ViewModel consumption
type Presenter struct {
	tx            *pop.Connection
	PeriodContext *periodSelectorContext
}

func NewPresenter(tx *pop.Connection) *Presenter {
	p := &Presenter{tx, nil}
	return p
}

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

func (p *Presenter) GetVendors() (*VendorsAPI, error) {
	vendors, err := logic.GetAllVendors(p.tx)
	if err != nil {
		return nil, err
	}

	apiVendors := NewVendorsAPI(vendors)
	return &apiVendors, nil
}

func (p *Presenter) GetVendor(id string) (*VendorAPI, error) {
	vendor, err := logic.GetVendor(id, p.tx)
	if err != nil {
		return nil, err
	}

	apiVendor := NewVendorAPI(vendor)
	return &apiVendor, nil
}

func (p *Presenter) InsertVendor(vendorAPI *VendorAPI) (*validate.Errors, error) {
	vendor, err := ConvertToModelVendor(vendorAPI)
	if err != nil {
		return nil, err
	}
	return logic.InsertVendor(vendor, p.tx)
}

func (p *Presenter) UpdateVendor(vendAPI *VendorAPI) (*validate.Errors, error) {
	vendor, err := ConvertToModelVendor(vendAPI)
	if err != nil {
		return nil, err
	}

	return logic.UpdateVendor(vendor, p.tx)
}

func (p *Presenter) GetInventoryItems() (*ItemsAPI, error) {
	items, err := logic.GetInventoryItems(p.tx)
	if err != nil {
		return nil, err
	}

	apiItems := NewItemsAPI(*items)

	return &apiItems, err
}

// GetPeriods gets the list of periods available to the user
func (p *Presenter) GetPeriods(startTime time.Time) []logic.Period {
	if p.PeriodContext == nil {
		p.setPeriodContext(startTime)
	}

	return p.PeriodContext.PeriodSelector.Periods
}

// GetSelectedPeriod gets the period that contains the startTime
func (p *Presenter) GetSelectedPeriod(startTime time.Time) logic.Period {
	if p.PeriodContext == nil {
		p.setPeriodContext(startTime)
	}

	return p.PeriodContext.PeriodSelector.GetPeriod(startTime)
}

// GetWeeks gets the list of weeks available to the user
func (p *Presenter) GetWeeks(startTime time.Time) []logic.Week {
	if p.PeriodContext == nil {
		p.setPeriodContext(startTime)
	}

	return p.PeriodContext.SelectedPeriod.Weeks
}

// GetSelectedWeek gets the week that contains the startTime
func (p *Presenter) GetSelectedWeek(startTime time.Time) logic.Week {
	if p.PeriodContext == nil {
		p.setPeriodContext(startTime)
	}

	return p.PeriodContext.PeriodSelector.GetWeek(startTime)
}

// GetYears gets the list of years available to the user
func (p *Presenter) GetYears() []int {
	if p.PeriodContext == nil {
		years, err := models.GetYears(p.tx)
		if err != nil {
			return []int{1989}
		}
		return years
	}
	return p.PeriodContext.Years
}

// setPeriodContext sets a period selector context for using the period selector UI
func (p *Presenter) setPeriodContext(date time.Time) {
	p.PeriodContext = &periodSelectorContext{}
	p.PeriodContext.PeriodSelector = logic.NewPeriodSelector(date.Year())
	p.PeriodContext.SelectedPeriod = p.PeriodContext.PeriodSelector.GetPeriod(date)
	p.PeriodContext.SelectedWeek = p.PeriodContext.SelectedPeriod.GetWeek(date)
	p.PeriodContext.SelectedYear = date.Year()
	p.PeriodContext.Years, _ = models.GetYears(p.tx)
}
