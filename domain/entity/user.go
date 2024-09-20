package entity

type User struct {
	UserID       string   `json:"user_id"`
	AppID        string   `json:"app_id"`
	Email        string   `json:"email"`
	Picture      string   `json:"picture"`
	Role         []string `json:"role"`
	GivenName    string   `json:"given_name"`
	FamilyName   string   `json:"family_name"`
	Name         string   `json:"name"`
	RefreshToken string   `json:"refresh_token"`
}

type LoginRequest struct {
	IdToken string `json:"id_token"`
	Email   string `json:"email"`
	AppID   string `json:"app_id"`
	Os      string `json:"os"`
}

type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ID           string   `json:"id"`
	Role         []string `json:"role"`
}

type RefreshTokenResp struct {
	AccessToken string   `json:"access_token"`
	ID          string   `json:"id"`
	Role        []string `json:"role"`
}

type RequestAccessToken struct {
	RefreshToken string `json:"refresh_token"`
}
