package handlers

import (
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"os"

	"pvz_service/objects"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func HandleValidationError(w http.ResponseWriter, err error) {
	if errList, ok := err.(govalidator.Errors); ok {
		for _, field := range errList {
			fmt.Printf("Validation error in field %v\n", field)
		} 
	}
	http.Error(w, "query params didn't pass validation", http.StatusBadRequest)
}

func RenderJSON(w http.ResponseWriter, object interface{}) {
	js, err := json.Marshal(object)
	if err != nil {
		http.Error(w, "Can't render JSON from object:", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func PasreJSON(r io.Reader, object interface{}) error {
	body, _ := io.ReadAll(r)
	if err := json.Unmarshal(body, object); err != nil {
		return err
	}
	return nil
}

func GetSecret() []byte {
	return []byte(os.Getenv("SECRET_KEY"))
}

func VerifyToken(tokenStr string) (*jwt.Token, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return GetSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func CreateToken(info *objects.UserDto) (string, error) {
	claims := Claims{
		Role: info.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   info.Id.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	unsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return unsigned.SignedString(GetSecret())
}
