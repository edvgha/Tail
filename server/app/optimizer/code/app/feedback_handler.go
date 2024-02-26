package app

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

func feedbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var request FeedBackRequest
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	response := update(&request)

	// Marshal the response data into JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to indicate JSON response
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}

func update(request *FeedBackRequest) FeedBackResponse {
	state, OK := Cache.Pop(request.ID)
	if !OK {
		return FeedBackResponse{Ack: false}
	}

	space := Spaces[state.ContextHash]
	log.Debug().Msgf("update: ctx: %s imp: %t", state.ContextHash, request.Impression)
	space.Update(state, request.Impression)
	return FeedBackResponse{Ack: true}
}
