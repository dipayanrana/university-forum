# University Discussion Forum

A modern web-based discussion forum built specifically for university communities. This platform enables students and faculty to engage in meaningful discussions, share knowledge, and build a collaborative learning environment.

## Features

- **User Authentication**
  - Secure registration and login system
  - Password encryption using bcrypt
  - Session management with cookies
  - User profile management

- **Discussion Management**
  - Create rich text discussions with titles and content
  - View all discussions in a paginated list
  - Individual discussion view with comments
  - Author information and timestamps

- **Comment System**
  - Add comments to discussions
  - Nested comment display
  - Comment author information and timestamps

- **Search Functionality**
  - Search across all discussions by keywords
  - Results display with highlighted matches
  - Filter search results by relevance
  - Real-time search suggestions

- **Responsive Design**
  - Mobile-first approach
  - Adapts to different screen sizes
  - Consistent experience across devices
  - Touch-friendly interface

- **Modern UI**
  - Clean and intuitive interface with Bootstrap 5
  - Consistent design language
  - Accessibility features
  - Interactive elements with JavaScript enhancements

- **Database**
  - SQLite for lightweight deployment
  - Structured data models for users, posts, and comments
  - Efficient querying and indexing
  - Data integrity with foreign key constraints

## Tech Stack

### Backend
- **Go (Golang)**: A fast, statically typed language with excellent concurrency support
  - Clean syntax and strong standard library
  - Built-in HTTP server capabilities
  - Excellent performance characteristics
  - Easy deployment with single binary

### Frontend
- **HTML5**: Semantic markup for structured content
- **CSS3**: Modern styling with flexbox and grid layouts
- **JavaScript**: Client-side interactivity and form validation
  - DOM manipulation for dynamic content updates
  - Event handling for user interactions
  - Fetch API for AJAX requests

### Frameworks & Libraries
- **Bootstrap 5**: Responsive CSS framework for consistent UI elements
  - Mobile-first approach
  - Pre-built components (navigation, cards, forms)
  - Customizable with Sass variables
- **Gorilla Mux**: Powerful URL router and dispatcher for Go
  - Pattern matching for URLs
  - Named URL parameters
  - Request middleware support
- **Gorilla Sessions**: Session management for web applications
  - Secure cookie-based sessions
  - Flash messages support
  - Customizable session duration

### Database
- **SQLite3**: Self-contained, serverless SQL database engine
  - Zero configuration required
  - ACID-compliant transactions
  - Suitable for moderate traffic levels
  - Single file storage for easy backups

### Security
- **Bcrypt**: Industry-standard password hashing algorithm
  - Automatic salt generation
  - Adjustable work factor for future-proofing
  - Protection against rainbow table attacks

## Installation

### Prerequisites
- Go 1.16 or higher installed
- Git for version control
- Basic knowledge of terminal/command line

### Step-by-Step Setup

1. **Clone the repository**:
```bash
git clone https://github.com/yourusername/university-forum.git
cd university-forum
```

2. **Install dependencies**:
```bash
go mod download
```

3. **Initialize the database** (optional, happens automatically on first run):
```bash
# The database will be created automatically when you run the application
# If you want to reset the database, simply delete forum.db
```

4. **Run the application**:
```bash
go run main.go
```

5. **Access the forum** at `http://localhost:8080` in your web browser

### Configuration Options

The application can be configured through environment variables:

- `PORT`: Web server port (default: 8080)
- `DB_PATH`: Path to SQLite database file (default: ./forum.db)
- `SESSION_KEY`: Secret key for session encryption (default: auto-generated)

Example with custom port:
```bash
PORT=9000 go run main.go
```

## Project Structure

```
university-forum/
├── main.go              # Main application entry point, router setup, and DB initialization
├── handlers/            # Request handlers organized by functionality
│   ├── auth.go         # Authentication handlers (login, register, logout)
│   └── posts.go        # Post and comment handlers (create, view, search)
├── models/             # Data models and database operations
│   ├── user.go         # User model definitions and methods
│   └── post.go         # Post and comment model definitions and methods
├── static/             # Static assets served directly to client
│   ├── css/           
│   │   └── style.css   # Custom styling beyond Bootstrap
│   └── js/
│       └── main.js     # Client-side interactivity and enhancements
├── templates/          # HTML templates using Go's template system
│   ├── layout.html     # Base template with navigation and page structure
│   ├── index.html      # Home page displaying recent discussions
│   ├── login.html      # User authentication form
│   ├── register.html   # New user registration form
│   ├── create-post.html# Form for creating new discussions
│   ├── view-post.html  # Detailed view of a discussion and its comments
│   └── search.html     # Search interface and results display
├── go.mod              # Go module definition and dependencies
├── go.sum              # Dependency verification checksums
└── forum.db            # SQLite database file storing all application data
```

### Key Components

- **Router (main.go)**: Uses Gorilla Mux to define routes and map them to handler functions
- **Authentication System**: Handles user registration, login, and session management using secure cookies
- **Discussion System**: Manages the creation, display, and interaction with discussion posts
- **Search Engine**: Provides full-text search capabilities across all discussion content
- **Database Layer**: Manages data persistence using SQLite with proper indexing for performance
- **Template System**: Renders dynamic HTML using Go's built-in template engine with layout inheritance

## Contributing

We welcome contributions from the community! Here's how you can help improve this project:

### Development Workflow

1. **Fork the repository** on GitHub
2. **Create a feature branch** with a descriptive name:
```bash
git checkout -b feature/search-functionality
```
3. **Make your changes** following our code style guidelines
4. **Write or update tests** if necessary
5. **Run tests** to ensure everything works:
```bash
go test ./...
```
6. **Commit your changes** with clear, descriptive messages:
```bash
git commit -m 'Add search functionality with keyword highlighting'
```
7. **Push to your branch**:
```bash
git push origin feature/search-functionality
```
8. **Open a Pull Request** on GitHub with a clear description of the changes

### Code Style Guidelines

- Follow standard Go formatting (use `gofmt` or `go fmt`)
- Use meaningful variable and function names
- Comment complex logic and public functions
- Follow MVC-like separation of concerns
- Write tests for new functionality

### Bug Reports

If you find a bug, please open an issue on GitHub with:
- A clear, descriptive title
- Detailed steps to reproduce the issue
- Expected vs. actual behavior
- Screenshots if applicable
- System information (OS, browser, etc.)

## License

This project is licensed under the MIT License - see the LICENSE file for details. 