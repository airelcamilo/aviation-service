package dto

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Data     interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(data interface{}, message string) Response {
	return Response{
		Success: true,
		Data:    data,
		Message: message,
	}
}

func NewErrorResponse(err interface{}, message string) Response {
	return Response{
		Success: false,
		Error:   err,
		Message: message,
	}
}
