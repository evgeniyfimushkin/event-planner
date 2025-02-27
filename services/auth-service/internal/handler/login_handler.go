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

        refreshToken, err := loginService.Login(req.Username, req.PassHash)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name: "refresh_token",
            Value: refreshToken,
            Path: "api/v1/auth/refresh",
            HttpOnly: true,
            Secure: true,
            SameSite: http.SameSiteStrictMode,
            MaxAge: 604800, //7days
        })
    }
}
