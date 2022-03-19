package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/overridesh/checkout-challenger/pkg/service/acquirer"
	pbPaymentGateway "github.com/overridesh/checkout-challenger/proto"
)

var (
	// IdempotencyPointerIsNil error for pointer is nil
	IdempotencyPointerIsNil error = errors.New("pointer is nil")
)

type RecoveryPoint string

const (
	FirstPoint     RecoveryPoint = "first_point"
	CreatedPoint   RecoveryPoint = "transaction_created"
	PurchasedPoint RecoveryPoint = "transaction_purchased"
	UpdatedPoint   RecoveryPoint = "transaction_updated"
)

type Idempotency struct {
	RecoveryPoint RecoveryPoint                           `json:"recovery_point"`
	Response      *pbPaymentGateway.CreatePaymentResponse `json:"response"`
	BankResponse  *acquirer.PaymentResponse               `json:"bank_response"`
	MerchantID    uuid.UUID                               `json:"merchant_id"`
	Amount        int32                                   `json:"amount"`
	Currency      string                                  `json:"currency"`
	FirstName     string                                  `json:"first_name"`
	LastName      string                                  `json:"last_name"`
	Number        string                                  `json:"number"`
	ExpiryMonth   string                                  `json:"expiry_month"`
	ExpiryYear    string                                  `json:"expiry_year"`
	Hash          string                                  `json:"hash"`
}

func (i Idempotency) CompareHash(hash string) bool {
	var (
		equal bool
	)

	if len(i.Hash) == 0 {
		return equal
	} else {
		equal = (hash == i.Hash)
	}

	return equal
}

func (i *Idempotency) SetMD5() error {
	if i == nil {
		return IdempotencyPointerIsNil
	}

	bytes, err := json.Marshal(Idempotency{
		MerchantID:  i.MerchantID,
		Amount:      i.Amount,
		Currency:    i.Currency,
		FirstName:   i.FirstName,
		LastName:    i.LastName,
		Number:      i.Number,
		ExpiryMonth: i.ExpiryMonth,
		ExpiryYear:  i.ExpiryYear,
	})
	if err != nil {
		return err
	}

	hashMD5 := md5.Sum(bytes)
	i.Hash = hex.EncodeToString(hashMD5[:])

	return nil
}

func NewRecoveryPoint(raw string) RecoveryPoint {
	var recover RecoveryPoint = FirstPoint

	switch strings.ToLower(raw) {
	case strings.ToLower(CreatedPoint.String()):
		recover = CreatedPoint
	case strings.ToLower(PurchasedPoint.String()):
		recover = PurchasedPoint
	case strings.ToLower(UpdatedPoint.String()):
		recover = UpdatedPoint
	}

	return recover
}

func (s RecoveryPoint) String() string {
	return string(s)
}
