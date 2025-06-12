package bff

type RegisterPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenPayload struct {
	RefreshToken string `json:"refresh_token"`
}

type VerifyTokenPayload struct {
	AccessToken string `json:"access_token"`
}

type VerifyTokenResponse struct {
	IsValid bool `json:"is_valid"`
	UserUUID string `json:"user_uuid"`
}