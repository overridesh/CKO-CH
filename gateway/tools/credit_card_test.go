package tools

import (
	"errors"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestLastFour(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output int
		err    error
	}{
		{
			name:   "LastFour",
			input:  uuid.NewV4().String(),
			output: 4,
			err:    nil,
		},
		{
			name:   "ShortInput",
			input:  "ABC",
			output: 0,
			err:    ErrCardNotEnough,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lastFour, err := LastFour(tt.input)
			if !errors.Is(err, tt.err) {
				t.Errorf("expect a error %v, but got %v", tt.err, err)
			}

			if tt.output != len(lastFour) {
				t.Errorf("expect a last four equal %d, but got %v", tt.output, len(lastFour))
			}
		})
	}
}

func TestCreditCardValidate(t *testing.T) {
	type creditCardTest struct {
		firstName   string
		lastName    string
		number      string
		cvv         string
		expiryMonth string
		expiryYear  string
	}

	tests := []struct {
		name   string
		input  creditCardTest
		output error
	}{
		{
			name: ErrCardNumberInvalid.Error(),
			input: creditCardTest{
				firstName:   "John",
				lastName:    "Doe",
				number:      "4444",
				cvv:         "123",
				expiryMonth: "10",
				expiryYear:  "2022",
			},
			output: ErrCardNumberInvalid,
		},
		{
			name: ErrFirstNameIsRequired.Error(),
			input: creditCardTest{
				lastName:    "Doe",
				number:      "4444444444445555",
				cvv:         "123",
				expiryMonth: "10",
				expiryYear:  "2022",
			},
			output: ErrFirstNameIsRequired,
		},
		{
			name: ErrLastNameIsRequired.Error(),
			input: creditCardTest{
				firstName:   "John",
				number:      "4444444444445555",
				cvv:         "123",
				expiryMonth: "10",
				expiryYear:  "2022",
			},
			output: ErrLastNameIsRequired,
		},
		{
			name: ErrExpiryMonthWrong.Error(),
			input: creditCardTest{
				firstName:   "John",
				lastName:    "Doe",
				number:      "4444444444445555",
				cvv:         "123",
				expiryMonth: "ABC",
				expiryYear:  "2022",
			},
			output: ErrExpiryMonthWrong,
		},
		{
			name: ErrExpiryYearWrong.Error(),
			input: creditCardTest{
				firstName:   "John",
				lastName:    "Doe",
				number:      "4444444444445555",
				cvv:         "123",
				expiryMonth: "10",
				expiryYear:  "0",
			},
			output: ErrExpiryYearWrong,
		},
		{
			name: "Error Nil",
			input: creditCardTest{
				firstName:   "John",
				lastName:    "Doe",
				number:      "4444444444445555",
				cvv:         "123",
				expiryMonth: "10",
				expiryYear:  "2022",
			},
			output: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CardValidate(
				tt.input.firstName,
				tt.input.lastName,
				tt.input.number,
				tt.input.cvv,
				tt.input.expiryMonth,
				tt.input.expiryYear,
			)
			if !errors.Is(err, tt.output) {
				t.Errorf("expect a  %v error, but got %v", tt.output, err)
			}
		})
	}
}
