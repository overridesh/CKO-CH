package bank_simulator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/overridesh/checkout-challenger/pkg/service/acquirer"
)

func TestBankSimulatorSuccess(t *testing.T) {
	mock := bankSimulatorResponse{
		Amount:   1000,
		Currency: "USD",
		Code:     "00",
		Message:  "ok",
		Status:   "approved",
		Source: bankSimulatorSource{
			Bin:         "4444",
			CardType:    "credit_card",
			ExpiryMonth: "20",
			ExpiryYear:  "2022",
			FirstName:   "John",
			LastName:    "Doe",
			Number:      "XXXX-XXXX-XXXX-5555",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(mock)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "internal server error")
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, string(bytes))
	}))

	defer ts.Close()

	var (
		bankSimulator acquirer.Acquirer = New(ts.URL, "")
	)

	response, err := bankSimulator.Purchase(&acquirer.PaymentRequest{})
	if err != nil {
		t.Errorf("expected a error nil, but got %v", err)
	}

	expected := acquirer.PaymentResponse{
		Amount:      mock.Amount,
		Currency:    mock.Currency,
		ExpiryMonth: mock.Source.ExpiryMonth,
		ExpiryYear:  mock.Source.ExpiryYear,
		FirstName:   mock.Source.FirstName,
		LastName:    mock.Source.LastName,
		Number:      mock.Source.Number,
		Code:        mock.Code,
		Summary:     mock.Message,
		CardBin:     mock.Source.Bin,
		CardType:    mock.Source.CardType,
		Status:      mock.Status,
		StatusCode:  http.StatusCreated,
	}

	if !cmp.Equal(response, &expected) {
		t.Errorf("expected items are equals, response: %v, expected: %v", response, expected)
	}
}

func TestBankSimulatorUnprocessEntity(t *testing.T) {
	mock := bankSimulatorResponse{
		Code:    "1000",
		Message: "cannot something",
		Status:  "failed",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(mock)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "internal server error")
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, string(bytes))
	}))

	defer ts.Close()

	var (
		bankSimulator acquirer.Acquirer = New(ts.URL, "")
	)

	response, err := bankSimulator.Purchase(&acquirer.PaymentRequest{})
	if err.Error() != mock.Message {
		t.Errorf("expected a error %v, but got %s", err, mock.Message)
	}

	expected := acquirer.PaymentResponse{
		Amount:      mock.Amount,
		Currency:    mock.Currency,
		ExpiryMonth: mock.Source.ExpiryMonth,
		ExpiryYear:  mock.Source.ExpiryYear,
		FirstName:   mock.Source.FirstName,
		LastName:    mock.Source.LastName,
		Number:      mock.Source.Number,
		Code:        mock.Code,
		Summary:     mock.Message,
		CardBin:     mock.Source.Bin,
		CardType:    mock.Source.CardType,
		Status:      mock.Status,
		StatusCode:  http.StatusUnprocessableEntity,
	}

	if !cmp.Equal(response, &expected) {
		t.Errorf("expected items are equals, response: %v, expected: %v", response, expected)
	}
}

func TestBankSimulatorWithoutURL(t *testing.T) {
	var (
		bankSimulator acquirer.Acquirer = New("", "")
	)

	response, err := bankSimulator.Purchase(&acquirer.PaymentRequest{})
	if response != nil {
		t.Errorf("expected a response nil, but got %v", response)
	}
	if err == nil {
		t.Error("expected a error not nil, but got nil")
	}
}
