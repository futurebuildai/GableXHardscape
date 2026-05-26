package feedback

import (
	"time"

	"github.com/google/uuid"
)

// Feedback represents a single feedback submission from the ERP or Partner Portal.
type Feedback struct {
	ID               uuid.UUID  `json:"id"`
	Source           string     `json:"source"`            // "ERP" or "PORTAL"
	Category         string     `json:"category"`          // Bug, UI/UX, Feature Request, Data Issue, Question, Other
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	PageURL          string     `json:"page_url,omitempty"`
	SubmittedByName  string     `json:"submitted_by_name,omitempty"`
	SubmittedByEmail string     `json:"submitted_by_email,omitempty"`
	UserID           *uuid.UUID `json:"user_id,omitempty"`
	Status           string     `json:"status"`   // NEW, ACKNOWLEDGED, IN_PROGRESS, RESOLVED, CLOSED
	Priority         string     `json:"priority"` // LOW, MEDIUM, HIGH, CRITICAL
	AdminNotes       string     `json:"admin_notes,omitempty"`
	ResolvedAt       *time.Time `json:"resolved_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ValidCategories enumerates the allowed category values.
var ValidCategories = map[string]bool{
	"Bug":             true,
	"UI/UX":           true,
	"Feature Request": true,
	"Data Issue":      true,
	"Question":        true,
	"Other":           true,
}

// ValidStatuses enumerates the allowed status values.
var ValidStatuses = map[string]bool{
	"NEW":          true,
	"ACKNOWLEDGED": true,
	"IN_PROGRESS":  true,
	"RESOLVED":     true,
	"CLOSED":       true,
}

// ValidPriorities enumerates the allowed priority values.
var ValidPriorities = map[string]bool{
	"LOW":      true,
	"MEDIUM":   true,
	"HIGH":     true,
	"CRITICAL": true,
}

// CreateFeedbackRequest is the payload for submitting new feedback.
type CreateFeedbackRequest struct {
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PageURL     string `json:"page_url,omitempty"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
}

// UpdateFeedbackRequest is the payload for admin status/priority updates.
type UpdateFeedbackRequest struct {
	Status     string `json:"status,omitempty"`
	Priority   string `json:"priority,omitempty"`
	AdminNotes string `json:"admin_notes,omitempty"`
}

// FeedbackListFilter holds query parameters for listing feedback.
type FeedbackListFilter struct {
	Status   string
	Category string
	Source   string
	Search   string
	Page     int
	Limit    int
}

// FeedbackListResponse wraps a paginated list of feedback items.
type FeedbackListResponse struct {
	Items      []Feedback     `json:"items"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	StatusCounts map[string]int `json:"status_counts"`
}
