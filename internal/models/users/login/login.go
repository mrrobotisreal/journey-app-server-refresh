package models_users_login

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID   int64  `json:"userID"`
	Firebase string `json:"firebaseID"`
	Username string `json:"username"`
	APIKey   string `json:"APIKey"`
}
