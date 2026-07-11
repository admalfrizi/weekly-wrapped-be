package response

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Success(message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  200,
		Message: message,
		Data:    data,
	}
}

func SuccessWithPagination(message string, data interface{}, meta PaginationMeta) APIResponse {
	return APIResponse{
		Status:  200,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

func Error(status int, message string, errors interface{}) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Errors:  errors,
	}
}