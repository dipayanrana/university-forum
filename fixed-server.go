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
	dbFixed    *sql.DB
	storeFixed *sessions.CookieStore
)

func init() {
	var err error
	// Initialize database
	dbFixed, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables
	createFixedTables()

	storeFixed = sessions.NewCookieStore([]byte("super-secret-key"))
}

func createFixedTables() {
	// Users table
	_, err := dbFixed.Exec(`
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
	_, err = dbFixed.Exec(`
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
	_, err = dbFixed.Exec(`
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

func homeFixedHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := storeFixed.Get(r, "session-name")
	posts, err := getFixedPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isAuthenticated, _ := session.Values["authenticated"].(bool)

	// Direct template parsing for home page
	t, err := template.ParseFiles(
		"./templates/layout.html",
		"./templates/index.html",
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"IsAuthenticated": isAuthenticated,
		"Posts":           posts,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func registerFixedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Just redirect for now
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

func loginFixedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Just authenticate for testing
		session, _ := storeFixed.Get(r, "session-name")
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

func logoutFixedHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := storeFixed.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func runFixedServer() {
	r := mux.NewRouter()

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	r.HandleFunc("/", homeFixedHandler).Methods("GET")
	r.HandleFunc("/register", registerFixedHandler).Methods("GET", "POST")
	r.HandleFunc("/login", loginFixedHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", logoutFixedHandler).Methods("GET")

	fmt.Println("Fixed server starting on port 7000...")
	log.Fatal(http.ListenAndServe(":7000", r))
}

func getFixedPosts() ([]PostFixed, error) {
	// Return empty posts for now
	return []PostFixed{}, nil
}

type PostFixed struct {
	ID         int64
	Title      string
	Content    string
	AuthorName string
	CreatedAt  string
}

func main() {
	runFixedServer()
}
