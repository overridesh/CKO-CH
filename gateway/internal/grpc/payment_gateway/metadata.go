package payment_gateway

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
)

var (
	merchantId     string = "merchant_id"
	idempotencyKey string = "idempotency_key"
)

// GetMerchantID get merchant id from metadata
func GetMerchantID(ctx context.Context) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, cannotGetMetadata.Err()
	}

	merchantIdSTR, ok := md[merchantId]
	if !ok || len(merchantIdSTR) <= 0 {
		return uuid.Nil, unauthorized.Err()
	}

	merchantID, err := uuid.FromString(merchantIdSTR[0])
	if err != nil || merchantID == uuid.Nil {
		return uuid.Nil, unauthorized.Err()
	}

	if merchantID == uuid.Nil {
		return uuid.Nil, unauthorized.Err()
	}

	return merchantID, nil
}

// GetIdempotencyKey if exists
func GetIdempotencyKey(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	key, ok := md[idempotencyKey]
	if !ok || len(key) <= 0 {
		return "", false
	}

	return key[0], true
}
