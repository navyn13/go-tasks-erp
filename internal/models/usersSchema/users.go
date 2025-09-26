package usersSchema

// LoginRequest and LoginResponse structs
type LoginResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignupRequest and SignupResponse structs
type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type SignupResponse struct {
	Message string `json:"message"`
}
