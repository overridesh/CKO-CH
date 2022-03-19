package payment_gateway

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/overridesh/checkout-challenger/pkg/storage/cache"
)

const (
	idempotencyKeyMetadata string = "X-Idempotency-Key"
)

/*
	*** IMPORTANT ***
	This is just a small example of Idempotency Key.
	For now it has no recovery points. Only the save in case of success.
*/

func IdempotencyKeyMiddleware(cache cache.Cache) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		idempotencyKey := md.Get(idempotencyKeyMetadata)

		merchantID, err := GetMerchantID(ctx)
		if err != nil {
			return nil, err
		}

		if len(idempotencyKey) > 0 {
			key := fmt.Sprintf("%s_%s", merchantID, idempotencyKey[0])
			md.Append("idempotency_key", key)
		}

		return handler(metadata.NewIncomingContext(ctx, md), req)
	}
}
