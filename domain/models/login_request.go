package login_request

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
