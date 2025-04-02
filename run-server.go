package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	runStore *sessions.CookieStore
	runDb    *sql.DB
)

// Post represents a forum post
type Post struct {
	ID         int64
	Title      string
	Content    string
	AuthorName string
	CreatedAt  string
}

// User represents a forum user
type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    string
}

// Comment represents a post comment
type Comment struct {
	ID         int64
	Content    string
	AuthorName string
	CreatedAt  string
}

func init() {
	var err error
	// Initialize database
	runDb, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables
	createRunTables()

	runStore = sessions.NewCookieStore([]byte("unique-key-123"))
}

func createRunTables() {
	// Users table
	_, err := runDb.Exec(`
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
	_, err = runDb.Exec(`
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
	_, err = runDb.Exec(`
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

func handleLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Validate input
		if username == "" || password == "" {
			renderLoginPage(w, "Username and password are required")
			return
		}

		// Get the user from the database
		var user User
		var passwordHash string
		err := runDb.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", username).
			Scan(&user.ID, &passwordHash)

		if err != nil {
			renderLoginPage(w, "Invalid username or password")
			return
		}

		// Compare passwords
		err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
		if err != nil {
			renderLoginPage(w, "Invalid username or password")
			return
		}

		// Set user as authenticated
		session, _ := runStore.Get(r, "session-name")
		session.Values["authenticated"] = true
		session.Values["user_id"] = user.ID
		session.Values["username"] = username
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	renderLoginPage(w, "")
}

func renderLoginPage(w http.ResponseWriter, errorMsg string) {
	t, err := template.ParseFiles(
		"./templates/layout.html",
		"./templates/login.html",
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := map[string]interface{}{
		"IsAuthenticated": false,
		"Username":        "",
		"ErrorMessage":    errorMsg,
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Execution error: %v", err), http.StatusInternalServerError)
		log.Printf("Execution error: %v", err)
	}
}

func handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate input
		if username == "" || email == "" || password == "" {
			renderRegisterPage(w, "All fields are required")
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error processing registration", http.StatusInternalServerError)
			return
		}

		// Insert the user into the database
		_, err = runDb.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)",
			username, email, string(hashedPassword))
		if err != nil {
			renderRegisterPage(w, "Username or email already taken")
			return
		}

		// Redirect to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	renderRegisterPage(w, "")
}

func renderRegisterPage(w http.ResponseWriter, errorMsg string) {
	t, err := template.ParseFiles(
		"./templates/layout.html",
		"./templates/register.html",
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := map[string]interface{}{
		"IsAuthenticated": false,
		"Username":        "",
		"ErrorMessage":    errorMsg,
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Execution error: %v", err), http.StatusInternalServerError)
		log.Printf("Execution error: %v", err)
	}
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	session, _ := runStore.Get(r, "session-name")
	isAuthenticated, _ := session.Values["authenticated"].(bool)
	username, _ := session.Values["username"].(string)

	// Get all posts
	posts, err := getRunPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles(
		"./templates/layout.html",
		"./templates/index.html",
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := map[string]interface{}{
		"IsAuthenticated": isAuthenticated,
		"Username":        username,
		"Posts":           posts,
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Execution error: %v", err), http.StatusInternalServerError)
		log.Printf("Execution error: %v", err)
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := runStore.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Values["user_id"] = nil
	session.Values["username"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	session, _ := runStore.Get(r, "session-name")
	isAuthenticated, _ := session.Values["authenticated"].(bool)
	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		title := r.FormValue("title")
		content := r.FormValue("content")

		// Validate input
		if title == "" || content == "" {
			renderCreatePostPage(w, "Title and content are required", r)
			return
		}

		// Get user ID from session
		userID, _ := session.Values["user_id"].(int64)

		// Insert post into database
		_, err := runDb.Exec("INSERT INTO posts (title, content, author_id) VALUES (?, ?, ?)",
			title, content, userID)
		if err != nil {
			renderCreatePostPage(w, "Error creating post", r)
			return
		}

		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	renderCreatePostPage(w, "", r)
}

func renderCreatePostPage(w http.ResponseWriter, errorMsg string, r *http.Request) {
	session, _ := runStore.Get(r, "session-name")
	isAuthenticated, _ := session.Values["authenticated"].(bool)
	username, _ := session.Values["username"].(string)

	t, err := template.ParseFiles(
		"./templates/layout.html",
		"./templates/create-post.html",
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := map[string]interface{}{
		"IsAuthenticated": isAuthenticated,
		"Username":        username,
		"ErrorMessage":    errorMsg,
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Execution error: %v", err), http.StatusInternalServerError)
		log.Printf("Execution error: %v", err)
	}
}

func getRunPosts() ([]Post, error) {
	rows, err := runDb.Query(`
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

		// Format date
		createdTime, _ := time.Parse("2006-01-02 15:04:05", p.CreatedAt)
		p.CreatedAt = createdTime.Format("Jan 02, 2006 at 3:04 PM")

		posts = append(posts, p)
	}
	return posts, nil
}

func handleViewPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	session, _ := runStore.Get(r, "session-name")
	isAuthenticated, _ := session.Values["authenticated"].(bool)
	username, _ := session.Values["username"].(string)

	// Get post details
	var post Post
	var authorID int64
	err := runDb.QueryRow(`
		SELECT p.id, p.title, p.content, p.created_at, p.author_id, u.username 
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.id = ?
	`, postID).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &authorID, &post.AuthorName)

	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Format date
	createdTime, _ := time.Parse("2006-01-02 15:04:05", post.CreatedAt)
	post.CreatedAt = createdTime.Format("Jan 02, 2006 at 3:04 PM")

	// Get comments for this post
	comments, err := getPostComments(post.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles(
		"./templates/layout.html",
		"./templates/view-post.html",
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := map[string]interface{}{
		"IsAuthenticated": isAuthenticated,
		"Username":        username,
		"Post":            post,
		"Comments":        comments,
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Execution error: %v", err), http.StatusInternalServerError)
		log.Printf("Execution error: %v", err)
	}
}

func handleAddComment(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	session, _ := runStore.Get(r, "session-name")
	isAuthenticated, _ := session.Values["authenticated"].(bool)
	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	postID := vars["id"]

	if r.Method == "POST" {
		content := r.FormValue("content")

		// Validate input
		if content == "" {
			http.Redirect(w, r, fmt.Sprintf("/post/%s", postID), http.StatusSeeOther)
			return
		}

		// Get user ID from session
		userID, _ := session.Values["user_id"].(int64)

		// Insert comment into database
		_, err := runDb.Exec("INSERT INTO comments (content, post_id, author_id) VALUES (?, ?, ?)",
			content, postID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Redirect to post page
	http.Redirect(w, r, fmt.Sprintf("/post/%s", postID), http.StatusSeeOther)
}

func getPostComments(postID int64) ([]Comment, error) {
	rows, err := runDb.Query(`
		SELECT c.id, c.content, c.created_at, u.username
		FROM comments c
		JOIN users u ON c.author_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err := rows.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.AuthorName)
		if err != nil {
			return nil, err
		}

		// Format date
		createdTime, _ := time.Parse("2006-01-02 15:04:05", c.CreatedAt)
		c.CreatedAt = createdTime.Format("Jan 02, 2006 at 3:04 PM")

		comments = append(comments, c)
	}
	return comments, nil
}

func handleUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileUsername := vars["username"]

	session, _ := runStore.Get(r, "session-name")
	isAuthenticated, _ := session.Values["authenticated"].(bool)
	currentUser, _ := session.Values["username"].(string)

	// Get user details
	var user User
	err := runDb.QueryRow(`
		SELECT id, username, email, created_at 
		FROM users
		WHERE username = ?
	`, profileUsername).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Format date
	createdTime, _ := time.Parse("2006-01-02 15:04:05", user.CreatedAt)
	user.CreatedAt = createdTime.Format("Jan 02, 2006")

	// Get user's posts
	posts, err := getUserPosts(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles(
		"./templates/layout.html",
		"./templates/profile.html",
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	isOwner := isAuthenticated && profileUsername == currentUser

	data := map[string]interface{}{
		"IsAuthenticated": isAuthenticated,
		"Username":        currentUser,
		"User":            user,
		"Posts":           posts,
		"IsOwner":         isOwner,
		"PostCount":       len(posts),
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Execution error: %v", err), http.StatusInternalServerError)
		log.Printf("Execution error: %v", err)
	}
}

func getUserPosts(userID int64) ([]Post, error) {
	rows, err := runDb.Query(`
		SELECT p.id, p.title, p.content, p.created_at, u.username
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.author_id = ?
		ORDER BY p.created_at DESC
	`, userID)
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

		// Format date
		createdTime, _ := time.Parse("2006-01-02 15:04:05", p.CreatedAt)
		p.CreatedAt = createdTime.Format("Jan 02, 2006 at 3:04 PM")

		// Truncate content for preview
		if len(p.Content) > 150 {
			p.Content = p.Content[:150] + "..."
		}

		posts = append(posts, p)
	}
	return posts, nil
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	session, _ := runStore.Get(r, "session-name")
	isAuthenticated, _ := session.Values["authenticated"].(bool)
	username, _ := session.Values["username"].(string)

	var posts []Post
	var err error

	if query != "" {
		posts, err = searchPosts(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	t, err := template.ParseFiles(
		"./templates/layout.html",
		"./templates/search.html",
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	data := map[string]interface{}{
		"IsAuthenticated": isAuthenticated,
		"Username":        username,
		"Query":           query,
		"Posts":           posts,
		"ResultCount":     len(posts),
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Execution error: %v", err), http.StatusInternalServerError)
		log.Printf("Execution error: %v", err)
	}
}

func searchPosts(query string) ([]Post, error) {
	searchQuery := "%" + query + "%"

	rows, err := runDb.Query(`
		SELECT p.id, p.title, p.content, p.created_at, u.username
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.title LIKE ? OR p.content LIKE ?
		ORDER BY p.created_at DESC
		LIMIT 20
	`, searchQuery, searchQuery)
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

		// Format date
		createdTime, _ := time.Parse("2006-01-02 15:04:05", p.CreatedAt)
		p.CreatedAt = createdTime.Format("Jan 02, 2006 at 3:04 PM")

		// Truncate content for preview
		if len(p.Content) > 150 {
			p.Content = p.Content[:150] + "..."
		}

		posts = append(posts, p)
	}
	return posts, nil
}

func runServer() {
	r := mux.NewRouter()

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	r.HandleFunc("/", handleHomePage).Methods("GET")
	r.HandleFunc("/login", handleLoginPage).Methods("GET", "POST")
	r.HandleFunc("/register", handleRegisterPage).Methods("GET", "POST")
	r.HandleFunc("/logout", handleLogout).Methods("GET")
	r.HandleFunc("/create-post", handleCreatePost).Methods("GET", "POST")
	r.HandleFunc("/post/{id:[0-9]+}", handleViewPost).Methods("GET")
	r.HandleFunc("/post/{id:[0-9]+}/comment", handleAddComment).Methods("POST")
	r.HandleFunc("/user/{username}", handleUserProfile).Methods("GET")
	r.HandleFunc("/search", handleSearch).Methods("GET")

	port := ":3999"
	fmt.Printf("Server starting on %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func main() {
	runServer()
}
