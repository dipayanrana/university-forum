package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"university-forum/handlers"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db          *sql.DB
	store       *sessions.CookieStore
	templateMap map[string]*template.Template
)

func loadTemplates() {
	templateMap = make(map[string]*template.Template)
	templateNames := []string{"index", "login", "register", "create-post", "view-post"}
	layoutFile := filepath.Join("templates", "layout.html")

	for _, name := range templateNames {
		contentFile := filepath.Join("templates", name+".html")
		tmpl, err := template.ParseFiles(layoutFile, contentFile)
		if err != nil {
			log.Fatalf("Failed to parse template %s: %v", name, err)
		}
		templateMap[name] = tmpl
	}
}

func init() {
	var err error
	// Initialize database
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables
	createTables()

	store = sessions.NewCookieStore([]byte("super-secret-key"))

	// Load templates
	loadTemplates()

	// Initialize handlers
	handlers.InitHandlers(db, store, nil)
}

func createTables() {
	// Users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Posts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			author_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (author_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Comments table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			post_id INTEGER NOT NULL,
			author_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (author_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/register", registerHandler).Methods("GET", "POST")
	r.HandleFunc("/login", loginHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")
	r.HandleFunc("/create-post", createPostHandler).Methods("GET", "POST")
	r.HandleFunc("/post/{id:[0-9]+}", viewPostHandler).Methods("GET")
	r.HandleFunc("/post/{id:[0-9]+}/comment", addCommentHandler).Methods("POST")

	port := ":7000"
	fmt.Printf("Server starting on %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	posts, err := getPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isAuthenticated, _ := session.Values["authenticated"].(bool)

	tmpl, ok := templateMap["index"]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		"IsAuthenticated": isAuthenticated,
		"Posts":           posts,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Process registration form
		_ = r.FormValue("username")
		_ = r.FormValue("email")
		_ = r.FormValue("password")

		// For testing purposes, just redirect
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Direct template parsing for register page
	t, err := template.ParseFiles(
		"./templates/layout.html",
		"./templates/register.html",
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"IsAuthenticated": false,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Process login form
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Basic validation
		if username == "" || password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		// For testing purposes, just authenticate any user
		session, _ := store.Get(r, "session-name")
		session.Values["authenticated"] = true
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Direct template parsing for login page
	t, err := template.ParseFiles(
		"./templates/layout.html",
		"./templates/login.html",
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"IsAuthenticated": false,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	// Code for create post...
	tmpl, ok := templateMap["create-post"]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err := tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		"IsAuthenticated": true,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewPostHandler(w http.ResponseWriter, r *http.Request) {
	// Code for view post...
	tmpl, ok := templateMap["view-post"]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err := tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		"IsAuthenticated": true,
		"Title":           "Sample Post",
		"Content":         "This is a sample post.",
		"AuthorName":      "Admin",
		"Comments":        []string{},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Code for adding comments...
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getPosts() ([]Post, error) {
	rows, err := db.Query(`
		SELECT p.id, p.title, p.content, p.created_at, u.username
		FROM posts p
		JOIN users u ON p.author_id = u.id
		ORDER BY p.created_at DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.AuthorName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

type Post struct {
	ID         int64
	Title      string
	Content    string
	AuthorName string
	CreatedAt  string
}
