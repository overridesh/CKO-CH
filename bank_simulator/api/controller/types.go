package controller

import "net/http"

// Headers simulate a x-auth-token for all request
type Headers struct {
	XAuthToken string `header:"X-Auth-Token" binding:"required"`
}

// Request struct for create a new payment
type Request struct {
	Amount      int32  `json:"amount" binding:"required"`
	Currency    string `json:"currency"  binding:"required"`
	ExpiryMonth string `json:"expiry_month"  binding:"required"`
	ExpiryYear  string `json:"expiry_year"  binding:"required"`
	FirstName   string `json:"first_name"  binding:"required"`
	LastName    string `json:"last_name"  binding:"required"`
	Number      string `json:"number"  binding:"required"`
	Cvv         string `json:"cvv"`
}

// Source card information
type Source struct {
	Bin         string `json:"bin"`
	CardType    string `json:"card_type"`
	ExpiryMonth string `json:"expiry_month"`
	ExpiryYear  string `json:"expiry_year"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Number      string `json:"number"`
}

// Response struct for success responses
type Response struct {
	Amount   int32  `json:"amount"`
	Currency string `json:"currency"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Status   string `json:"status"`
	Source   Source `json:"source"`
}

type Error struct {
	Code       string
	Message    string
	StatusCode int
}

// MapErrors Possible errors
var MapErrors map[int]Error = map[int]Error{
	1: {
		Code:       "100000",
		Message:    "Not authenticated",
		StatusCode: http.StatusUnauthorized,
	},
	2: {
		Code:       "100001",
		Message:    "Authentication could not be performed",
		StatusCode: http.StatusUnauthorized,
	},
	3: {
		Code:       "100002",
		Message:    "Attempted authentication",
		StatusCode: http.StatusUnprocessableEntity,
	},
	4: {
		Code:       "100003",
		Message:    "Authentication rejected",
		StatusCode: http.StatusUnprocessableEntity,
	},
	5: {
		Code:       "100004",
		Message:    "Card not enrolled",
		StatusCode: http.StatusUnprocessableEntity,
	},
	6: {
		Code:       "100005",
		Message:    "Error message during scheme communication",
		StatusCode: http.StatusUnprocessableEntity,
	},
	7: {
		Code:       "100006",
		Message:    "Card number is invalid",
		StatusCode: http.StatusBadRequest,
	},
	8: {
		Code:       "100007",
		Message:    "Expiration is invalid",
		StatusCode: http.StatusBadRequest,
	},
	9: {
		Code:       "100008",
		Message:    "Payload is invalid",
		StatusCode: http.StatusBadRequest,
	},
}

// TestCards to simulate actions
var TestCards map[string]int = map[string]int{
	"4539628347117863": 1,
	"5309961755464047": 1,
	"4024007186645015": 2,
	"5234106378657904": 2,
	"4556574722325580": 3,
	"5558468902774508": 3,
	"4275765574319271": 4,
	"5596061690670931": 4,
	"4484070000035519": 5,
	"5352151570003404": 5,
	"4452927588210665": 6,
	"5291144083573579": 6,
	"4485040371536584": 0,
	"4543474002249996": 0,
	"5588686116426417": 0,
	"5436031030606378": 0,
	"5199992312641465": 0,
	"345678901234564":  0,
}
