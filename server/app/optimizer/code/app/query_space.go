package app

import (
	"encoding/json"
	"net/http"
	"tail.server/app/optimizer/code/space"
)

func spaceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	queryParams := r.URL.Query()
	ctx := queryParams.Get("ctx")

	response := prepareSpace(ctx)

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

func prepareSpace(ctx string) space.LearnedEstimation {
	s, exists := Spaces[ctx]
	if !exists {
		return space.LearnedEstimation{}
	}

	return s.WC()
}
