package presentation

import (
	"buddhabowls/models"
	"encoding/json"
	"github.com/gobuffalo/uuid"
)

// VendorAPI is a struct for serving vendor information to the ui
type VendorAPI struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email,string,omitempty"`
	PhoneNumber  string   `json:"phone_number,string,omitempty"`
	Contact      string   `json:"contact,string,omitempty"`
	ShippingCost float64  `json:"shipping_cost"`
	Items        ItemsAPI `json:"Items"`
}

type VendorsAPI []VendorAPI

func (v VendorAPI) String() string {
	jo, _ := json.Marshal(v)
	return string(jo)
}

func (v VendorsAPI) String() string {
	jo, _ := json.Marshal(v)
	return string(jo)
}

// NewVendorAPI converts a vendor to an api vendor
func NewVendorAPI(vendor *models.Vendor) VendorAPI {
	return VendorAPI{
		ID:           vendor.ID.String(),
		Name:         vendor.Name,
		Email:        vendor.Email.String,
		PhoneNumber:  vendor.PhoneNumber.String,
		Contact:      vendor.Contact.String,
		ShippingCost: vendor.ShippingCost,
		Items:        NewItemsAPI(&vendor.Items),
	}
}

// NewVendorsAPI converts a vendor slice to an api vendor slice
func NewVendorsAPI(vendors *models.Vendors) VendorsAPI {
	apis := make([]VendorAPI, len(*vendors))
	for i, vendor := range *vendors {
		apis[i] = NewVendorAPI(&vendor)
	}

	return apis
}

// SelectValue returns the ID for select input tags
func (v VendorAPI) SelectValue() interface{} {
	return v.ID
}

// SelectLabel returs the name for select input tags
func (v VendorAPI) SelectLabel() string {
	if v.ID == "" {
		return "- Select a vendor -"
	}
	return v.Name
}

func ConvertToModelVendor(vendAPI *VendorAPI) (*models.Vendor, error) {
	id := uuid.UUID{}
	if len(vendAPI.ID) > 0 {
		var err error
		id, err = uuid.FromString(vendAPI.ID)
		if err != nil {
			return nil, err
		}
	}

	items, err := ConvertToModelVendorItems(vendAPI.Items, id)
	if err != nil {
		return nil, err
	}

	return &models.Vendor{
		ID:           id,
		Name:         vendAPI.Name,
		Email:        models.StringToNullString(vendAPI.Email),
		Contact:      models.StringToNullString(vendAPI.Contact),
		PhoneNumber:  models.StringToNullString(vendAPI.PhoneNumber),
		ShippingCost: vendAPI.ShippingCost,
		Items:        *items,
	}, nil
}
