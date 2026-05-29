package handlers

import (
	"construction-ar-backend/internal/auth"
	"construction-ar-backend/internal/middleware"
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  interface{} `json:"user"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mock Authentication - In real app, check hash in DB
	// For this SaaS version, we assume any valid email works for demo
	token, err := auth.GenerateToken("user-123", "org-alpha", "Admin")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User: map[string]string{
			"id": "user-123",
			"email": req.Email,
			"role": "Admin",
		},
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(claims)
}
