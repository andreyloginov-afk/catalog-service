package mapper

import (
	"errors"

	"github.com/andreyloginov-afk/catalog-service/internal/app/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorToGRPC(err error) error {
	var appErr *entity.AppError

	if errors.As(err, &appErr) {
		return status.Error(appErr.GRPCCode(), appErr.Message)
	}
	return status.Error(codes.Internal, "internal error")
}
