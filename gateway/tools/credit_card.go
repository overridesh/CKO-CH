package tools

import (
	"errors"
	"strings"
)

var (
	ErrCardNotEnough       error = errors.New("card number is not long enough")
	ErrCardNumberInvalid   error = errors.New("card number is invalid")
	ErrFirstNameIsRequired error = errors.New("first_name is required")
	ErrLastNameIsRequired  error = errors.New("last_name is required")
	ErrExpiryMonthWrong    error = errors.New("expiry_month need two numbers")
	ErrExpiryYearWrong     error = errors.New("expiry_year need four numbers")
)

func CardValidate(
	firstName,
	lastName,
	number,
	cvv,
	expiryMonth,
	expiryYear string,
) error {
	if len(strings.TrimSpace(firstName)) == 0 {
		return ErrFirstNameIsRequired
	}
	if len(strings.TrimSpace(lastName)) == 0 {
		return ErrLastNameIsRequired
	}
	if len(strings.TrimSpace(expiryMonth)) != 2 {
		return ErrExpiryMonthWrong
	}
	if len(strings.TrimSpace(expiryYear)) != 4 {
		return ErrExpiryYearWrong
	}

	numberLength := len(strings.TrimSpace(number))
	if !(numberLength >= 16 && numberLength <= 19) {
		return ErrCardNumberInvalid
	}
	return nil
}

func LastFour(number string) (string, error) {
	if len(number) < 4 {
		return "", ErrCardNotEnough
	}
	return number[len(number)-4:], nil
}
