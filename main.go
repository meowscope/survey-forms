package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"example.com/m/internal/handlers"
	"example.com/m/internal/models"
	"example.com/m/internal/repository"
)

var tempDB = []models.Survey{}

func main() {
	mux := http.NewServeMux()
	db, err := repository.OpenDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = repository.InitSchema(db)
	if err != nil {
		log.Fatalf("failed at db initialization: %v", err)
	}

	def_handler := &handlers.Handler{Mu: &sync.RWMutex{}, DB: db, TempDB: &tempDB}

	mux.HandleFunc("/", def_handler.DefaultHandler)
	mux.HandleFunc("POST /surveys", def_handler.CreateSurvey)
	mux.HandleFunc("DELETE /surveys", def_handler.DeleteSurvey)
	mux.HandleFunc("GET /surveys", def_handler.GetSurveys)

	fmt.Printf("Server should be running on 8080 port now.\n")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
