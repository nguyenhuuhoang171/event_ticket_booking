package dto

// LoginResponse represents the login response payload.
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

// SignupResponse represents the signup response payload.
type SignupResponse struct {
	AccessToken string `json:"access_token"`
}

// LogoutResponse represents the logout response payload.
type LogoutResponse struct {
}
