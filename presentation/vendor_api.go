package presentation

import (
	"buddhabowls/models"
	"errors"
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

// ConvertToAPI converts a vendor to an api vendor
func (v *VendorAPI) ConvertToAPI(m interface{}) error {
	vendor, ok := m.(models.Vendor)
	if !ok {
		return errors.New("Must supply Vendor type")
	}

	v.ID = vendor.ID.String()
	v.Name = vendor.Name
	v.Email = vendor.Email.String
	v.PhoneNumber = vendor.PhoneNumber.String
	v.Contact = vendor.Contact.String
	v.ShippingCost = vendor.ShippingCost
	v.Items = ItemsAPI{}
	v.Items.ConvertToAPI(vendor.Items)

	return nil
}

// ConvertToAPI converts a vendor slice to an api vendor slice
func (v *VendorsAPI) ConvertToAPI(m interface{}) error {
	vendors, ok := m.(models.Vendors)
	if !ok {
		return errors.New("Must supply Vendors type")
	}

	apis := VendorsAPI{}
	for _, vendor := range vendors {
		api := VendorAPI{}
		if err := api.ConvertToAPI(vendor); err != nil {
			return err
		}
		apis = append(apis, api)
	}

	v = &apis
	return nil
}
