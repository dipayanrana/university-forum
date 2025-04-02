# University Discussion Forum

A modern web-based discussion forum built specifically for university communities. This platform enables students and faculty to engage in meaningful discussions, share knowledge, and build a collaborative learning environment.

## Features

- User Authentication (Register/Login)
- Create and View Discussions
- Comment System
- Search Functionality
- Responsive Design
- Modern UI with Bootstrap
- SQLite Database

## Tech Stack

- Backend: Go (Golang)
- Frontend: HTML, CSS, JavaScript
- CSS Framework: Bootstrap 5
- Database: SQLite
- Dependencies:
  - gorilla/mux: URL router and dispatcher
  - gorilla/sessions: Session management
  - mattn/go-sqlite3: SQLite driver
  - golang.org/x/crypto/bcrypt: Password hashing

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/university-forum.git
cd university-forum
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

4. Access the forum at `http://localhost:8080`

## Project Structure

```
university-forum/
├── main.go              # Main application file
├── handlers/            # Request handlers
│   ├── auth.go         # Authentication handlers
│   └── posts.go        # Post and comment handlers
├── static/             # Static files
│   ├── css/           
│   │   └── style.css   # Custom styles
│   └── js/
│       └── main.js     # Client-side JavaScript
├── templates/          # HTML templates
│   ├── layout.html     # Base template
│   ├── index.html      # Home page
│   ├── login.html      # Login page
│   ├── register.html   # Registration page
│   ├── create-post.html# Create post page
│   ├── view-post.html  # View post page
│   └── search.html     # Search page
└── forum.db           # SQLite database
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 