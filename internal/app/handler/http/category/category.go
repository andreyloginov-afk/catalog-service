package hcategory

import (
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
	svcCategory service.Category
}

func NewHandler(svcCategory service.Category) rhandler.Category {
	return &handler{svcCategory: svcCategory}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	// декодирование JSON
	var req entity.RequestCategoryCreate

	// валидация параметров
	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.HandleError(w, r, err)
		return
	}
	// вызов сервиса
	category, err := h.svcCategory.Create(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAlreadyExists):
			httph.HandleError(w, r, err)
		default:
			httph.HandleError(w, r, err)
		}
		return
	}

	// конвертация entity ResponsDTO
	resp := entity.ResponseCategory{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	httph.SendEncoded(w, r, http.StatusCreated, resp)
}

func (h *handler) GetByGUID(w http.ResponseWriter, r *http.Request) {
	// извлекает guid
	vars := mux.Vars(r)
	guid, err := uuid.FromString(vars["guid"])
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}
	// вызывает сервис
	_, err = h.svcCategory.GetByGUID(r.Context(), guid)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAlreadyExists):
			httph.HandleError(w, r, err)
		default:
			httph.HandleError(w, r, err)
		}
		return
	}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	//
	vars := mux.Vars(r)
	guid, err := uuid.FromString(vars["guid"])
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}

	var req entity.RequestCategoryUpdate

	//
	if err := binding.ScanAndValidateJSON(r, &req); err != nil {
		httph.HandleError(w, r, err)
		return
	}

	_, err = h.svcCategory.Update(r.Context(), guid, req)
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

	err = h.svcCategory.Delete(r.Context(), guid)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			httph.HandleError(w, r, err)
		case errors.Is(err, entity.ErrCategoryHasProducts):
			httph.HandleError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	category, err := h.svcCategory.List(r.Context())
	if err != nil {
		httph.HandleError(w, r, err)
		return
	}

	resp := make([]entity.ResponseCategory, 0, len(category))
	for _, c := range category {
		resp = append(resp, entity.ResponseCategory{
			GUID:      c.GUID,
			Name:      c.Name,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}

	httph.SendEncoded(w, r, http.StatusOK, resp)
}
