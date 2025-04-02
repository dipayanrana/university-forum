package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Post struct {
	ID         int64
	Title      string
	Content    string
	AuthorID   int64
	AuthorName string
	CreatedAt  string
	Comments   []Comment
}

type Comment struct {
	ID         int64
	Content    string
	AuthorID   int64
	AuthorName string
	CreatedAt  string
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		title := r.FormValue("title")
		content := r.FormValue("content")
		userID := session.Values["user_id"].(int64)

		if title == "" || content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("INSERT INTO posts (title, content, author_id) VALUES (?, ?, ?)",
			title, content, userID)
		if err != nil {
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "create-post.html", nil)
}

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var post Post
	err = db.QueryRow(`
		SELECT p.id, p.title, p.content, p.author_id, u.username, p.created_at
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.id = ?
	`, postID).Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.AuthorName, &post.CreatedAt)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Get comments
	rows, err := db.Query(`
		SELECT c.id, c.content, c.author_id, u.username, c.created_at
		FROM comments c
		JOIN users u ON c.author_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC
	`, postID)
	if err != nil {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.AuthorID, &comment.AuthorName, &comment.CreatedAt)
		if err != nil {
			continue
		}
		post.Comments = append(post.Comments, comment)
	}

	templates.ExecuteTemplate(w, "view-post.html", post)
}

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	if content == "" {
		http.Error(w, "Comment content is required", http.StatusBadRequest)
		return
	}

	userID := session.Values["user_id"].(int64)
	_, err = db.Exec("INSERT INTO comments (content, post_id, author_id) VALUES (?, ?, ?)",
		content, postID, userID)
	if err != nil {
		http.Error(w, "Error adding comment", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+vars["id"], http.StatusSeeOther)
}
