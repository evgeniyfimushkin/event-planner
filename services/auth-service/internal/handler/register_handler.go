package handler

import (
	"auth-service/internal/service"
	"encoding/json"
	"net/http"
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

        if err := registerService.Register(req.Username, req.Email, req.PassHash); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
    }
}

