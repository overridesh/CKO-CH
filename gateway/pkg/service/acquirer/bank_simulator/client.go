package bank_simulator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/overridesh/checkout-challenger/pkg/service/acquirer"
)

// bankSimulatorResponse request
type bankSimulatorRequest struct {
	Amount      int32  `json:"amount"`
	Currency    string `json:"currency"`
	ExpiryMonth string `json:"expiry_month"`
	ExpiryYear  string `json:"expiry_year"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Number      string `json:"number"`
}

// bankSimulatorResponse response from acquirer
type bankSimulatorResponse struct {
	Amount   int32               `json:"amount"`
	Currency string              `json:"currency"`
	Code     string              `json:"code"`
	Message  string              `json:"message"`
	Status   string              `json:"status"`
	Source   bankSimulatorSource `json:"source"`
}

// Source payment method
type bankSimulatorSource struct {
	Bin         string `json:"bin"`
	CardType    string `json:"card_type"`
	ExpiryMonth string `json:"expiry_month"`
	ExpiryYear  string `json:"expiry_year"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Number      string `json:"number"`
}

const (
	timeout          int = 65
	transportTimeout int = 10
)

type bankSimulator struct {
	client  *http.Client
	baseURL string
	apikey  string
}

func New(baseURL, apikey string) acquirer.Acquirer {
	// add a timeout for SSL handshakes and dialer
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(transportTimeout) * time.Second,
		}).Dial,
		TLSHandshakeTimeout: time.Duration(transportTimeout) * time.Second,
	}

	// add timeout to overall request
	client := http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: netTransport,
	}

	return &bankSimulator{
		client:  &client,
		baseURL: baseURL,
		apikey:  apikey,
	}
}

func (bs *bankSimulator) Purchase(payload *acquirer.PaymentRequest) (*acquirer.PaymentResponse, error) {
	var (
		method   string = http.MethodPost
		endpoint string = "/transactions"
	)

	rawData, err := json.Marshal(bankSimulatorRequest{
		Amount:      payload.Amount,
		Currency:    payload.Currency,
		ExpiryMonth: payload.ExpiryMonth,
		ExpiryYear:  payload.ExpiryYear,
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Number:      payload.Number,
	})
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(rawData)

	request, err := http.NewRequest(method, fmt.Sprintf("%s%s", bs.baseURL, endpoint), body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("X-Auth-Token", bs.apikey)

	req, err := bs.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer req.Body.Close()

	var response bankSimulatorResponse
	if err := json.NewDecoder(req.Body).Decode(&response); err != nil {
		return nil, err
	}

	if req.StatusCode == http.StatusCreated {
		return &acquirer.PaymentResponse{
			StatusCode:  req.StatusCode,
			Amount:      response.Amount,
			Currency:    response.Currency,
			ExpiryMonth: response.Source.ExpiryMonth,
			ExpiryYear:  response.Source.ExpiryYear,
			FirstName:   response.Source.FirstName,
			LastName:    response.Source.LastName,
			Number:      response.Source.Number,
			Code:        response.Code,
			Summary:     response.Message,
			Status:      response.Status,
			CardBin:     response.Source.Bin,
			CardType:    response.Source.CardType,
		}, nil
	}

	return &acquirer.PaymentResponse{
		Code:       response.Code,
		Summary:    response.Message,
		Status:     response.Status,
		StatusCode: req.StatusCode,
	}, errors.New(response.Message)
}
