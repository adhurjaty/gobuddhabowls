package actions

import (
	"buddhabowls/helpers"
	"buddhabowls/presentation"
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type transactionResult struct {
	Errors       []apiError    `json:"errors,omitempty"`
	Transactions []transaction `json:"transactions"`
	Cursor       string        `json:"cursor,omitempty"`
}

type apiError struct {
	Field  string `json:"field"`
	Detail string `json:"detail`
}

type transaction struct {
	TransactionTime time.Time `json:"created_at"`
	Tenders         []tender  `json:"tenders,omitempty"`
	Refunds         []refund  `json:"refunds,omitempty"`
}

type tender struct {
	Amount money `json:"amount_money"`
	Tip    money `json:"tip_money"`
	Fee    money `json:"processing_fee_money"`
}

type money struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type refund struct {
	Amount money `json:"amount_money"`
	Fee    money `json:"processing_fee_money"`
}

type SquareSale struct {
	TransactionTime time.Time
	Amount          float64
	Tip             float64
	Fee             float64
}

type SquareSales []SquareSale

func (s SquareSales) String() string {
	sj, _ := json.Marshal(s)
	return string(sj)
}

const squareURLBase = "https://connect.squareup.com/v2/locations/"

func ListSales(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	presenter := presentation.NewPresenter(tx)
	startTime, endTime, err := setPeriodSelector(c, presenter)
	if err != nil {
		return errors.WithStack(err)
	}

	sales, err := getSquareSales(startTime, endTime)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Set("sales", *sales)

	return c.Render(200, r.HTML("sales/index"))
}

func getTransactionURL(locationID string, startTime time.Time, endTime time.Time) string {
	startParam := helpers.RFC3339Date(startTime)
	endParam := helpers.RFC3339Date(endTime)

	return fmt.Sprintf("%s%s/transactions?begin_time=%s&end_time=%s",
		squareURLBase, locationID, startParam, endParam)
}

func addCursorToURL(url string, cursor string) string {
	if cursor == "" {
		return url
	}
	return url + "&cursor=" + cursor
}

func sendGetRequest(url string, token string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	return client.Do(req)
}

func getSquareSales(startTime time.Time, endTime time.Time) (*SquareSales, error) {

	locationID := "69VJ030ANYAGV"
	// remember to change the token when pushing to remote
	squareToken := "sq0atp-Zo5ieRMqg6UpcSsAzSLEJQ"
	cursor := ""
	allTransactions := &[]transaction{}

	for {
		transactionURL := getTransactionURL(locationID, startTime, endTime)
		transactionURL = addCursorToURL(transactionURL, cursor)

		resp, err := sendGetRequest(transactionURL, squareToken)
		if err != nil {
			return nil, err
		}
		jsonBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		results := &transactionResult{}
		json.Unmarshal(jsonBytes, results)
		if results.Errors != nil {
			return nil, errors.New(fmt.Sprintf("%s: %s",
				results.Errors[0].Field, results.Errors[0].Detail))
		}

		*allTransactions = append(*allTransactions, results.Transactions...)

		cursor = results.Cursor
		if cursor == "" {
			break
		}
	}

	allSales := toSales(allTransactions)

	return allSales, nil
}

func toSales(transactions *[]transaction) *SquareSales {
	sales := &SquareSales{}
	for _, transaction := range *transactions {
		sale := SquareSale{
			TransactionTime: transaction.TransactionTime,
		}

		for _, tender := range transaction.Tenders {
			sale.Amount += float64(tender.Amount.Amount) / 100.0
			sale.Fee += float64(tender.Fee.Amount) / 100.0
			sale.Tip += float64(tender.Tip.Amount) / 100.0
		}

		for _, refund := range transaction.Refunds {
			sale.Amount -= float64(refund.Amount.Amount) / 100.0
			sale.Fee += float64(refund.Fee.Amount) / 100.
		}

		*sales = append(*sales, sale)
	}

	return sales
}
