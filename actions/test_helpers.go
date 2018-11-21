package actions

import (
	"buddhabowls/models"
	"database/sql"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"time"
)

func login(as *ActionSuite) {
	db := as.DB
	session := as.Session
	user := &models.User{
		Email:        "some@email.com",
		PasswordHash: "lksjf;lej;lkjfsa;eaaklesmae;fmae",
	}
	db.Create(user)

	dbUser := &models.User{}
	db.First(dbUser)
	session.Set("current_user_id", dbUser.ID.String())
}

func createPO(db *pop.Connection, orderTime time.Time) (*models.PurchaseOrder, error) {
	vendor, err := createVendor(db)
	if err != nil {
		return nil, err
	}

	purchaseOrder := &models.PurchaseOrder{
		Vendor:    *vendor,
		VendorID:  vendor.ID,
		OrderDate: nulls.Time{Time: orderTime, Valid: true},
	}
	if err = db.Create(purchaseOrder); err != nil {
		return nil, err
	}

	item := &models.OrderItem{
		InventoryItem:   vendor.Items[0].InventoryItem,
		InventoryItemID: vendor.Items[0].InventoryItemID,
		Count:           4,
		Price:           vendor.Items[0].Price,
		OrderID:         purchaseOrder.ID,
	}

	if err = db.Create(item); err != nil {
		return nil, err
	}

	purchaseOrder.Items = models.OrderItems{*item}

	return purchaseOrder, nil
}

func createPOMultipleItems(db *pop.Connection, orderTime time.Time) (*models.PurchaseOrder, error) {
	purchaseOrder, err := createPO(db, orderTime)
	if err != nil {
		return nil, err
	}

	newItem, err := createVendorItem(db, &purchaseOrder.Vendor, "yet_another_item")
	item := &models.OrderItem{
		InventoryItem:   newItem.InventoryItem,
		InventoryItemID: newItem.InventoryItemID,
		Count:           6,
		Price:           newItem.Price,
		OrderID:         purchaseOrder.ID,
	}

	if err = db.Create(item); err != nil {
		return nil, err
	}

	purchaseOrder.Items = append(purchaseOrder.Items, *item)

	return purchaseOrder, nil
}

func receiveOrder(db *pop.Connection, purchaseOrder *models.PurchaseOrder, date time.Time) error {
	purchaseOrder.ReceivedDate = nulls.Time{Time: date, Valid: true}
	return db.Update(purchaseOrder)
}

func createVendor(db *pop.Connection) (*models.Vendor, error) {
	vendor := &models.Vendor{
		Name:         "test_vendor",
		ShippingCost: 2.5,
	}

	err := db.Create(vendor)
	if err != nil {
		return nil, err
	}

	invItem, err := createInventoryItem(db, "test_item")
	if err != nil {
		return nil, err
	}

	item := &models.VendorItem{
		InventoryItem:   *invItem,
		InventoryItemID: invItem.ID,
		Vendor:          *vendor,
		VendorID:        vendor.ID,
		Conversion:      1,
		Price:           5,
		PurchasedUnit:   sql.NullString{String: "EA", Valid: true},
	}
	if err = db.Create(item); err != nil {
		return nil, err
	}

	vendor.Items = models.VendorItems{*item}

	return vendor, nil
}

func createInventoryItem(db *pop.Connection, name string) (*models.InventoryItem, error) {
	category := &models.InventoryItemCategory{
		Name:  "test_category",
		Index: 0,
	}
	err := db.Create(category)
	if err != nil {
		return nil, err
	}
	item := &models.InventoryItem{
		Name:       name,
		Category:   *category,
		CategoryID: category.ID,
		IsActive:   true,
		Index:      0,
	}
	err = db.Create(item)

	return item, err
}

func createVendorItem(db *pop.Connection, vendor *models.Vendor, name string) (*models.VendorItem, error) {
	invItem, err := createInventoryItem(db, name)
	if err != nil {
		return nil, err
	}

	item := &models.VendorItem{
		InventoryItem:   *invItem,
		InventoryItemID: invItem.ID,
		Vendor:          *vendor,
		VendorID:        vendor.ID,
		Conversion:      4,
		Price:           8.99,
		PurchasedUnit:   sql.NullString{String: "Quart", Valid: true},
	}
	if err = db.Create(item); err != nil {
		return nil, err
	}

	return item, nil
}
