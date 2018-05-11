package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/validate/validators"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type PurchaseOrder struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	VendorID     uuid.UUID  `json:"vendor_id" db:"vendor_id"`
	Vendor       Vendor     `belongs_to:"vendors" db:"-"`
	OrderDate    nulls.Time `json:"order_date" db:"order_date"`
	ReceivedDate nulls.Time `json:"received_date" db:"received_date"`
	ShippingCost float64    `json:"shipping_cost" db:"shipping_cost"`
	Items        OrderItems `has_many:"order_items" db:"-" fk_id:"order_id"`
}

// String is not required by pop and may be deleted
func (p PurchaseOrder) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// PurchaseOrders is not required by pop and may be deleted
type PurchaseOrders []PurchaseOrder

// String is not required by pop and may be deleted
func (p PurchaseOrders) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *PurchaseOrder) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		validate.ValidatorFunc(func(errors *validate.Errors) {
			if !p.OrderDate.Valid {
				errors.Add(validators.GenerateKey(p.Vendor.Name), "Must have order date")
			}
		}),
		validate.ValidatorFunc(func(errors *validate.Errors) {
			if p.OrderDate.Valid && p.ReceivedDate.Valid && p.OrderDate.Time.Unix() > p.ReceivedDate.Time.Unix() {
				errors.Add(validators.GenerateKey(p.Vendor.Name+" "+p.OrderDate.Time.String()),
					"Received date must be after order date")
			}
		}),
		validate.ValidatorFunc(func(errors *validate.Errors) {
			if p.Items == nil || len(p.Items) == 0 {
				errors.Add(validators.GenerateKey(p.Vendor.Name+" Items"),
					"Must have items with count > 0")
			}
		}),
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *PurchaseOrder) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *PurchaseOrder) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// GetCost gets the total cost of the purchase order
func (p PurchaseOrder) GetCost() float64 {
	var cost float64
	for _, item := range p.Items {
		cost += item.Price * item.Count
	}

	return cost + p.ShippingCost
}

// GetCategoryCosts gets a map of category -> cost map for the order
func (p PurchaseOrder) GetCategoryCosts() CategoryBreakdown {
	catCosts := map[InventoryItemCategory]float64{}

	for _, item := range p.Items {
		// fmt.Println(item)
		catCosts = AddToCategoryMap(catCosts, item.InventoryItem.Category, item.Price*item.Count)
	}

	return FromCategoryMap(catCosts)
}

// LoadPurchaseOrder gets purchase order and sub-components matching the given ID
func LoadPurchaseOrder(tx *pop.Connection, id string) (PurchaseOrder, error) {
	po := PurchaseOrder{}

	if err := tx.Eager().Find(&po, id); err != nil {
		return po, err
	}

	for i := 0; i < len(po.Items); i++ {
		if err := tx.Eager().Find(&po.Items[i], po.Items[i].ID); err != nil {
			return po, err
		}
		if err := tx.Eager().Find(&po.Items[i].InventoryItem, po.Items[i].InventoryItemID); err != nil {
			return po, err
		}
	}

	return po, nil
}

// LoadPurchaseOrders gets the purchase orders with the specified query
// including all sub-components
func LoadPurchaseOrders(q *pop.Query) (*PurchaseOrders, error) {
	poList := &PurchaseOrders{}

	if err := q.All(poList); err != nil {
		return nil, err
	}

	// I don't love the fact that I need to load the nested models manually
	// TODO: look for a solution to eager loading nested objects
	for _, po := range *poList {
		for i := 0; i < len(po.Items); i++ {
			if err := q.Connection.Eager().Find(&po.Items[i], po.Items[i].ID); err != nil {
				return nil, err
			}
			if err := q.Connection.Eager().Find(&po.Items[i].InventoryItem, po.Items[i].InventoryItemID); err != nil {
				return nil, err
			}
		}

		po.Items.Sort()
	}

	return poList, nil
}

// GetYears gets the years for which there is company data
func GetYears(tx *pop.Connection) ([]int, error) {
	yearResult := make([]int, 50)
	// Search for just the years in purchase orders
	q := tx.RawQuery("SELECT DISTINCT EXTRACT(YEAR FROM order_date) FROM purchase_orders ORDER BY EXTRACT(YEAR FROM order_date) ASC")

	if err := q.All(&yearResult); err != nil {
		return nil, err
	}

	// throw away extra allocated data. Probably a better way to do this
	years := []int{}
	for _, val := range yearResult {
		if val > 2000 {
			years = append(years, val)
		}
	}

	return years, nil
}
