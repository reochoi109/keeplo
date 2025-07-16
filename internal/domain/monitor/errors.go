package monitor

import "errors"

var (
	ErrMonitorNotFound      = errors.New("monitor not found")
	ErrPermissionDenied     = errors.New("permission denied")
	ErrInvalidMonitorData   = errors.New("invalid monitor data")
	ErrMonitorAlreadyExists = errors.New("monitor already exists")
	ErrMonitorInactive      = errors.New("monitor is inactive")
)
