package acquirer

// PaymentRequest expected purchase request
type PaymentRequest struct {
	Amount      int32  `json:"-"`
	Currency    string `json:"-"`
	ExpiryMonth string `json:"-"`
	ExpiryYear  string `json:"-"`
	FirstName   string `json:"-"`
	LastName    string `json:"-"`
	Number      string `json:"-"`
}

// PaymentResponse expected purchase response
type PaymentResponse struct {
	Amount      int32  `json:"-"`
	Currency    string `json:"-"`
	ExpiryMonth string `json:"-"`
	ExpiryYear  string `json:"-"`
	FirstName   string `json:"-"`
	LastName    string `json:"-"`
	Number      string `json:"-"`
	Code        string `json:"-"`
	Summary     string `json:"-"`
	CardBin     string `json:"-"`
	CardType    string `json:"-"`
	Status      string `json:"-"`
	StatusCode  int    `json:"-"`
}

// Acquirer service factory
type Acquirer interface {
	Purchase(*PaymentRequest) (*PaymentResponse, error)
}
