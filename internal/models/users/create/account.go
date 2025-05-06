package models_users_create

type CreateAccountRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

type CreateAccountResponse struct {
	UserID int64  `json:"userID"`
	APIKey string `json:"APIKey"`
}
