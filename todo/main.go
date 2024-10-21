package main

import (
	"fmt"
	"log"
	"net/http"
	"todo/handlers" // Use correct path to the handlers package
)

func main() {
	http.HandleFunc("/add/todo", handlers.AddHandler)
	http.HandleFunc("/get/all", handlers.GetHandler)
	http.HandleFunc("/get/todo/", handlers.GetByIDHandler)
	http.HandleFunc("/edit/todo/", handlers.EditByID)
	http.HandleFunc("/delete/todo/", handlers.DeletebyID)

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
