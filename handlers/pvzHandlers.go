package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"pvz_service/objects"
	"pvz_service/services"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type PvzApiHandler struct {
	service services.PvzService
	receptionService services.ReceptionService
	productService services.ProductService
}

func NewPvzHandler(svc services.PvzService, recpSvc services.ReceptionService, prodSvc services.ProductService) *PvzApiHandler {
	return &PvzApiHandler{
		service: svc,
		receptionService: recpSvc,
		productService: prodSvc,
	}
}

func (h *PvzApiHandler) SetUpRoutes(r *mux.Router) {
	uuidRegexp := `[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}`
	Employee := RequireRole("employee")
	closeHandler := Employee(AuthMiddleware(http.HandlerFunc(h.CloseReceptionHandler)))
	deleteHandler := Employee(AuthMiddleware(http.HandlerFunc(h.DeleteProductHandler)))
	r.HandleFunc("", h.PvzOps).Methods("GET", "POST")
	r.Handle(fmt.Sprintf("/{pvzId:%s}/close_last_reception", uuidRegexp), closeHandler).Methods("POST")
	r.Handle(fmt.Sprintf("/{pvzId:%s}/delete_last_product", uuidRegexp), deleteHandler).Methods("POST")
}

func (h *PvzApiHandler) PvzOps(w http.ResponseWriter, r *http.Request) {
	Moderator := RequireRole("moderator")
	switch r.Method {
	case "GET":
		AuthMiddleware(http.HandlerFunc(h.FilterPvzHandler)).ServeHTTP(w, r)
	case "POST":
		Moderator(AuthMiddleware(http.HandlerFunc(h.CreatePvzHandler))).ServeHTTP(w, r)
	}
}

func (h *PvzApiHandler) FilterPvzHandler(w http.ResponseWriter, r *http.Request) {
	pvzQuery := &objects.PvzQuery{ StartDate: "1970-01-01", EndDate: "2100-01-01", Page: 1, Limit: 10 }
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if err := decoder.Decode(pvzQuery, r.URL.Query()); err != nil {
		fmt.Println("bad request query")
		http.Error(w, "Can't parse query params", http.StatusInternalServerError)
		return 
	}

	result, err := h.service.FilterPvz(r.Context(), *pvzQuery)
	switch err {
	case nil:
		RenderJSON(w, result)
	case sql.ErrNoRows:
		http.Error(w, "No pvz found", http.StatusNotFound)
		return
	default:
		http.Error(w, "Error: " + err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PvzApiHandler) CreatePvzHandler(w http.ResponseWriter, r *http.Request) {
	pvzData := &objects.PvzDto{}
	if err := PasreJSON(r.Body, pvzData); err != nil {
		http.Error(w, "Can't decode request body", http.StatusBadRequest)
		return
	}
	if _, err := govalidator.ValidateStruct(pvzData); err != nil {
		HandleValidationError(w, err)
		return
	}
	id, err := h.service.CreatePvz(r.Context(), *pvzData)
	if err != nil {
		http.Error(w, "Error during pvz creation", http.StatusInternalServerError)
		return
	}

	RenderJSON(w, &objects.PvzDto{
		Id: id,
		RegistrationDate: pvzData.RegistrationDate,
		City: pvzData.City,
	})
}

func checkParamProvided(w http.ResponseWriter, r *http.Request, param string) (string, bool) {
	params := mux.Vars(r)
	p, ok := params[param]
	if !ok {
		fmt.Println("no param " + param)
		http.Error(w, param + " url variable is required", http.StatusBadRequest)
		return "", false
	}
	return p, true
}

func (h *PvzApiHandler) CloseReceptionHandler(w http.ResponseWriter, r *http.Request) {
	pvzIdStr, ok := checkParamProvided(w, r, "pvzId")
	if !ok {
		return
	}
	
	pvzId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		http.Error(w, "Invalid pvzId param", http.StatusBadRequest)
		return
	}
	
	receptionDto, err := h.receptionService.CloseReception(r.Context(), pvzId)
	if err != nil {
		http.Error(w, "Reception close failed", http.StatusInternalServerError)
		return
	}

	RenderJSON(w, receptionDto)
}

func (h *PvzApiHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	pvzIdStr, ok := checkParamProvided(w, r, "pvzId")
	if !ok {
		return
	}

	pvzId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		http.Error(w, "Invalid pvzId param", http.StatusBadRequest)
		return
	}

	if err := h.productService.DeleteLast(r.Context(), pvzId); err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return 
	}

	w.WriteHeader(200)
	w.Write([]byte("Last product deleted"))
}
