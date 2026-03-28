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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.svcProduct.Create(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAlreadyExists):
			httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, entity.ErrNotFound):
			httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		default:
			httph.ErrorApply(w, http.StatusInternalServerError, err.Error())
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
		httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.svcProduct.GetByGUID(r.Context(), guid)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			httph.ErrorApply(w, http.StatusNotFound, err.Error())
		default:
			httph.ErrorApply(w, http.StatusInternalServerError, err.Error())
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

	httph.SendEncoded(w, r, http.StatusOK, resp)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	guid, err := uuid.FromString(vars["guid"])
	if err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		return
	}

	var req entity.RequestProductUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.svcProduct.Update(r.Context(), guid, req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			httph.ErrorApply(w, http.StatusNotFound, err.Error())
		case errors.Is(err, entity.ErrAlreadyExists):
			httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		default:
			httph.ErrorApply(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	guid, err := uuid.FromString(vars["guid"])
	if err != nil {
		httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.svcProduct.Delete(r.Context(), guid)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			httph.ErrorApply(w, http.StatusNotFound, err.Error())
		case errors.Is(err, entity.ErrCategoryHasProducts):
			httph.ErrorApply(w, http.StatusBadRequest, err.Error())
		default:
			httph.ErrorApply(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.svcProduct.List(r.Context())
	if err != nil {
		httph.ErrorApply(w, http.StatusInternalServerError, err.Error())
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
