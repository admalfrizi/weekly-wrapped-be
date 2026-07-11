package response

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Success(message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  200,
		Message: message,
		Data:    data,
	}
}

func Error(status int, message string, errors interface{}) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Errors:  errors,
	}
}