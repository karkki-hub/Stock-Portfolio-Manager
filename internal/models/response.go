package models

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(message string, data interface{}) APIResponse { // Helper function to create a success response
	return APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string) APIResponse { // Helper function to create an error response
	return APIResponse{
		Status:  "error",
		Message: message,
	}
}
