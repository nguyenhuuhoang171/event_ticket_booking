package dto

// LoginRequest represents the login request payload.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// SignupRequest represents the signup request payload.
type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LogoutRequest represents the logout request payload.
type LogoutRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}
