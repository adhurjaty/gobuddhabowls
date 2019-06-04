package presentation

import (
	"buddhabowls/logic"
	"time"
)

func (p *Presenter) GetNewPrepItems() (*ItemsAPI, error) {
	items, err := p.GetPrepItems()
	if err != nil {
		return nil, err
	}

	if err = p.populateLatestPrepItems(items); err != nil {
		return nil, err
	}

	clearItemIds(items)

	return items, nil
}

func (p *Presenter) GetPrepItems() (*ItemsAPI, error) {
	items, err := logic.GetPrepItems(p.tx)
	if err != nil {
		return nil, err
	}

	apiItems := NewItemsAPI(items)

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
			item.BatchRecipe = latestItem.BatchRecpie
		}
	}

	return nil
}
