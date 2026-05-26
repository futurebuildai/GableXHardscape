package feedback

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/futurebuildai/gablexhardscape/pkg/database"
	"github.com/google/uuid"
)

// Repository provides data access for the feedback table.
type Repository struct {
	db *database.DB
}

// NewRepository creates a new feedback repository.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new feedback row and returns the created record.
func (r *Repository) Create(ctx context.Context, fb *Feedback) (*Feedback, error) {
	query := `
		INSERT INTO feedback (source, category, title, description, page_url,
			submitted_by_name, submitted_by_email, user_id, status, priority)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, source, category, title, description, page_url,
			COALESCE(submitted_by_name, ''), COALESCE(submitted_by_email, ''),
			user_id, status, priority, COALESCE(admin_notes, ''),
			resolved_at, created_at, updated_at
	`
	var result Feedback
	err := r.db.GetExecutor(ctx).QueryRow(ctx, query,
		fb.Source, fb.Category, fb.Title, fb.Description, fb.PageURL,
		fb.SubmittedByName, fb.SubmittedByEmail, fb.UserID, fb.Status, fb.Priority,
	).Scan(
		&result.ID, &result.Source, &result.Category, &result.Title, &result.Description,
		&result.PageURL, &result.SubmittedByName, &result.SubmittedByEmail,
		&result.UserID, &result.Status, &result.Priority, &result.AdminNotes,
		&result.ResolvedAt, &result.CreatedAt, &result.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create feedback: %w", err)
	}
	return &result, nil
}

// GetByID retrieves a single feedback item by UUID.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Feedback, error) {
	query := `
		SELECT id, source, category, title, description, COALESCE(page_url, ''),
			COALESCE(submitted_by_name, ''), COALESCE(submitted_by_email, ''),
			user_id, status, priority, COALESCE(admin_notes, ''),
			resolved_at, created_at, updated_at
		FROM feedback WHERE id = $1
	`
	var fb Feedback
	err := r.db.GetExecutor(ctx).QueryRow(ctx, query, id).Scan(
		&fb.ID, &fb.Source, &fb.Category, &fb.Title, &fb.Description,
		&fb.PageURL, &fb.SubmittedByName, &fb.SubmittedByEmail,
		&fb.UserID, &fb.Status, &fb.Priority, &fb.AdminNotes,
		&fb.ResolvedAt, &fb.CreatedAt, &fb.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("feedback not found: %w", err)
	}
	return &fb, nil
}

// List returns a paginated, filtered list of feedback items.
func (r *Repository) List(ctx context.Context, filter FeedbackListFilter) ([]Feedback, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if filter.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Category != "" {
		where = append(where, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, filter.Category)
		argIdx++
	}
	if filter.Source != "" {
		where = append(where, fmt.Sprintf("source = $%d", argIdx))
		args = append(args, filter.Source)
		argIdx++
	}
	if filter.Search != "" {
		where = append(where, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	// Count total matching rows.
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM feedback WHERE %s", whereClause)
	var total int
	if err := r.db.GetExecutor(ctx).QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count feedback: %w", err)
	}

	// Paginated fetch.
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset := 0
	if filter.Page > 1 {
		offset = (filter.Page - 1) * limit
	}
	dataQuery := fmt.Sprintf(`
		SELECT id, source, category, title, description, COALESCE(page_url, ''),
			COALESCE(submitted_by_name, ''), COALESCE(submitted_by_email, ''),
			user_id, status, priority, COALESCE(admin_notes, ''),
			resolved_at, created_at, updated_at
		FROM feedback WHERE %s
		ORDER BY created_at DESC
		LIMIT %d OFFSET %d
	`, whereClause, limit, offset)

	rows, err := r.db.GetExecutor(ctx).Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list feedback: %w", err)
	}
	defer rows.Close()

	items := make([]Feedback, 0)
	for rows.Next() {
		var fb Feedback
		if err := rows.Scan(
			&fb.ID, &fb.Source, &fb.Category, &fb.Title, &fb.Description,
			&fb.PageURL, &fb.SubmittedByName, &fb.SubmittedByEmail,
			&fb.UserID, &fb.Status, &fb.Priority, &fb.AdminNotes,
			&fb.ResolvedAt, &fb.CreatedAt, &fb.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan feedback: %w", err)
		}
		items = append(items, fb)
	}

	return items, total, nil
}

// Update modifies an existing feedback item's status, priority, or admin notes.
func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateFeedbackRequest) (*Feedback, error) {
	sets := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIdx := 1

	if req.Status != "" {
		sets = append(sets, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, req.Status)
		argIdx++
		if req.Status == "RESOLVED" {
			sets = append(sets, fmt.Sprintf("resolved_at = $%d", argIdx))
			args = append(args, time.Now())
			argIdx++
		}
	}
	if req.Priority != "" {
		sets = append(sets, fmt.Sprintf("priority = $%d", argIdx))
		args = append(args, req.Priority)
		argIdx++
	}
	if req.AdminNotes != "" {
		sets = append(sets, fmt.Sprintf("admin_notes = $%d", argIdx))
		args = append(args, req.AdminNotes)
		argIdx++
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE feedback SET %s WHERE id = $%d
		RETURNING id, source, category, title, description, COALESCE(page_url, ''),
			COALESCE(submitted_by_name, ''), COALESCE(submitted_by_email, ''),
			user_id, status, priority, COALESCE(admin_notes, ''),
			resolved_at, created_at, updated_at
	`, strings.Join(sets, ", "), argIdx)

	var fb Feedback
	err := r.db.GetExecutor(ctx).QueryRow(ctx, query, args...).Scan(
		&fb.ID, &fb.Source, &fb.Category, &fb.Title, &fb.Description,
		&fb.PageURL, &fb.SubmittedByName, &fb.SubmittedByEmail,
		&fb.UserID, &fb.Status, &fb.Priority, &fb.AdminNotes,
		&fb.ResolvedAt, &fb.CreatedAt, &fb.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update feedback: %w", err)
	}
	return &fb, nil
}

// CountByStatus returns a map of status → count for dashboard badges.
func (r *Repository) CountByStatus(ctx context.Context) (map[string]int, error) {
	query := `SELECT status, COUNT(*) FROM feedback GROUP BY status`
	rows, err := r.db.GetExecutor(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to count feedback by status: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return counts, nil
}
