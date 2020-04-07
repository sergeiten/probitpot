package probit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (p *Probit) Token() error {
	req, err := p.tokenRequest()
	if err != nil {
		return err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("response status code is %d", resp.StatusCode))
	}

	var tokenResponse tokenResponse

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return err
	}

	p.accessToken = tokenResponse.AccessToken

	return nil
}

func (p *Probit) tokenRequest() (*http.Request, error) {
	requestBody, err := json.Marshal(map[string]string{
		"grant_type": "client_credentials",
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", TokenURL+"/token", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+p.basicToken)

	return req, nil
}
