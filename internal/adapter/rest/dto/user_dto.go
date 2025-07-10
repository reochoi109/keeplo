package dto

// Request --------------------------------------
type SignupRequest struct {
	Email         string `json:"email"`
	NickName      string `json:"nickname"`
	Password      string `json:"password"`
	CheckPassword string `json:"check_password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserInfoRequest struct {
	NickName      string `json:"nickname"`
	Password      string `json:"password"`
	CheckPassword string `json:"check_password"`
}

type DuplicateEmailRequest struct {
	Email string `json:"email"`
}

type CheckPasswordRequest struct {
	Password string `json:"password"`
}

// Response --------------------------------------
