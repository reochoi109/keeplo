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
	Token  string `json:"token" example:"eyJhbGciOiJIUzI1NiIsIn..."`
	UserID string `json:"user_id" example:"user-uuid-string"`
	Email  string `json:"email" example:"user@example.com"`
}

type UserResponse struct {
	ID    string `json:"id" example:"user-uuid-string"`
	Email string `json:"email" example:"user@example.com"`
}

type DuplicateEmailResponse struct {
	IsDuplicate bool `json:"is_duplicate" example:"true"`
}

func NewLoginResponse(token, userID, email string) LoginResponse {
	return LoginResponse{
		Token:  token,
		UserID: userID,
		Email:  email,
	}
}

func NewUserResponse(id, email string) UserResponse {
	return UserResponse{
		ID:    id,
		Email: email,
	}
}

func NewDuplicateEmailResponse(duplicate bool) DuplicateEmailResponse {
	return DuplicateEmailResponse{
		IsDuplicate: duplicate,
	}
}
