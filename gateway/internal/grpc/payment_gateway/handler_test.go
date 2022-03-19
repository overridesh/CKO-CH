package payment_gateway

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/overridesh/checkout-challenger/internal/model"
	"github.com/overridesh/checkout-challenger/internal/repository"
	mockRepository "github.com/overridesh/checkout-challenger/pkg/mock"
	"github.com/overridesh/checkout-challenger/pkg/service/acquirer"
	"github.com/overridesh/checkout-challenger/pkg/storage/cache"
	pbPaymentGateway "github.com/overridesh/checkout-challenger/proto"
	uuid "github.com/satori/go.uuid"
)

func dialer(
	repository repository.PaymentGatewayRepository,
	bankSimulatorService acquirer.Acquirer,
	cache cache.Cache,
) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pbPaymentGateway.RegisterPaymentGatewayServiceServer(
		server,
		NewGRPC(repository, bankSimulatorService, cache),
	)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestGetPayment(t *testing.T) {
	tests := []struct {
		name   string
		input  func() (*pbPaymentGateway.GetPaymentDetailsResponse, error)
		output *status.Status
	}{
		{
			name: "GetPayment_Success",
			input: func() (*pbPaymentGateway.GetPaymentDetailsResponse, error) {
				paymentGatewayRepository := new(mockRepository.PaymentGatewayRepository)
				tx := model.Transaction{
					Id:         uuid.NewV4(),
					MerchantID: uuid.NewV4(),
				}

				paymentGatewayRepository.On("GetByID", tx.Id, tx.MerchantID).Return(&tx, nil)

				ctx := context.Background()
				conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(paymentGatewayRepository, nil, nil)))
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()

				md := metadata.NewOutgoingContext(ctx, metadata.MD{
					"merchant_id": []string{tx.MerchantID.String()},
				})

				client := pbPaymentGateway.NewPaymentGatewayServiceClient(conn)

				return client.GetPayment(md, &pbPaymentGateway.GetPaymentDetailsRequest{
					Id: tx.Id.String(),
				})
			},
			output: nil,
		},
		{
			name: "GetPayment_Unauthorized",
			input: func() (*pbPaymentGateway.GetPaymentDetailsResponse, error) {
				paymentGatewayRepository := new(mockRepository.PaymentGatewayRepository)
				tx := model.Transaction{
					Id:         uuid.NewV4(),
					MerchantID: uuid.NewV4(),
				}

				ctx := context.Background()
				conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(paymentGatewayRepository, nil, nil)))
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()

				client := pbPaymentGateway.NewPaymentGatewayServiceClient(conn)

				return client.GetPayment(ctx, &pbPaymentGateway.GetPaymentDetailsRequest{
					Id: tx.Id.String(),
				})
			},
			output: unauthorized,
		},
		{
			name: "GetPayment_IdMustBeUUID",
			input: func() (*pbPaymentGateway.GetPaymentDetailsResponse, error) {
				paymentGatewayRepository := new(mockRepository.PaymentGatewayRepository)

				ctx := context.Background()
				conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(paymentGatewayRepository, nil, nil)))
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()

				md := metadata.NewOutgoingContext(ctx, metadata.MD{
					"merchant_id": []string{uuid.NewV4().String()},
				})

				client := pbPaymentGateway.NewPaymentGatewayServiceClient(conn)

				return client.GetPayment(md, &pbPaymentGateway.GetPaymentDetailsRequest{
					Id: "ABC",
				})
			},
			output: IdMustBeUUID,
		},
		{
			name: "GetPayment_IdMusBeValidUUID",
			input: func() (*pbPaymentGateway.GetPaymentDetailsResponse, error) {
				paymentGatewayRepository := new(mockRepository.PaymentGatewayRepository)

				ctx := context.Background()
				conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(paymentGatewayRepository, nil, nil)))
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()

				md := metadata.NewOutgoingContext(ctx, metadata.MD{
					"merchant_id": []string{uuid.NewV4().String()},
				})

				client := pbPaymentGateway.NewPaymentGatewayServiceClient(conn)

				return client.GetPayment(md, &pbPaymentGateway.GetPaymentDetailsRequest{
					Id: uuid.Nil.String(),
				})
			},
			output: IdMusBeValidUUID,
		},
		{
			name: "GetPayment_NotFound",
			input: func() (*pbPaymentGateway.GetPaymentDetailsResponse, error) {
				paymentGatewayRepository := new(mockRepository.PaymentGatewayRepository)
				tx := model.Transaction{
					Id:         uuid.NewV4(),
					MerchantID: uuid.NewV4(),
				}

				paymentGatewayRepository.On("GetByID", tx.Id, tx.MerchantID).Return(nil, repository.ErrTransactionNotFound)

				ctx := context.Background()
				conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(paymentGatewayRepository, nil, nil)))
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()

				md := metadata.NewOutgoingContext(ctx, metadata.MD{
					"merchant_id": []string{tx.MerchantID.String()},
				})

				client := pbPaymentGateway.NewPaymentGatewayServiceClient(conn)

				return client.GetPayment(md, &pbPaymentGateway.GetPaymentDetailsRequest{
					Id: tx.Id.String(),
				})
			},
			output: transactionNotFound,
		},
		{
			name: "GetPayment_Internal",
			input: func() (*pbPaymentGateway.GetPaymentDetailsResponse, error) {
				paymentGatewayRepository := new(mockRepository.PaymentGatewayRepository)
				tx := model.Transaction{
					Id:         uuid.NewV4(),
					MerchantID: uuid.NewV4(),
				}

				paymentGatewayRepository.On("GetByID", tx.Id, tx.MerchantID).Return(nil, repository.ErrMerchantNotFound)

				ctx := context.Background()
				conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(paymentGatewayRepository, nil, nil)))
				if err != nil {
					log.Fatal(err)
				}
				defer conn.Close()

				md := metadata.NewOutgoingContext(ctx, metadata.MD{
					"merchant_id": []string{tx.MerchantID.String()},
				})

				client := pbPaymentGateway.NewPaymentGatewayServiceClient(conn)

				return client.GetPayment(md, &pbPaymentGateway.GetPaymentDetailsRequest{
					Id: tx.Id.String(),
				})
			},
			output: internalServerErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.input()
			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.output.Code() {
						t.Errorf("error code: expected %v, received %v", codes.InvalidArgument, er.Code())
					}
					if er.Message() != tt.output.Message() {
						t.Errorf("error message: expected %v, received %v", tt.output.Message(), er.Message())
					}
				}
			}
		})
	}
}
