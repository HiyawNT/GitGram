package models

// PushEvent represents the GitHub webhook push event payload
type PushEvent struct {
	Ref        string     `json:"ref"`
	Repository Repository `json:"repository"`
	Pusher     Pusher     `json:"pusher"`
	Commits    []Commit   `json:"commits"`
}

// Repository contains information about the GitHub repository
type Repository struct {
	FullName string `json:"full_name"`
	Name     string `json:"name"`
	HTMLURL  string `json:"html_url"`
}

// Pusher contains information about the person who pushed the commits
type Pusher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Commit represents a single commit in the push event
type Commit struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	URL     string `json:"url"`
	Author  Author `json:"author"`
}

// Author contains information about the commit author
type Author struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}