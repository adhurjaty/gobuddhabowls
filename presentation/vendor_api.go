package presentation

// VendorAPI is a struct for serving vendor information to the ui
type VendorAPI struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email,string,omitempty"`
	PhoneNumber  string    `json:"phone_number,string,omitempty"`
	Contact      string    `json:"contact,string,omitempty"`
	ShippingCost float64   `json:"shipping_cost"`
	Items        []ItemAPI `json:"Items"`
}
