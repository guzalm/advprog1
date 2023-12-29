package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type JsonRequest struct {
	Message string `json:"message"`
}

type JsonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

const (
	successStatus      = "success"
	errorStatus        = "error"
	expectedFieldName  = "message"
	invalidJSONMessage = "Invalid JSON message"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePostRequest(w, r)
	case http.MethodGet:
		handleGetRequest(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		handleJSONRequest(w, r)
	default:
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
	}
}

func handleJSONRequest(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if message, ok := requestData[expectedFieldName].(string); !ok {
		http.Error(w, fmt.Sprintf("%s: Missing or invalid expected field '%s'", invalidJSONMessage, expectedFieldName), http.StatusBadRequest)
	} else {
		log.Printf("Received message: %v", message)
		fmt.Println("Server console log: Received message:", message)

		jsonResponse := JsonResponse{
			Status:  successStatus,
			Message: "Data successfully received",
		}

		if jsonResponseBytes, err := json.Marshal(jsonResponse); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponseBytes)
		}
	}
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "go.html")
}

func main() {
	http.HandleFunc("/", handleRequest)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
