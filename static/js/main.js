// Enable Bootstrap tooltips
document.addEventListener('DOMContentLoaded', function() {
    var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    var tooltipList = tooltipTriggerList.map(function(tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl);
    });
});

// Add confirmation for delete actions
document.addEventListener('click', function(e) {
    if (e.target && e.target.classList.contains('delete-confirm')) {
        if (!confirm('Are you sure you want to delete this item?')) {
            e.preventDefault();
        }
    }
});

// Auto-resize textareas
document.addEventListener('DOMContentLoaded', function() {
    const textareas = document.querySelectorAll('textarea');
    textareas.forEach(textarea => {
        textarea.addEventListener('input', function() {
            this.style.height = 'auto';
            this.style.height = (this.scrollHeight) + 'px';
        });
    });
});

// Real-time search suggestions
document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.querySelector('.search-input');
    const suggestionsContainer = document.querySelector('.search-suggestions');
    
    if (searchInput && suggestionsContainer) {
        let debounceTimer;
        
        searchInput.addEventListener('input', function() {
            clearTimeout(debounceTimer);
            const query = this.value.trim();
            
            if (query.length < 2) {
                suggestionsContainer.innerHTML = '';
                suggestionsContainer.classList.remove('active');
                return;
            }
            
            // Debounce to avoid too many requests
            debounceTimer = setTimeout(() => {
                fetch(`/api/search-suggestions?q=${encodeURIComponent(query)}`)
                    .then(response => response.json())
                    .then(data => {
                        suggestionsContainer.innerHTML = '';
                        
                        if (data.suggestions && data.suggestions.length > 0) {
                            data.suggestions.forEach(suggestion => {
                                const div = document.createElement('div');
                                div.className = 'suggestion-item';
                                div.innerHTML = `<a href="/search?q=${encodeURIComponent(suggestion)}">${highlightMatch(suggestion, query)}</a>`;
                                suggestionsContainer.appendChild(div);
                            });
                            suggestionsContainer.classList.add('active');
                        } else {
                            suggestionsContainer.classList.remove('active');
                        }
                    })
                    .catch(error => {
                        console.error('Error fetching suggestions:', error);
                    });
            }, 300);
        });
        
        // Hide suggestions when clicking outside
        document.addEventListener('click', function(e) {
            if (!searchInput.contains(e.target) && !suggestionsContainer.contains(e.target)) {
                suggestionsContainer.classList.remove('active');
            }
        });
    }
});

// Highlight matching text in search results
function highlightMatch(text, query) {
    const regex = new RegExp(`(${query})`, 'gi');
    return text.replace(regex, '<mark>$1</mark>');
}

// Form validation
document.addEventListener('DOMContentLoaded', function() {
    const forms = document.querySelectorAll('.needs-validation');
    
    Array.from(forms).forEach(form => {
        form.addEventListener('submit', event => {
            if (!form.checkValidity()) {
                event.preventDefault();
                event.stopPropagation();
            }
            
            form.classList.add('was-validated');
        }, false);
    });
});

// Lazy load comments for posts
document.addEventListener('DOMContentLoaded', function() {
    const commentSection = document.querySelector('.comments-section');
    const loadMoreBtn = document.querySelector('.load-more-comments');
    
    if (commentSection && loadMoreBtn) {
        let page = 1;
        
        loadMoreBtn.addEventListener('click', function() {
            const postId = this.dataset.postId;
            page++;
            
            fetch(`/post/${postId}/comments?page=${page}`)
                .then(response => response.json())
                .then(data => {
                    if (data.comments && data.comments.length > 0) {
                        const fragment = document.createDocumentFragment();
                        
                        data.comments.forEach(comment => {
                            const commentEl = document.createElement('div');
                            commentEl.className = 'comment';
                            commentEl.innerHTML = `
                                <div class="comment-header">
                                    <span class="comment-author">${comment.authorName}</span>
                                    <span class="comment-date">${comment.createdAt}</span>
                                </div>
                                <div class="comment-body">
                                    <p>${comment.content}</p>
                                </div>
                            `;
                            fragment.appendChild(commentEl);
                        });
                        
                        commentSection.appendChild(fragment);
                        
                        if (data.hasMore === false) {
                            loadMoreBtn.style.display = 'none';
                        }
                    } else {
                        loadMoreBtn.style.display = 'none';
                    }
                })
                .catch(error => {
                    console.error('Error loading comments:', error);
                });
        });
    }
});

// Dark mode toggle
document.addEventListener('DOMContentLoaded', function() {
    const darkModeToggle = document.querySelector('.dark-mode-toggle');
    const htmlElement = document.documentElement;
    
    // Check for saved theme preference or respect OS preference
    const savedTheme = localStorage.getItem('theme');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    
    // Set initial theme
    if (savedTheme === 'dark' || (savedTheme === null && prefersDark)) {
        htmlElement.setAttribute('data-bs-theme', 'dark');
        if (darkModeToggle) {
            darkModeToggle.innerHTML = '<i class="bi bi-sun"></i>';
        }
    } else {
        htmlElement.setAttribute('data-bs-theme', 'light');
        if (darkModeToggle) {
            darkModeToggle.innerHTML = '<i class="bi bi-moon"></i>';
        }
    }
    
    // Toggle theme when button is clicked
    if (darkModeToggle) {
        darkModeToggle.addEventListener('click', function() {
            const currentTheme = htmlElement.getAttribute('data-bs-theme');
            const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
            
            htmlElement.setAttribute('data-bs-theme', newTheme);
            localStorage.setItem('theme', newTheme);
            
            // Update icon
            darkModeToggle.innerHTML = newTheme === 'dark' 
                ? '<i class="bi bi-sun"></i>' 
                : '<i class="bi bi-moon"></i>';
        });
    }
});

// Notification system
document.addEventListener('DOMContentLoaded', function() {
    const notificationBell = document.querySelector('.notification-bell');
    const notificationDropdown = document.querySelector('.notification-dropdown');
    
    if (notificationBell && notificationDropdown) {
        // Check for new notifications periodically
        function checkNotifications() {
            fetch('/api/notifications/unread-count')
                .then(response => response.json())
                .then(data => {
                    const counter = notificationBell.querySelector('.counter');
                    if (data.count > 0) {
                        if (!counter) {
                            const newCounter = document.createElement('span');
                            newCounter.className = 'counter badge bg-danger';
                            newCounter.textContent = data.count;
                            notificationBell.appendChild(newCounter);
                        } else {
                            counter.textContent = data.count;
                        }
                    } else if (counter) {
                        counter.remove();
                    }
                })
                .catch(error => {
                    console.error('Error checking notifications:', error);
                });
        }
        
        // Load notifications when clicked
        notificationBell.addEventListener('click', function(e) {
            e.preventDefault();
            
            fetch('/api/notifications')
                .then(response => response.json())
                .then(data => {
                    notificationDropdown.innerHTML = '';
                    
                    if (data.notifications && data.notifications.length > 0) {
                        data.notifications.forEach(notification => {
                            const item = document.createElement('a');
                            item.className = notification.read ? 'dropdown-item' : 'dropdown-item unread';
                            item.href = notification.link;
                            item.innerHTML = `
                                <div class="notification-item">
                                    <div class="notification-content">${notification.message}</div>
                                    <div class="notification-time">${notification.timeAgo}</div>
                                </div>
                            `;
                            
                            // Mark as read when clicked
                            item.addEventListener('click', function() {
                                fetch(`/api/notifications/${notification.id}/mark-read`, {
                                    method: 'POST'
                                }).catch(error => {
                                    console.error('Error marking notification as read:', error);
                                });
                            });
                            
                            notificationDropdown.appendChild(item);
                        });
                    } else {
                        const emptyItem = document.createElement('div');
                        emptyItem.className = 'dropdown-item text-center';
                        emptyItem.textContent = 'No notifications';
                        notificationDropdown.appendChild(emptyItem);
                    }
                    
                    // Mark all as read button
                    const markAllBtn = document.createElement('div');
                    markAllBtn.className = 'dropdown-item text-center mark-all-read';
                    markAllBtn.textContent = 'Mark all as read';
                    markAllBtn.addEventListener('click', function() {
                        fetch('/api/notifications/mark-all-read', {
                            method: 'POST'
                        })
                        .then(() => {
                            checkNotifications();
                            const unreadItems = notificationDropdown.querySelectorAll('.unread');
                            unreadItems.forEach(item => item.classList.remove('unread'));
                        })
                        .catch(error => {
                            console.error('Error marking all notifications as read:', error);
                        });
                    });
                    
                    notificationDropdown.appendChild(markAllBtn);
                })
                .catch(error => {
                    console.error('Error loading notifications:', error);
                });
        });
        
        // Initial check and periodic updates
        checkNotifications();
        setInterval(checkNotifications, 60000); // Check every minute
    }
}); 