package presentation

import (
	"buddhabowls/logic"
	"buddhabowls/models"
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
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

func (p *Presenter) GetInventories(startTime time.Time, endTime time.Time) (*InventoriesAPI, error) {
	inventories, err := logic.GetInventories(startTime, endTime, p.tx)
	if err != nil {
		return nil, err
	}

	vendors, err := p.GetVendors()
	if err != nil {
		return nil, err
	}

	apiInv := NewInventoriesAPI(inventories, vendors)

	return &apiInv, nil
}

func (p *Presenter) GetLatestInventory(date time.Time) (*InventoryAPI, error) {
	inventory, err := logic.GetLatestInventory(date, p.tx)
	if err != nil {
		return nil, err
	}

	vendors, err := p.GetVendors()
	if err != nil {
		return nil, err
	}

	apiInv := NewInventoryAPI(inventory, vendors)
	return &apiInv, nil
}

func (p *Presenter) GetInventory(id string) (*InventoryAPI, error) {
	inventory, err := logic.GetInventory(id, p.tx)
	if err != nil {
		return nil, err
	}

	apiInv := NewInventoryAPI(inventory, nil)
	return &apiInv, nil
}

func (p *Presenter) InsertInventory(invAPI *InventoryAPI) (*validate.Errors, error) {
	inventory, err := ConvertToModelInventory(invAPI)
	if err != nil {
		return nil, err
	}

	verrs, err := p.updateInvVendorItems(invAPI.Items)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	return logic.InsertInventory(inventory, p.tx)
}

func (p *Presenter) UpdateInventory(invAPI *InventoryAPI) (*validate.Errors, error) {
	inventory, err := ConvertToModelInventory(invAPI)
	if err != nil {
		return nil, err
	}

	verrs, err := p.updateInvVendorItems(invAPI.Items)
	if err != nil || verrs.HasAny() {
		return verrs, err
	}

	return logic.UpdateInventory(inventory, p.tx)
}

func (p *Presenter) updateInvVendorItems(invItems ItemsAPI) (*validate.Errors, error) {
	vendorItems := models.VendorItems{}
	for _, item := range invItems {
		subItem, ok := item.VendorItemMap[item.SelectedVendor]
		if ok {
			vendorID, err := uuid.FromString(subItem.SelectedVendor)
			if err != nil {
				continue
			}
			vItem, err := ConvertToModelVendorItem(subItem, vendorID)
			if err != nil {
				continue
			}

			vendorItems = append(vendorItems, *vItem)
		}
	}

	return logic.UpdateVendorItems(&vendorItems, p.tx)
}

func (p *Presenter) DestroyInventory(invAPI *InventoryAPI) error {
	inventory, err := ConvertToModelInventory(invAPI)
	if err != nil {
		return err
	}

	return logic.DestroyInventory(inventory, p.tx)
}

func (p *Presenter) GetInventoryItems() (*ItemsAPI, error) {
	items, err := logic.GetInventoryItems(p.tx)
	if err != nil {
		return nil, err
	}

	apiItems := NewItemsAPI(*items)

	return &apiItems, err
}

func (p *Presenter) GetNewInventoryItems() (*ItemsAPI, error) {
	// get base inventory items
	items, err := p.GetInventoryItems()
	if err != nil {
		return nil, err
	}

	// populate them based on latest inventory
	if err = p.populateLatestInvItems(items); err != nil {
		return nil, err
	}
	// populate the latest selected vendor
	p.populateSelectedVendors(items)

	clearItemIds(items)

	return items, nil
}

func (p *Presenter) populateLatestInvItems(items *ItemsAPI) error {
	latestInv, err := p.GetLatestInventory(time.Now())
	if err != nil {
		return err
	}

	j := 0
	for i := 0; i < len(*items) && j < len(latestInv.Items); i++ {
		item := &(*items)[i]
		latestItem := &latestInv.Items[j]
		if item.InventoryItemID == latestItem.InventoryItemID {
			item.Count = latestItem.Count
			item.VendorItemMap = latestItem.VendorItemMap
			// default behavior, will probably be re-set in the next function
			item.SetSelectedVendor(latestItem.SelectedVendor)
		} else if item.Index > latestItem.Index {
			i--
		} else if item.Index < latestItem.Index {
			continue
		}
		j++
	}

	return nil
}

func (p *Presenter) populateSelectedVendors(items *ItemsAPI) {
	for i := 0; i < len(*items); i++ {
		item := &(*items)[i]
		vendor, err := logic.GetLatestVendor(item.InventoryItemID, p.tx)
		if err != nil {
			continue
		}
		fmt.Println(vendor)

		item.SetSelectedVendor(vendor.Name)
	}
}

func clearItemIds(items *ItemsAPI) {
	for _, item := range *items {
		item.ID = ""
	}
}

func (p *Presenter) GetInventoryItem(id string) (*ItemAPI, error) {
	item, err := logic.GetInventoryItem(id, p.tx)
	if err != nil {
		return nil, err
	}

	apiItem := NewItemAPI(item)
	return &apiItem, nil
}

func (p *Presenter) UpdateInventoryItem(item *ItemAPI) (*validate.Errors, error) {
	invItem, err := ConvertToModelInventoryItem(item)
	if err != nil {
		return nil, err
	}

	strVendorID := item.VendorItemMap[item.SelectedVendor].SelectedVendor
	vendorID, err := uuid.FromString(strVendorID)
	if err != nil {
		return nil, err
	}
	vendorItem, err := ConvertToModelVendorItem(*item, vendorID)
	if err != nil {
		return nil, err
	}

	orderID, err := p.getLatestOrderID(item, strVendorID)
	if err != nil {
		return nil, err
	}
	orderItem, err := ConvertToModelOrderItem(*item, orderID)
	if err != nil {
		return nil, err
	}

	verrs, err := logic.UpdateInventoryItem(invItem, p.tx)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	verrs, err = logic.UpdateVendorItem(vendorItem, p.tx)
	if verrs.HasAny() || err != nil {
		return verrs, err
	}

	return logic.UpdateOrderItem(orderItem, p.tx)
}

func (p *Presenter) getLatestOrderID(item *ItemAPI, vendorID string) (uuid.UUID, error) {
	order, err := logic.GetLatestOrder(item.InventoryItemID, vendorID, p.tx)
	if err != nil {
		return uuid.UUID{}, err
	}

	return order.ID, nil
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
