package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><body><h1>Hello World!</h1><p>This is a test page.</p></body></html>")
	})

	fmt.Println("Test server starting on :8888...")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
