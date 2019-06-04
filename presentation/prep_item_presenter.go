package presentation

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
