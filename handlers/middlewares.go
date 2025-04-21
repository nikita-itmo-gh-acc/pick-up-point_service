package handlers

import (
	"context"
	"net/http"
	"pvz_service/logger"
	"time"
)

type responseWrapper struct {
	http.ResponseWriter
	status      int
}

func (rw *responseWrapper) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}


func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenStr := tokenHeader[len("Bearer "):]
		token, err := VerifyToken(tokenStr)
		if err != nil {
			http.Error(w, "Token verification failed", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(*Claims)
        ctx := context.WithValue(r.Context(), "role", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            foundRole, ok := r.Context().Value("role").(string)
            if !ok || foundRole != role {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &responseWrapper{ ResponseWriter: w, status: 200 }
		next.ServeHTTP(wrapper, r)
		logger.Info.Printf("[%s] %s - %d, %s %s\n",
			r.Method, r.URL.Path, wrapper.status, r.RemoteAddr, time.Since(start))
	})
}
