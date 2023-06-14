package main

import (
	"log"
	"net/http"

	"assignment_3/handlers"
	"assignment_3/services"
)

func main() {
	
	go services.UpdateData()
	router := handlers.SetupRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}