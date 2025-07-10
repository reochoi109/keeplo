package dto

type RegisterMonitorRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    string `json:"port"`
}

type UpdateMonitorRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    string `json:"port"`
}
