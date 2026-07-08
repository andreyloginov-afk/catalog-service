package mcatv1

import (
	"github.com/andreyloginov-afk/catalog-service/internal/app/entity"
	catalogv1 "github.com/andreyloginov-afk/catalog-service/internal/pkg/grpc/gen/catalog/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProductToProto(p entity.Product) *catalogv1.Product {
	var desc string
	if p.Description != nil {
		desc = *p.Description
	}
	return &catalogv1.Product{
		Guid:         p.GUID.String(),
		CreatedAt:    timestamppb.New(p.CreatedAt),
		Price:        int64(p.Price),
		Name:         p.Name,
		Description:  desc,
		CategoryGuid: p.CategoryGuid.String(),
		UpdatedAt:    timestamppb.New(p.UpdatedAt),
	}
}
