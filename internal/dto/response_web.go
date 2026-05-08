package dto

type ResponseWeb[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type OldResponseWeb[T any] struct {
	Message      string `json:"message"`
	Success      bool   `json:"success"`
	ResponseData T      `json:"response_data"`
}
