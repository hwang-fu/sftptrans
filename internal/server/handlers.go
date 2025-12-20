package server

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:`
}

type APIError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
