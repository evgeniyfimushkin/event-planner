package handler

import (
	"auth-service/internal/service"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

type RegisterRequest struct {
    Username string `json:"username"`
    PassHash string `json:"passhash"`
    Email string `json:"email,omitempty"`
}

func Register(registerService *service.RegisterService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req RegisterRequest

        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        user, err := registerService.Register(req.Username, req.Email, req.PassHash)
        if err != nil {
            var pgErr *pgconn.PgError
            if errors.As(err, &pgErr) && pgErr.Code == "23505" {
                http.Error(w, "User already exists", http.StatusConflict) // 409 Conflict
                return
            }
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        user.PassHash = ""
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)
    }
}

