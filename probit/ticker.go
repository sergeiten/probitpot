package probit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Ticker struct {
	Last        string `json:"last"`
	Low         string `json:"low"`
	High        string `json:"high"`
	Change      string `json:"change"`
	BaseVolume  string `json:"base_volume"`
	QuoteVolume string `json:"quote_volume"`
	MarketID    string `json:"market_id"`
	Time        string `json:"time"`
}

type TickerResponse struct {
	Data []Ticker `json:"data"`
}

func (p *Probit) Ticker(marketID string) (*TickerResponse, error) {
	req, err := p.tickerRequest(marketID)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("response status code is %d", resp.StatusCode))
	}

	var tickerResponse TickerResponse

	err = json.NewDecoder(resp.Body).Decode(&tickerResponse)
	if err != nil {
		return nil, err
	}

	return &tickerResponse, nil
}

func (p *Probit) tickerRequest(marketID string) (*http.Request, error) {
	u := fmt.Sprintf(ApiURL+"/api/exchange/v1/ticker?market_ids=%s", marketID)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+p.accessToken)

	return req, nil
}
