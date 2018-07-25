package presentation_test

import (
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"testing"
)

func Test_GetPurchaseOrders(t *testing.T) {

}

func Test_GetPurchaseOrder(t *testing.T) {

}

func Test_InsertPurchaseOrder(t *testing.T) {

}

func Test_UpdatePurchaseOrder(t *testing.T) {

}

func Test_DestroyPurchaseOrder(t *testing.T) {

}

func Test_GetVendors(t *testing.T) {

}

func Test_GetVendor(t *testing.T) {

}

func Test_GetPeriods(t *testing.T) {

}

func Test_GetSelectedPeriod(t *testing.T) {

}

func Test_GetWeeks(t *testing.T) {

}

func Test_GetSelectedWeek(t *testing.T) {

}

func Test_GetYears(t *testing.T) {

}

func getTX() (*pop.Connection, error) {
	return pop.Connect(envy.Get("GO_ENV", "test"))
}
