package actions

import (
	"buddhabowls/presentation"
	"encoding/json"
)

func getItemsFromParams(itemsParamJSON string) (presentation.ItemsAPI, error) {
	items := presentation.ItemsAPI{}

	err := json.Unmarshal([]byte(itemsParamJSON), &items)
	if err != nil {
		return items, err
	}

	return items, nil
}
