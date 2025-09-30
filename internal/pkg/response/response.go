package response

import (
	"encoding/json"
	"net/http"
)

// Response representa uma resposta padronizada da API
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// JSON envia uma resposta JSON
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Erro ao encodar resposta", http.StatusInternalServerError)
		}
	}
}

// Success envia uma resposta de sucesso
func Success(w http.ResponseWriter, status int, message string, data any) {
	JSON(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error envia uma resposta de erro
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{
		Success: false,
		Error:   message,
	})
}

// ErrorWithDetails envia uma resposta de erro com detalhes
func ErrorWithDetails(w http.ResponseWriter, status int, message, code, details string) {
	JSON(w, status, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
		Details: details,
	})
}
