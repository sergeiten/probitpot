package probit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const TimeInForceGTC = "gtc"     // Good Till Cancel
const TimeInForceIOC = "ioc"     // Immediate Or Cancel
const TimeInForceFok = "fok"     // Fill Or Fill
const TimeInForceGTCPO = "gtcpo" // Good Till Cancel and Post Only

const SideSell = "sell"
const SideBuy = "buy"

const TypeLimit = "limit"
const TypeMarket = "market"

type NewOrderParams struct {
	MarketID      string `json:"market_id"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	TimeInForce   string `json:"time_in_force"`
	LimitPrice    string `json:"limit_price,omitempty"`
	Cost          string `json:"cost,omitempty"`
	Quantity      string `json:"quantity"`
	ClientOrderID string `json:"client_order_id,omitempty"`
}

type NewOrderResponse struct {
	Data struct {
		ID                string `json:"id"`
		UserID            string `json:"user_id"`
		MarketID          string `json:"market_id"`
		Side              string `json:"side"`
		Type              string `json:"type"`
		Quantity          string `json:"quantity"`
		LimitPrice        string `json:"limit_price"`
		TimeInForce       string `json:"time_in_force"`
		FilledCost        string `json:"filled_cost"`
		FilledQuantity    string `json:"filled_quantity"`
		CancelledQuantity string `json:"cancelled_quantity"`
		OpenQuantity      string `json:"open_quantity"`
		Status            string `json:"status"`
		Time              string `json:"time"`
		ClientOrderID     string `json:"client_order_id"`
	} `json:"data"`
}

func (p *Probit) Sell(marketID, ttype, limitPrice, quantity string) (*NewOrderResponse, error) {
	params := NewOrderParams{
		MarketID:    marketID,
		Type:        ttype,
		Side:        SideSell,
		TimeInForce: TimeInForceGTC,
		LimitPrice:  limitPrice,
		Quantity:    quantity,
	}

	return p.newOrderDoRequest(params)
}

func (p *Probit) Buy(marketID, ttype, limitPrice, quantity, clientOrderID string) (*NewOrderResponse, error) {
	params := NewOrderParams{
		MarketID:      marketID,
		Type:          ttype,
		Side:          SideBuy,
		TimeInForce:   TimeInForceGTC,
		LimitPrice:    limitPrice,
		Quantity:      quantity,
		ClientOrderID: clientOrderID,
	}

	return p.newOrderDoRequest(params)
}

func (p *Probit) newOrderDoRequest(params NewOrderParams) (*NewOrderResponse, error) {
	requestBody, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	req, err := p.newOrderRequest(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var respError responseError

		_ = json.NewDecoder(resp.Body).Decode(&respError)

		return nil, errors.New(fmt.Sprintf("response status code: %d, error code: %s, request body: %s", resp.StatusCode, respError.ErrorCode, requestBody))
	}

	var newOrderResponse NewOrderResponse

	err = json.NewDecoder(resp.Body).Decode(&newOrderResponse)
	if err != nil {
		return nil, err
	}

	return &newOrderResponse, nil
}

func (p *Probit) newOrderRequest(reqBody []byte) (*http.Request, error) {

	req, err := http.NewRequest("POST", ApiURL+"/api/exchange/v1/new_order", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+p.accessToken)

	return req, nil
}
