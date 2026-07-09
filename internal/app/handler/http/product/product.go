package hproduct

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andreyloginov-afk/catalog-service/internal/app/entity"
	rhandler "github.com/andreyloginov-afk/catalog-service/internal/app/handler/http"
	"github.com/andreyloginov-afk/catalog-service/internal/app/service"
	"github.com/andreyloginov-afk/catalog-service/internal/pkg/http/binding"
	"github.com/andreyloginov-afk/catalog-service/internal/pkg/http/httph"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

type handler struct {
	svcProduct service.Product
}

func NewHandler(svcProduct service.Product) rhandler.Product {
	return &handler{svcProduct: svcProduct}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req entity.RequestProductCreate

	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.HandleError(w, r, err)
		return
	}

	product, err := h.svcProduct.Create(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAlreadyExists):
			httph.HandleError(w, r, err)
		case errors.Is(err, entity.ErrNotFound):
			httph.HandleError(w, r, err)
		default:
			httph.HandleError(w, r, err)
		}
		return
	}

	resp := entity.ResponseProduct{
		GUID:         product.GUID,
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		CategoryGUID: product.CategoryGuid,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}

	httph.SendEncoded(w, r, http.StatusCreated, resp)
}

func (h *handler) GetByGUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	guid, err := uuid.FromString(vars["guid"])
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}

	products, err := h.svcProduct.GetByGUIDs(r.Context(), []uuid.UUID{guid})
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}
	if len(products) == 0 {
		httph.HandleError(w, r, entity.ErrNotFound)
		return
	}
	product := products[0]

	resp := entity.ResponseProduct{
		GUID:         product.GUID,
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		CategoryGUID: product.CategoryGuid,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}

	httph.SendEncoded(w, r, http.StatusOK, resp)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	guid, err := uuid.FromString(vars["guid"])
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}

	var req entity.RequestProductUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httph.HandleError(w, r, err)
		return
	}

	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.HandleError(w, r, err)
		return
	}

	_, err = h.svcProduct.Update(r.Context(), guid, req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			httph.HandleError(w, r, err)
		case errors.Is(err, entity.ErrAlreadyExists):
			httph.HandleError(w, r, err)
		default:
			httph.HandleError(w, r, err)
		}
		return
	}
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	guid, err := uuid.FromString(vars["guid"])
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}

	err = h.svcProduct.Delete(r.Context(), guid)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			httph.HandleError(w, r, err)
		case errors.Is(err, entity.ErrCategoryHasProducts):
			httph.HandleError(w, r, err)
		default:
			httph.HandleError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.svcProduct.List(r.Context(), entity.RequestProductList{})
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}

	resp := make([]entity.ResponseProduct, 0, len(products))
	for _, p := range products {
		resp = append(resp, entity.ResponseProduct{
			GUID:         p.GUID,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			CategoryGUID: p.CategoryGuid,
			CreatedAt:    p.CreatedAt,
			UpdatedAt:    p.UpdatedAt,
		})
	}

	httph.SendEncoded(w, r, http.StatusOK, resp)
}
