package main

import (
	"Concurrent_downloading/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/fetch", handlers.GetImage)

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
