package model

import (
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestCompareHash(t *testing.T) {
	tests := []struct {
		name   string
		input  func() bool
		expect bool
	}{
		{
			name: "CompareHash_Success",
			input: func() bool {
				var (
					merchantId uuid.UUID = uuid.NewV4()
				)
				firstItem := &Idempotency{
					MerchantID:  merchantId,
					Amount:      100,
					Currency:    "USD",
					FirstName:   "John",
					LastName:    "Doe",
					Number:      "4485040371536584",
					ExpiryMonth: "10",
					ExpiryYear:  "2222",
				}
				firstItem.SetMD5()

				secondItem := &Idempotency{
					MerchantID:  merchantId,
					Amount:      100,
					Currency:    "USD",
					FirstName:   "John",
					LastName:    "Doe",
					Number:      "4485040371536584",
					ExpiryMonth: "10",
					ExpiryYear:  "2222",
				}
				secondItem.SetMD5()
				return firstItem.CompareHash(secondItem.Hash)
			},
			expect: true,
		},
		{
			name: "CompareHash_Failed",
			input: func() bool {
				firstItem := &Idempotency{
					MerchantID:  uuid.NewV4(),
					Amount:      100,
					Currency:    "USD",
					FirstName:   "John",
					LastName:    "Doe",
					Number:      "4485040371536584",
					ExpiryMonth: "10",
					ExpiryYear:  "2222",
				}
				firstItem.SetMD5()

				secondItem := &Idempotency{
					MerchantID:  uuid.NewV4(),
					Amount:      100,
					Currency:    "USD",
					FirstName:   "John",
					LastName:    "Doe",
					Number:      "4485040371536584",
					ExpiryMonth: "10",
					ExpiryYear:  "2222",
				}
				secondItem.SetMD5()
				return firstItem.CompareHash(secondItem.Hash)
			},
			expect: false,
		},
		{
			name: "CompareHash_Success",
			input: func() bool {
				firstItem := &Idempotency{
					Hash: "",
				}
				return firstItem.CompareHash("")
			},
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input()
			if result != tt.expect {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", result, tt.expect)
			}
		})
	}
}

func TestSetMD5(t *testing.T) {
	tests := []struct {
		name   string
		input  *Idempotency
		expect error
	}{
		{
			name:   "SetMD5_Success",
			input:  nil,
			expect: IdempotencyPointerIsNil,
		},
		{
			name: "SetMD5_Success",
			input: &Idempotency{
				MerchantID:  uuid.NewV4(),
				Amount:      100,
				Currency:    "USD",
				FirstName:   "John",
				LastName:    "Doe",
				Number:      "4485040371536584",
				ExpiryMonth: "10",
				ExpiryYear:  "2222",
			},
			expect: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.SetMD5()
			if err != tt.expect {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", err, tt.expect)
			}
		})
	}
}

func TestNewRecoveryPoint(t *testing.T) {
	tests := []struct {
		name   string
		input  RecoveryPoint
		expect RecoveryPoint
	}{
		{
			name:   "RecoveryPoint",
			input:  NewRecoveryPoint(""),
			expect: FirstPoint,
		},
		{
			name:   "RecoveryCreated",
			input:  NewRecoveryPoint("TRANSACTION_CREATED"),
			expect: CreatedPoint,
		},
		{
			name:   "RecoveryPurchased",
			input:  NewRecoveryPoint("transaction_purchased"),
			expect: PurchasedPoint,
		},
		{
			name:   "RecoveryUpdated",
			input:  NewRecoveryPoint("TransacTion_Updated"),
			expect: UpdatedPoint,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input != tt.expect {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", tt.input, tt.expect)
			}
		})
	}
}
