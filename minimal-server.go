package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var (
	sessionStore *sessions.CookieStore
)

func init() {
	sessionStore = sessions.NewCookieStore([]byte("minimal-key"))
}

func serveLoginPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"./templates/layout.html",
		"./templates/login.html",
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	err = t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"IsAuthenticated": false,
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Execution error: %v", err), http.StatusInternalServerError)
		log.Printf("Execution error: %v", err)
	}
}

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body><h1>Home Page</h1><p>This is a minimal home page. <a href='/login'>Login</a></p></body></html>")
}

func runMinimalServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", serveHomePage).Methods("GET")
	r.HandleFunc("/login", serveLoginPage).Methods("GET")

	fmt.Println("Minimal server starting on port 5000...")
	log.Fatal(http.ListenAndServe(":5000", r))
}

func main() {
	runMinimalServer()
}
