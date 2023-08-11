package response

type RegisterResponse struct {
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}
