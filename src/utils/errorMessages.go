package utils

import (
	"encoding/json"
	"net/http"
)

func WrongMethod(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)

	response := map[string]interface{}{
		"status":  405,
		"success": false,
		"message": "Method not Allowed",
	}

	json.NewEncoder(w).Encode(response)
}

func UnAuthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	response := map[string]interface{}{
		"status":  401,
		"success": false,
		"message": message,
	}

	json.NewEncoder(w).Encode(response)
}

func InvalidInput(w http.ResponseWriter, message ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	responseMessage := "Invalid Request Payload"
	if len(message) > 0 {
		responseMessage = message[0]
	}

	response := map[string]interface{}{
		"status":  400,
		"success": false,
		"message": responseMessage,
	}

	json.NewEncoder(w).Encode(response)
}

func InternalServerError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	response := map[string]interface{}{
		"status":  500,
		"success": false,
		"message": message,
	}

	json.NewEncoder(w).Encode(response)
}
