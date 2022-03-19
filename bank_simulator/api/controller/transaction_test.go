package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	creditcard "github.com/durango/go-credit-card"
	"github.com/gin-gonic/gin"
)

var (
	xAuthTokenFake string = "Fake-Key"
)

func TestTransactionRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Router with gin
	router := gin.Default()
	router.POST("/transactions", TransactionControllerFunc)

	tests := []testCase{
		{
			name:       "transaction_unauthorized",
			path:       "/transactions",
			method:     http.MethodPost,
			statusCode: http.StatusUnauthorized,
			expect:     nil,
		},
		{
			name:       "transaction_invalid_number",
			path:       "/transactions",
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			headers: map[string]string{
				"X-Auth-Token": xAuthTokenFake,
			},
			input: `{
				"amount": 1000,
				"currency": "USD",
				"expiry_month": "10",
				"expiry_year": "2022",
				"first_name": "John",
				"last_name": "Doe",
				"number": "444"
			}`,
			expect: nil,
		},
		{
			name:       "transaction_invalid_number_luhn_algorithm",
			path:       "/transactions",
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			headers: map[string]string{
				"X-Auth-Token": xAuthTokenFake,
			},
			input: `{
				"amount": 1000,
				"currency": "USD",
				"expiry_month": "10",
				"expiry_year": "2022",
				"first_name": "John",
				"last_name": "Doe",
				"number": "79927398713"
			}`,
			expect: nil,
		},
		{
			name:       "transaction_wrong_expiry_date",
			path:       "/transactions",
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			headers: map[string]string{
				"X-Auth-Token": xAuthTokenFake,
			},
			input: `{
				"amount": 1000,
				"currency": "USD",
				"expiry_month": "99",
				"expiry_year": "2000",
				"first_name": "John",
				"last_name": "Doe",
				"number": "345678901234564"
			}`,
			expect: nil,
		},
		{
			name:       "transaction_wrong_amount",
			path:       "/transactions",
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			headers: map[string]string{
				"X-Auth-Token": xAuthTokenFake,
			},
			input: `{
				"amount": "ASD",
				"currency": "USD",
				"expiry_month": "99",
				"expiry_year": "2000",
				"first_name": "John",
				"last_name": "Doe",
				"number": "345678901234564"
			}`,
			expect: nil,
		},
		{
			name:       "transaction_wrong_amount",
			path:       "/transactions",
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			headers: map[string]string{
				"X-Auth-Token": xAuthTokenFake,
			},
			input: `{
				"amount": "ASD",
				"currency": "USD",
				"expiry_month": "99",
				"expiry_year": "2000",
				"first_name": "John",
				"last_name": "Doe",
			}`,
			expect: nil,
		},
	}

	for cardNumber, value := range TestCards {
		var (
			amount      int32           = 1000
			currency    string          = "USD"
			expiryMonth string          = "10"
			expiryYear  string          = "2022"
			firstName   string          = "John"
			lastName    string          = "Doe"
			card        creditcard.Card = creditcard.Card{
				Number: cardNumber,
				Cvv:    "123",
				Month:  expiryMonth,
				Year:   expiryYear,
			}
		)

		if value > 0 {
			tests = append(tests, testCase{
				name:       fmt.Sprintf("transaction_success_%s", card),
				path:       "/transactions",
				method:     http.MethodPost,
				statusCode: MapErrors[value].StatusCode,
				headers: map[string]string{
					"X-Auth-Token": xAuthTokenFake,
				},
				input: fmt.Sprintf(`{
					"amount": %d,
					"currency": "%s",
					"expiry_month": "%s",
					"expiry_year": "%s",
					"first_name": "%s",
					"last_name": "%s",
					"number": "%s"
				}`, amount, currency, expiryMonth, expiryYear, firstName, lastName, cardNumber),
				expect: gin.H{
					"code":    MapErrors[value].Code,
					"message": MapErrors[value].Message,
				},
			})
		} else {
			lastFour, _ := card.LastFour()
			tests = append(tests, testCase{
				name:       fmt.Sprintf("transaction_success_%s", card),
				path:       "/transactions",
				method:     http.MethodPost,
				statusCode: http.StatusCreated,
				headers: map[string]string{
					"X-Auth-Token": xAuthTokenFake,
				},
				input: fmt.Sprintf(`{
					"amount": %d,
					"currency": "%s",
					"expiry_month": "%s",
					"expiry_year": "%s",
					"first_name": "%s",
					"last_name": "%s",
					"number": "%s"
				}`, amount, currency, expiryMonth, expiryYear, firstName, lastName, cardNumber),
				expect: gin.H{
					"amount":   amount,
					"currency": currency,
					"code":     00,
					"message":  "ok",
					"status":   statusApproved,
					"source": gin.H{
						"bin":          cardNumber[0:5],
						"card_type":    "credit_card",
						"expiry_month": expiryMonth,
						"expiry_year":  expiryYear,
						"first_name":   firstName,
						"last_name":    lastName,
						"number":       fmt.Sprintf("XXXX-XXXX-XXXX-%s", lastFour),
					},
				},
			})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := testRequest(router, tt.method, tt.path, tt.input, tt.headers)
			statusCode := w.Result().StatusCode

			if statusCode != tt.statusCode {
				t.Fatalf("got: %v output: %v", statusCode, tt.statusCode)
			}

			if tt.expect != nil {
				var response map[string]interface{}
				err := json.Unmarshal([]byte(w.Body.Bytes()), &response)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}

}
