package actions

import (
	"buddhabowls/helpers"
	"buddhabowls/logic"
	"buddhabowls/presentation"
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type transactionResult struct {
	TransactionTime time.Time     `json:"created_at"`
	Tax             money         `json:"tax_money"`
	Tip             money         `json:"tip_money"`
	Discount        money         `json:"discount_money"`
	Fee             money         `json:"processing_fee_money"`
	Refund          money         `json:"refunded_money"`
	Items           []itemization `json:"itemizations"`
}

type itemization struct {
	Name   string     `json:"name"`
	Count  string     `json:"quantity"`
	Amount money      `json:"single_quantity_money"`
	Extras []modifier `json:"modifiers"`
}

type modifier struct {
	Name   string `json:"name"`
	Amount money  `json:"applied_money"`
}

type money struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency_code"`
}

type SalesSummary struct {
	Tips    float64     `json:"tips"`
	Fees    float64     `json:"fees"`
	Tax     float64     `json:"tax"`
	Refunds float64     `json:"refunds"`
	Sales   SquareSales `json:"Sales"`
}

type SquareSale struct {
	Name   string  `json:"name"`
	Count  int     `json:"count"`
	Amount float64 `json:"amount"`
}

type SquareSales []SquareSale

func (s SalesSummary) String() string {
	sj, _ := json.Marshal(s)
	return string(sj)
}

const squareURLBase = "https://connect.squareup.com/v1/"

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

	user, err := presenter.GetUser(c.Session().
		Get("current_user_id").(uuid.UUID).String())
	if err != nil {
		return c.Error(404, err)
	}

	sales, err := getSquareSales(user, startTime, endTime)
	if err != nil {
		sales = &SalesSummary{
			Sales: SquareSales{},
		}
		c.Set("errors", "Invalid Square token or location ID")
	}

	c.Set("sales", *sales)

	return c.Render(200, r.HTML("sales/index"))
}

func getTransactionURL(locationID string, timezone string, startTime time.Time, endTime time.Time) string {
	startParam := getDateString(startTime, timezone)
	endParam := getDateString(endTime, timezone)

	return fmt.Sprintf("%s%s/payments?begin_time=%s&end_time=%s",
		squareURLBase, locationID, startParam, endParam)
}

func getDateString(d time.Time, timezone string) string {
	offsetTime := logic.OffsetStart(d)
	location, err := time.LoadLocation(timezone)
	if err != nil {
		location = &time.Location{}
	}
	locationTime := time.Date(offsetTime.Year(), offsetTime.Month(),
		offsetTime.Day(), offsetTime.Hour(), offsetTime.Minute(),
		offsetTime.Second(), offsetTime.Nanosecond(), location)
	return helpers.RFC3339Date(locationTime)
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

func getSquareSales(user *presentation.UserAPI, startTime time.Time, endTime time.Time) (*SalesSummary, error) {

	transactionURL := getTransactionURL(user.SquareLocation, user.Timezone,
		startTime, endTime)

	resp, err := sendGetRequest(transactionURL, user.SquareToken)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	results := &[]transactionResult{}
	err = json.Unmarshal(jsonBytes, results)
	if err != nil {
		return nil, err
	}

	salesSummary := &SalesSummary{
		Sales: SquareSales{},
	}
	for _, result := range *results {
		for _, item := range result.Items {
			sale := itemizationToSale(item)
			addItemToSummary(sale, &salesSummary.Sales)

			for _, extra := range item.Extras {
				sale = extraToSale(extra)
				addItemToSummary(sale, &salesSummary.Sales)
			}
		}

		// fees are negative. want a positive number
		salesSummary.Fees -= float64(result.Fee.Amount) / 100.0
		salesSummary.Tax += float64(result.Tax.Amount) / 100.0
		salesSummary.Refunds += float64(result.Refund.Amount) / 100.0
		salesSummary.Tips += float64(result.Tip.Amount) / 100.0
	}

	return salesSummary, nil
}

func itemizationToSale(item itemization) SquareSale {
	count, err := strconv.ParseFloat(item.Count, 64)
	if err != nil {
		count = 1
	}
	return SquareSale{
		Name:   item.Name,
		Amount: float64(item.Amount.Amount) / 100.0,
		Count:  int(count),
	}
}

func extraToSale(extra modifier) SquareSale {
	return SquareSale{
		Name:   "Extra: " + extra.Name,
		Amount: float64(extra.Amount.Amount) / 100.0,
		Count:  1,
	}
}

func addItemToSummary(sale SquareSale, allSales *SquareSales) {
	for i, prevSale := range *allSales {
		if sale.Name == prevSale.Name {
			(*allSales)[i].Count += sale.Count
			(*allSales)[i].Amount += sale.Amount
			return
		}
	}

	*allSales = append(*allSales, sale)
}
