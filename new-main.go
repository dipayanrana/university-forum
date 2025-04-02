package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db        *sql.DB
	store     *sessions.CookieStore
	templates *template.Template
)

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

	// Parse templates
	templates = template.Must(template.ParseGlob("templates/*.html"))
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

	fmt.Println("Server starting on port 7000...")
	log.Fatal(http.ListenAndServe(":7000", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	posts, err := getPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isAuthenticated, _ := session.Values["authenticated"].(bool)

	data := map[string]interface{}{
		"IsAuthenticated": isAuthenticated,
		"Posts":           posts,
	}

	err = templates.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Process registration form
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// For testing purposes, just redirect
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := templates.ExecuteTemplate(w, "layout", map[string]interface{}{
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

		// For testing purposes, just authenticate any user
		session, _ := store.Get(r, "session-name")
		session.Values["authenticated"] = true
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := templates.ExecuteTemplate(w, "layout", map[string]interface{}{
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
	err := templates.ExecuteTemplate(w, "layout", map[string]interface{}{
		"IsAuthenticated": true,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewPostHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "layout", map[string]interface{}{
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
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getPosts() ([]Post, error) {
	// Return empty posts for now
	return []Post{}, nil
}

type Post struct {
	ID         int64
	Title      string
	Content    string
	AuthorName string
	CreatedAt  string
}
