package ghcatalogv1

import (
	"context"

	"github.com/andreyloginov-afk/catalog-service/internal/app/entity"
	"github.com/andreyloginov-afk/catalog-service/internal/app/mapper"
	mcatv1 "github.com/andreyloginov-afk/catalog-service/internal/app/mapper/catalog/v1"
	"github.com/andreyloginov-afk/catalog-service/internal/app/service"
	catalogv1 "github.com/andreyloginov-afk/catalog-service/internal/pkg/grpc/gen/catalog/v1"
	"github.com/gofrs/uuid"
)

type handler struct {
	catalogv1.UnimplementedCatalogServiceServer
	srv service.Product
}

func NewHandler(srv service.Product) catalogv1.CatalogServiceServer {
	return &handler{srv: srv}
}

func (h *handler) GetProduct(ctx context.Context, req *catalogv1.GetProductRequest) (*catalogv1.GetProductResponse, error) {
	rawGuid := req.GetGuid()
	guid, err := uuid.FromString(rawGuid)
	if err != nil {
		return nil, mapper.ErrorToGRPC(entity.ErrIncorrectParameters)
	}
	products, err := h.srv.GetByGUIDs(ctx, []uuid.UUID{guid})
	if err != nil {
		return nil, mapper.ErrorToGRPC(err)
	}
	if len(products) == 0 {
		return nil, mapper.ErrorToGRPC(entity.ErrNotFound)
	}
	return &catalogv1.GetProductResponse{
		Product: mcatv1.ProductToProto(products[0]),
	}, nil
}

func (h *handler) GetProducts(ctx context.Context, req *catalogv1.GetProductsRequest) (*catalogv1.GetProductsResponse, error) {
	if len(req.GetGuids()) == 0 {
		return &catalogv1.GetProductsResponse{}, nil
	}
	guids := make([]uuid.UUID, 0, len(req.GetGuids()))
	for _, rawGuid := range req.GetGuids() {
		guid, err := uuid.FromString(rawGuid)
		if err != nil {
			return nil, mapper.ErrorToGRPC(entity.ErrIncorrectParameters)
		}
		guids = append(guids, guid)
	}
	product, err := h.srv.GetByGUIDs(ctx, guids)
	if err != nil {
		return nil, mapper.ErrorToGRPC(err)
	}
	foundSet := make(map[uuid.UUID]struct{}, len(product))
	protoProducts := make([]*catalogv1.Product, 0, len(product))
	for _, p := range product {
		foundSet[p.GUID] = struct{}{}
		protoProducts = append(protoProducts, mcatv1.ProductToProto(p))
	}

	missingGuids := make([]string, 0)
	for _, guid := range guids {
		if _, ok := foundSet[guid]; !ok {
			missingGuids = append(missingGuids, guid.String())
		}
	}

	return &catalogv1.GetProductsResponse{
		Products:     protoProducts,
		MissingGuids: missingGuids,
	}, nil
}
