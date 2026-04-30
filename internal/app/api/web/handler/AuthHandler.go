package handler

import (
	"encoding/json"
	"net/http"

	"service/internal/pkg/middleware"
)

// AuthHandler handles OAuth2 authentication endpoints.
type AuthHandler struct{}

// Redirect starts the OAuth2 authorization flow by redirecting the browser
// to the GX authorization server.
func (h AuthHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	middleware.Redirect(w, r)
}

// Callback handles the authorization code callback from the GX server,
// exchanges the code for a token, and redirects to the frontend.
func (h AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	middleware.Callback(w, r)
}

// Me returns the authenticated employee's profile.
func (h AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	emp := middleware.GetEmployee(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}
