package dto

type ResponseWeb[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}
