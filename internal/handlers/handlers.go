package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/m/internal/models"
)

type Handler struct {
	TempDB *[]models.Survey
}

func (h *Handler) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

// Creates a new survey using a struct Survey
func (h *Handler) CreateSurvey(w http.ResponseWriter, r *http.Request) {
	new_survey := models.Survey{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&new_survey)
	if err != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	}
	err = models.ValidateSurvey(new_survey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	*h.TempDB = append(*h.TempDB, new_survey)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]any{
		"message": "survey successfully created",
		"survey":  new_survey,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
