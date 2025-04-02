package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you accessed: %s", r.URL.Path)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><body><h1>Login Page</h1><form><input type='text' name='username'><input type='password' name='password'><button type='submit'>Login</button></form></body></html>")
	})

	fmt.Println("Server starting on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
