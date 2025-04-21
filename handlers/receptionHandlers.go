package handlers

import (
	"fmt"
	"net/http"
	"pvz_service/objects"
	"pvz_service/services"

	"github.com/asaskevich/govalidator"
)

type ReceptionApiHandler struct {
	service services.ReceptionService
}

func NewReceptionHandler(svc services.ReceptionService) *ReceptionApiHandler {
	return &ReceptionApiHandler{
		service: svc,
	}
}

func (h *ReceptionApiHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	receptionData := &objects.ReceptionDto{}

	fmt.Println("Im here!")
	if err := PasreJSON(r.Body, receptionData); err != nil {
		http.Error(w, "Can't decode request body", http.StatusBadRequest)
		return
	}

	fmt.Println(receptionData)
	if _, err := govalidator.ValidateStruct(receptionData); err != nil {
		HandleValidationError(w, err)
		return
	}

	receptionDataUpd, err := h.service.StartReception(r.Context(), *receptionData)
	if err != nil {
		http.Error(w, "Reception Creation Failed", http.StatusInternalServerError)
		return
	}

	RenderJSON(w, receptionDataUpd)
}
