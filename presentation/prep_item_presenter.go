package presentation

import (
	"buddhabowls/logic"
	"fmt"
	"time"

	"github.com/gobuffalo/validate"
)

func (p *Presenter) GetNewPrepItems() (*ItemsAPI, error) {
	items, err := p.GetPrepItems()
	if err != nil {
		return nil, err
	}

	if err = p.populateLatestPrepItems(items); err != nil {
		return nil, err
	}

	return items, nil
}

func (p *Presenter) GetPrepItems() (*ItemsAPI, error) {
	items, err := logic.GetPrepItems(p.tx)
	if err != nil {
		return nil, err
	}

	apiItems := NewItemsAPI(items)
	err = p.populatePrepItemCosts(&apiItems)

	fmt.Println(apiItems)
	return &apiItems, err
}

func (p *Presenter) populateLatestPrepItems(items *ItemsAPI) error {
	latestInv, err := p.GetLatestInventory(time.Now())
	if err != nil {
		return err
	}

	for i := 0; i < len(*items); i++ {
		item := &(*items)[i]
		for _, latestItem := range latestInv.PrepItems {
			if item.BatchRecipeID == latestItem.BatchRecipeID {
				item.Count = latestItem.Count
				break
			}
		}
	}

	return nil
}

func (p *Presenter) populatePrepItemCosts(items *ItemsAPI) error {
	for i, item := range *items {
		cost, err := p.getItemRecipeCost(item)
		if err != nil {
			return err
		}
		(*items)[i].Price = cost * item.Conversion
	}

	return nil
}

func (p *Presenter) GetPrepItem(id string) (*ItemAPI, error) {
	item, err := logic.GetPrepItem(id, p.tx)
	if err != nil {
		return nil, err
	}

	apiItem := NewItemAPI(item)
	return &apiItem, nil
}

func (p *Presenter) UpdatePrepItem(item *ItemAPI) (*validate.Errors, error) {
	prepItem, err := ConvertToModelPrepItem(item)
	if err != nil {
		return validate.NewErrors(), err
	}

	return logic.UpdatePrepItem(prepItem, p.tx)
}
