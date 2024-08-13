package main

import (
	"log"
	"net/http"
	"voting-system/routes"
)

func main() {
	router := routes.SetupRouter()
	log.Println("Server is running on port 5000...")
	log.Fatal(http.ListenAndServe(":5000", router))
}
