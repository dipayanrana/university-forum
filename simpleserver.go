package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templates map[string]*template.Template

func loadTemplates() {
	templates = make(map[string]*template.Template)
	templateNames := []string{"index", "login", "register", "create-post", "view-post"}
	layoutFile := filepath.Join("templates", "layout.html")

	for _, name := range templateNames {
		contentFile := filepath.Join("templates", name+".html")
		tmpl, err := template.ParseFiles(layoutFile, contentFile)
		if err != nil {
			log.Fatalf("Failed to parse template %s: %v", name, err)
		}
		templates[name] = tmpl
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := templates["index"]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"PageTitle":       "University Forum",
		"IsAuthenticated": false,
	}

	err := tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := templates["login"]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"PageTitle":       "Login",
		"IsAuthenticated": false,
	}

	err := tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, ok := templates["register"]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"PageTitle":       "Register",
		"IsAuthenticated": false,
	}

	err := tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	loadTemplates()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Simple server starting on :9999...")
	log.Fatal(http.ListenAndServe(":9999", nil))
}
