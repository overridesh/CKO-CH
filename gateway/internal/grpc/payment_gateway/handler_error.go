package payment_gateway

import (
	"github.com/overridesh/checkout-challenger/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	unauthorized        *status.Status = status.New(codes.InvalidArgument, "unauthorized")
	cannotGetMetadata   *status.Status = status.New(codes.InvalidArgument, "cannot get headers")
	IdMustBeUUID        *status.Status = status.New(codes.InvalidArgument, "the id must be uuid")
	IdMusBeValidUUID    *status.Status = status.New(codes.InvalidArgument, "the id must be a valid uuid")
	internalServerErr   *status.Status = status.New(codes.Internal, "internal server error")
	transactionNotFound *status.Status = status.New(codes.NotFound, repository.ErrTransactionNotFound.Error())
)
