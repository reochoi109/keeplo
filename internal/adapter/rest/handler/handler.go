package handler

import (
	"keeplo/internal/application/monitor"
	"keeplo/internal/application/user"
)

type Handler struct {
	UserService    user.Service
	MonitorService monitor.Service
}

func NewHandler(userService user.Service, monitorService monitor.Service) *Handler {
	return &Handler{
		UserService:    userService,
		MonitorService: monitorService,
	}
}
