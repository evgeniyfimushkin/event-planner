package handler

import (
	"auth-service/internal/service"
	"encoding/json"
	"net/http"
)


func Refresh(refreshService *service.RefreshService) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("refresh_token")
        if err != nil {
            http.Error(w, "refresh token is missing", http.StatusUnauthorized)
            return
        }

        accessToken, err := refreshService.Refresh(cookie.Value)
        if err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name: "access_token",
            Value: accessToken,
            Path: "/",
            HttpOnly: true,
            Secure: true,
            SameSite: http.SameSiteStrictMode,
        })

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
    }
}
