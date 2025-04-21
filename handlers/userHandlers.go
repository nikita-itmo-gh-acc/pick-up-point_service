package handlers

import (
	"fmt"
	"net/http"
	"pvz_service/objects"
	"pvz_service/services"

	"github.com/gorilla/mux"
)

func SendToken(w http.ResponseWriter, info *objects.UserDto) {
	tokenStr, err := CreateToken(info)
    if err != nil {
		http.Error(w, "Token creation failed", http.StatusInternalServerError)
		fmt.Println(err)
		return
    }

	token := &objects.TokenDto{ Token: tokenStr }
	RenderJSON(w, token)
}

type UserApiHandler struct {
	service services.UserService
}

func NewUserHandler(svc services.UserService) *UserApiHandler {
    return &UserApiHandler{service: svc}
}

func (h *UserApiHandler) SetUpRoutes(r *mux.Router) {
    r.HandleFunc("/register", h.RegisterHandler).Methods("POST")
    r.HandleFunc("/login", h.LoginHandler).Methods("POST")
	r.HandleFunc("/dummyLogin", h.DummyHandler).Methods("POST")
}

func (h *UserApiHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	regData := &objects.UserDto{}
	if err := PasreJSON(r.Body, regData); err != nil {
		http.Error(w, "Can't decode request body", http.StatusBadRequest)
		return
	}

	id, err := h.service.Register(r.Context(), *regData)
	if err != nil {
		http.Error(w, "Can't register user", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	RenderJSON(w, &objects.UserDto{ 
		Id: id, 
		Email: regData.Email, 
		Role: regData.Role,
	})
}

func (h *UserApiHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginData := &objects.UserDto{}
	if err := PasreJSON(r.Body, loginData); err != nil {
		http.Error(w, "Can't decode request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Login(r.Context(), *loginData); err != nil {
		switch err.Error() {
		case "wrong password":
			http.Error(w, "Wrong password", http.StatusForbidden)
			return
		default:
			http.Error(w, "Login failed", http.StatusInternalServerError)
			return
		}
	}

	SendToken(w, loginData)
}

func (h *UserApiHandler) DummyHandler(w http.ResponseWriter, r *http.Request) {
	dummy := &objects.UserDto{}
	if err := PasreJSON(r.Body, dummy); err != nil {
		http.Error(w, "Can't decode request body", http.StatusBadRequest)
		return
	}

	SendToken(w, dummy)
}
