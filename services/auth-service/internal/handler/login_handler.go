package handler

import (
	"auth-service/internal/service"
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
    Username string `json:"username"`
    PassHash string `json:"passhash"`
}

func Login(loginService *service.LoginService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req LoginRequest

        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        token, err := loginService.Login(req.Username, req.PassHash)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"token":token})
    }
}
