package handler

import "keeplo/internal/application/user"

type Handler struct {
	UserService user.Service
}

func NewHandler(userService user.Service) *Handler {
	return &Handler{
		UserService: userService,
	}
}
