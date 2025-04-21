package handlers

import (
	"net/http"
	"pvz_service/objects"
	"pvz_service/services"

	"github.com/asaskevich/govalidator"
)

type ProductApiHandler struct {
	service services.ProductService
}

func NewProductHandler(svc services.ProductService) *ProductApiHandler {
	return &ProductApiHandler{
		service: svc,
	}
}

func (h *ProductApiHandler) AddHandler(w http.ResponseWriter, r *http.Request) {
	productData := &objects.AddProductDto{}
	if err := PasreJSON(r.Body, productData); err != nil {
		http.Error(w, "Can't decode request body", http.StatusBadRequest)
		return
	}

	if _, err := govalidator.ValidateStruct(productData); err != nil {
		HandleValidationError(w, err)
		return
	}

	productDataUpd, err := h.service.Add(r.Context(), *productData)
	if err != nil {
		http.Error(w, "Product Add Failed", http.StatusInternalServerError)
		return
	}

	RenderJSON(w, productDataUpd)
}