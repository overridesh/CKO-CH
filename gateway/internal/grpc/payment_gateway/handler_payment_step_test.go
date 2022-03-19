package payment_gateway

import (
	"context"
	"database/sql"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"

	"github.com/overridesh/checkout-challenger/internal/model"
	"github.com/overridesh/checkout-challenger/internal/repository"
	mockPkg "github.com/overridesh/checkout-challenger/pkg/mock"
	pbPaymentGateway "github.com/overridesh/checkout-challenger/proto"
)

func TestSuccessTransactionStep(t *testing.T) {
	var (
		firstUUID uuid.UUID = uuid.NewV4()
	)
	tests := []struct {
		name   string
		input  func() string
		expect uuid.UUID
	}{
		{
			name: "TestSuccessStep_Success_WithEmpty_Transaction",
			input: func() string {
				svc := NewGRPC(nil, nil, nil)
				return svc.SuccessTransactionStep(&model.Transaction{
					Id: firstUUID,
				}).Id
			},
			expect: firstUUID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input()
			if result != tt.expect.String() {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", result, tt.expect)
			}
		})
	}
}

func TestPaymentStep(t *testing.T) {
	var (
		firstUUID uuid.UUID = uuid.NewV4()
	)

	tests := []struct {
		name   string
		input  func() (*model.Transaction, error)
		expect error
	}{
		{
			name: "TestPaymentStepStep_WithoutError",
			input: func() (*model.Transaction, error) {
				tx := model.Transaction{
					Id: firstUUID,
				}
				paymentGatewayRepository := new(mockPkg.PaymentGatewayRepository)
				paymentGatewayRepository.On("Create", tx).Return(&tx, nil)
				svc := NewGRPC(paymentGatewayRepository, nil, nil)
				return svc.CreateTransaction(context.Background(), tx)
			},
			expect: nil,
		},
		{
			name: "TestPaymentStepStep_WithError",
			input: func() (*model.Transaction, error) {
				paymentGatewayRepository := new(mockPkg.PaymentGatewayRepository)
				paymentGatewayRepository.On("Create", model.Transaction{}).Return(nil, sql.ErrNoRows)
				svc := NewGRPC(paymentGatewayRepository, nil, nil)
				return svc.CreateTransaction(context.Background(), model.Transaction{})
			},
			expect: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.input()
			if err != tt.expect {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", err, tt.expect)
			}
		})
	}
}

func TestUpdateTransactionStep(t *testing.T) {
	tests := []struct {
		name   string
		input  func() error
		expect error
	}{
		{
			name: "TestUpdateTransactionStep_WithoutError",
			input: func() error {
				paymentGatewayRepository := new(mockPkg.PaymentGatewayRepository)
				paymentGatewayRepository.On("Update", &model.Transaction{}, mock.AnythingOfType("func(*model.Transaction) error")).Return(nil)
				svc := NewGRPC(paymentGatewayRepository, nil, nil)
				return svc.UpdateTransactionStep(context.Background(), &model.Transaction{})
			},
			expect: nil,
		},
		{
			name: "TestUpdateTransactionStep_WithError",
			input: func() error {
				paymentGatewayRepository := new(mockPkg.PaymentGatewayRepository)
				paymentGatewayRepository.On("Update", &model.Transaction{}, mock.AnythingOfType("func(*model.Transaction) error")).Return(repository.ErrTransactionNotFound)
				svc := NewGRPC(paymentGatewayRepository, nil, nil)
				return svc.UpdateTransactionStep(context.Background(), &model.Transaction{})
			},
			expect: repository.ErrTransactionNotFound,
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

func TestFromPuchasedTransactionPoint(t *testing.T) {
	tests := []struct {
		name   string
		input  func() (*pbPaymentGateway.CreatePaymentResponse, error)
		expect error
	}{
		{
			name: "TestFromPuchasedTransactionPoint_WithoutError",
			input: func() (*pbPaymentGateway.CreatePaymentResponse, error) {
				paymentGatewayRepository := new(mockPkg.PaymentGatewayRepository)
				paymentGatewayRepository.On("Update", &model.Transaction{}, mock.AnythingOfType("func(*model.Transaction) error")).Return(nil)
				svc := NewGRPC(paymentGatewayRepository, nil, nil)
				return svc.FromPuchasedTransactionPoint(context.Background(), &model.Transaction{}, nil)
			},
			expect: nil,
		},
		{
			name: "TestFromPuchasedTransactionPoint_WithError",
			input: func() (*pbPaymentGateway.CreatePaymentResponse, error) {
				paymentGatewayRepository := new(mockPkg.PaymentGatewayRepository)
				paymentGatewayRepository.On("Update", &model.Transaction{}, mock.AnythingOfType("func(*model.Transaction) error")).Return(repository.ErrTransactionNotFound)
				svc := NewGRPC(paymentGatewayRepository, nil, nil)
				return svc.FromPuchasedTransactionPoint(context.Background(), &model.Transaction{}, nil)
			},
			expect: repository.ErrTransactionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.input()
			if err != tt.expect {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", err, tt.expect)
			}
		})
	}
}

func TestFromCreatedTrasactionPoint(t *testing.T) {
	tests := []struct {
		name   string
		input  func() (*pbPaymentGateway.CreatePaymentResponse, error)
		expect error
	}{
		{
			name: "TestFromCreatedTrasactionPoint_WithoutError",
			input: func() (*pbPaymentGateway.CreatePaymentResponse, error) {
				paymentGatewayRepository := new(mockPkg.PaymentGatewayRepository)
				paymentGatewayRepository.On("Update", &model.Transaction{}, mock.AnythingOfType("func(*model.Transaction) error")).Return(nil)
				svc := NewGRPC(paymentGatewayRepository, nil, nil)

				return svc.FromCreatedTrasactionPoint(context.Background(), &model.Transaction{})
			},
			expect: nil,
		},
		{
			name: "TestFromCreatedTrasactionPoint_WithError",
			input: func() (*pbPaymentGateway.CreatePaymentResponse, error) {
				paymentGatewayRepository := new(mockPkg.PaymentGatewayRepository)
				paymentGatewayRepository.On("Update", &model.Transaction{}, mock.AnythingOfType("func(*model.Transaction) error")).Return(repository.ErrTransactionNotFound)
				svc := NewGRPC(paymentGatewayRepository, nil, nil)

				return svc.FromCreatedTrasactionPoint(context.Background(), &model.Transaction{})
			},
			expect: repository.ErrTransactionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.input()
			if err != tt.expect {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", err, tt.expect)
			}
		})
	}
}

func TestFromFromFirstPoint(t *testing.T) {
	tests := []struct {
		name   string
		input  func() (*pbPaymentGateway.CreatePaymentResponse, error)
		expect error
	}{
		{
			name: "TestFromFromFirstPoint_WithoutError",
			input: func() (*pbPaymentGateway.CreatePaymentResponse, error) {
				paymentGatewayRepository := new(mockPkg.PaymentGatewayRepository)
				tx := model.Transaction{
					Id: uuid.NewV4(),
				}

				paymentGatewayRepository.On("Create", tx).Return(&tx, nil)
				paymentGatewayRepository.On("Update", &tx, mock.AnythingOfType("func(*model.Transaction) error")).Return(nil)

				svc := NewGRPC(paymentGatewayRepository, nil, nil)

				return svc.FromFirstPoint(context.Background(), tx)
			},
			expect: nil,
		},
		{
			name: "TestFromFromFirstPoint_WithError_Create",
			input: func() (*pbPaymentGateway.CreatePaymentResponse, error) {
				paymentGatewayRepository := new(mockPkg.PaymentGatewayRepository)
				tx := model.Transaction{
					Id: uuid.NewV4(),
				}

				paymentGatewayRepository.On("Create", tx).Return(nil, sql.ErrNoRows)

				svc := NewGRPC(paymentGatewayRepository, nil, nil)

				return svc.FromFirstPoint(context.Background(), tx)
			},
			expect: sql.ErrNoRows,
		},
		{
			name: "TestFromFromFirstPoint_WithError_Update",
			input: func() (*pbPaymentGateway.CreatePaymentResponse, error) {
				paymentGatewayRepository := new(mockPkg.PaymentGatewayRepository)
				tx := model.Transaction{
					Id: uuid.NewV4(),
				}

				paymentGatewayRepository.On("Create", tx).Return(&tx, nil)
				paymentGatewayRepository.On("Update", &tx, mock.AnythingOfType("func(*model.Transaction) error")).Return(sql.ErrConnDone)

				svc := NewGRPC(paymentGatewayRepository, nil, nil)

				return svc.FromFirstPoint(context.Background(), tx)
			},
			expect: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.input()
			if err != tt.expect {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", err, tt.expect)
			}
		})
	}
}
