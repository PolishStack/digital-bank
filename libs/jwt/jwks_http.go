package jwt

import (
    "encoding/json"
    "net/http"
)

// HTTP handler to return JWKS JSON
func JWKSHandler(s Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        set, err := s.JWKS()
        if err != nil {
            http.Error(w, "internal", http.StatusInternalServerError)
            return
        }
        b, err := json.Marshal(set)
        if err != nil {
            http.Error(w, "internal", http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.Write(b)
    }
}
