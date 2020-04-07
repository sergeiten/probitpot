package probit

import (
	"encoding/base64"
	"errors"
	"net/http"
	"time"
)

const (
	ApiURL        = "https://api.probit.com"
	TokenURL      = "https://accounts.probit.com"
	MinTotalPrice = 1000
)

type responseError struct {
	ErrorCode string            `json:"errorCode"`
	Message   string            `json:"message"`
	Details   map[string]string `json:"details"`
}

type Probit struct {
	client          *http.Client
	clientID        string
	clientSecretKey string
	basicToken      string
	accessToken     string
}

func NewProbit(clientID, clientSecretKey string) (*Probit, error) {
	if clientID == "" {
		return nil, errors.New("client id is empty")
	}

	if clientSecretKey == "" {
		return nil, errors.New("client secret key is empty")
	}

	basicToken := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecretKey))

	client := newTimeoutClient(30*time.Second, 30*time.Second)

	return &Probit{
		client:          client,
		clientID:        clientID,
		clientSecretKey: clientSecretKey,
		basicToken:      basicToken,
	}, nil
}
