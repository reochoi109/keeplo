package dto

// Request --------------------------------------
type SignupRequest struct {
	Email         string `json:"email" binding:"required,email"`
	NickName      string `json:"nickname" binding:"required,min=2,max=20"`
	Password      string `json:"password" binding:"required,min=8"`
	CheckPassword string `json:"check_password" binding:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateNicknameRequest struct {
	NickName string `json:"nickname" binding:"required,min=2,max=20"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}
type DuplicateEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type CheckPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

// Response --------------------------------------

type LoginResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}
