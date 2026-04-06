package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"example.com/m/internal/dto"
	"example.com/m/internal/models"
	"example.com/m/internal/repository"
)

func TestDefaultHandler(t *testing.T) {
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(DefaultHandler)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	expected := `"There is nothing here."` + "\n"
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v, want %v", recorder.Body.String(), expected)
	}
}

func TestCreateSurvey(t *testing.T) {
	dbPath := "./test.db"
	db, err := repository.OpenDB_test()
	if err != nil {
		t.Fatalf("failed at db open: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
		if err := os.Remove(dbPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("failed to delete remove test db %s: %v", dbPath, err)
		}
	})
	err = repository.InitSchema(db)
	if err != nil {
		t.Fatalf("failed at db initialization: %v", err)
	}
	def_handler := &Handler{DB: db}
	recorder := httptest.NewRecorder()
	request := dto.RequestCreateSurvey{Name: "Survey Name", Description: "Survey Description", Questions_list: []dto.RequestCreateQuestion{
		{
			Description: "Normal Description",
			Type:        0,
			IsMandatory: true,
			Choices: []models.Answer_choice{
				{
					Description: "Answer Choice description",
				},
			},
		},
		{
			Description: "Second question",
			Type:        1,
			IsMandatory: false,
		},
	}}
	decoder, err := json.Marshal(request) // marshal the request into json and put into decoder
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/survey", bytes.NewReader(decoder))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	handler := http.HandlerFunc(def_handler.CreateSurvey)
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusCreated)
	}
	expected := map[string]any{
		"message": "survey successfully created",
	}
	response := map[string]any{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed while unmarshalling a response: %v", err)
	}
	if response["message"] != expected["message"] {
		t.Errorf("response returned an unexpected result, got %v, want %v", response["message"], expected["message"])
	}
	if response["survey"] == nil {
		t.Errorf("no survey is present in response")
	}

}
