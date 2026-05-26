package feedback

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// Service provides business logic for the feedback module.
type Service struct {
	repo     *Repository
	notifier *Notifier
	logger   *slog.Logger
}

// NewService creates a new feedback service.
func NewService(repo *Repository, notifier *Notifier, logger *slog.Logger) *Service {
	return &Service{
		repo:     repo,
		notifier: notifier,
		logger:   logger,
	}
}

// SubmitFeedback validates and creates a new feedback item, then triggers
// a Google Chat notification.
func (s *Service) SubmitFeedback(ctx context.Context, req CreateFeedbackRequest, source string, userID *uuid.UUID) (*Feedback, error) {
	// Validate required fields.
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.Description == "" {
		return nil, fmt.Errorf("description is required")
	}
	if !ValidCategories[req.Category] {
		return nil, fmt.Errorf("invalid category: %s", req.Category)
	}

	fb := &Feedback{
		Source:           source,
		Category:         req.Category,
		Title:            req.Title,
		Description:      req.Description,
		PageURL:          req.PageURL,
		SubmittedByName:  req.Name,
		SubmittedByEmail: req.Email,
		UserID:           userID,
		Status:           "NEW",
		Priority:         "MEDIUM",
	}

	created, err := s.repo.Create(ctx, fb)
	if err != nil {
		return nil, fmt.Errorf("failed to submit feedback: %w", err)
	}

	// Fire-and-forget notification to Google Chat.
	if s.notifier != nil {
		s.notifier.NotifyNewFeedback(created)
	}

	s.logger.Info("New feedback submitted",
		"id", created.ID,
		"source", source,
		"category", req.Category,
		"title", req.Title,
	)

	return created, nil
}

// ListFeedback returns a paginated, filtered list of feedback items with
// aggregate status counts.
func (s *Service) ListFeedback(ctx context.Context, filter FeedbackListFilter) (*FeedbackListResponse, error) {
	items, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	counts, err := s.repo.CountByStatus(ctx)
	if err != nil {
		s.logger.Warn("Failed to count feedback by status", "error", err)
		counts = make(map[string]int)
	}

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	return &FeedbackListResponse{
		Items:        items,
		Total:        total,
		Page:         filter.Page,
		Limit:        limit,
		StatusCounts: counts,
	}, nil
}

// GetFeedback retrieves a single feedback item by ID.
func (s *Service) GetFeedback(ctx context.Context, id uuid.UUID) (*Feedback, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateFeedback allows admins to update status, priority, or admin notes.
func (s *Service) UpdateFeedback(ctx context.Context, id uuid.UUID, req UpdateFeedbackRequest) (*Feedback, error) {
	if req.Status != "" && !ValidStatuses[req.Status] {
		return nil, fmt.Errorf("invalid status: %s", req.Status)
	}
	if req.Priority != "" && !ValidPriorities[req.Priority] {
		return nil, fmt.Errorf("invalid priority: %s", req.Priority)
	}
	return s.repo.Update(ctx, id, req)
}
