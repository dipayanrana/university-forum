package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func serveLogin(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"./templates/layout.html",
		"./templates/login.html",
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"IsAuthenticated": false,
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Execution error: %v", err), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Home page - <a href='/login'>Login</a>")
	})
	http.HandleFunc("/login", serveLogin)

	fmt.Println("Server starting on port 6000...")
	log.Fatal(http.ListenAndServe(":6000", nil))
}
