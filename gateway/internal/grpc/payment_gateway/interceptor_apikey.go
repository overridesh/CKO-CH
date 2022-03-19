package payment_gateway

import (
	"context"
	"errors"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/overridesh/checkout-challenger/internal/repository"
)

const (
	authorization string = "Authorization"
)

/*
	ValidateAPIKey validate if the x-api-key is valid and set merchant-id for the request

	TODO: This should not go to the database as we are doing now.
	This could lead to a DDOS attack. But for the example case there is no problem.
*/

func ValidateAPIKey(merchantRepository repository.MerchantRepository) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, cannotGetMetadata.Err()
		}

		apiKeyHeader := md.Get(authorization)
		if len(apiKeyHeader) <= 0 {
			return nil, unauthorized.Err()
		}

		id, err := uuid.FromString(apiKeyHeader[0])
		if err != nil || id == uuid.Nil {
			return nil, unauthorized.Err()
		}

		id, err = merchantRepository.GetIDByApiKey(id)
		if err != nil {
			if errors.Is(err, repository.ErrMerchantNotFound) {
				return nil, unauthorized.Err()
			}
			return nil, internalServerErr.Err()
		}

		md.Append("merchant_id", id.String())
		return handler(metadata.NewIncomingContext(ctx, md), req)
	}
}
