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

// ConvertToAPI converts a purchase order to an api purchase order
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

// ConvertToAPI converts a purchase order slice to an api purchase order slice
func (v *VendorsAPI) ConvertToAPI(m interface{}) error {
	vendors, ok := m.(models.Vendors)
	if !ok {
		return errors.New("Must supply Vendors type")
	}

	apis := VendorsAPI{}
	for i, vendor := range vendors {
		api := VendorAPI{}
		if err := api.ConvertToAPI(vendor); err != nil {
			return err
		}
		apis = append(apis, api)
	}

	v = &apis
	return nil
}
