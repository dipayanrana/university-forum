// Wait for the DOM to be fully loaded
document.addEventListener('DOMContentLoaded', function() {
    // Enable Bootstrap tooltips
    var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl);
    });

    // Character counter for textareas
    const commentTextareas = document.querySelectorAll('textarea[name="content"]');
    if (commentTextareas.length > 0) {
        commentTextareas.forEach(textarea => {
            // Create and append counter element
            const counterDiv = document.createElement('div');
            counterDiv.className = 'text-muted small text-end mt-1';
            textarea.parentNode.appendChild(counterDiv);

            // Update counter on input
            textarea.addEventListener('input', function() {
                const remaining = this.value.length;
                counterDiv.textContent = `${remaining} characters`;

                // Visual feedback when approaching character limit
                if (remaining > 500) {
                    counterDiv.classList.add('text-warning');
                } else {
                    counterDiv.classList.remove('text-warning');
                }
            });

            // Trigger input event to initialize counter
            textarea.dispatchEvent(new Event('input'));
        });
    }

    // Add confirmation for delete actions
    const deleteButtons = document.querySelectorAll('.btn-delete');
    if (deleteButtons.length > 0) {
        deleteButtons.forEach(button => {
            button.addEventListener('click', function(e) {
                if (!confirm('Are you sure you want to delete this? This action cannot be undone.')) {
                    e.preventDefault();
                }
            });
        });
    }

    // Search highlighting
    const highlightSearchResults = () => {
        const urlParams = new URLSearchParams(window.location.search);
        const searchQuery = urlParams.get('q');
        
        if (searchQuery && searchQuery.length > 2) {
            const postContents = document.querySelectorAll('.card-text');
            const postTitles = document.querySelectorAll('.card-title');
            
            const regex = new RegExp(searchQuery, 'gi');
            
            // Function to highlight matches
            const highlightMatches = (element) => {
                const text = element.innerHTML;
                if (regex.test(text)) {
                    element.innerHTML = text.replace(
                        regex, 
                        match => `<span class="search-highlight">${match}</span>`
                    );
                }
            };
            
            // Apply to titles and content
            postTitles.forEach(highlightMatches);
            postContents.forEach(highlightMatches);
        }
    };
    
    // Run highlighting if we're on the search page
    if (window.location.pathname === '/search') {
        highlightSearchResults();
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