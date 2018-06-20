package presentation

import (
	"buddhabowls/models"
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

// NewVendorAPI converts a vendor to an api vendor
func NewVendorAPI(vendor *models.Vendor) VendorAPI {
	return VendorAPI{
		ID:           vendor.ID.String(),
		Name:         vendor.Name,
		Email:        vendor.Email.String,
		PhoneNumber:  vendor.PhoneNumber.String,
		Contact:      vendor.Contact.String,
		ShippingCost: vendor.ShippingCost,
		Items:        NewItemsAPI(vendor.Items),
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
