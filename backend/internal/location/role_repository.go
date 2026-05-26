package location

import (
	"context"
	"errors"
	"fmt"

	"github.com/futurebuildai/gablexhardscape/pkg/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresRoleRepository struct {
	db *database.DB
}

func NewRoleRepository(db *database.DB) *PostgresRoleRepository {
	return &PostgresRoleRepository{db: db}
}

// GetStaffRole fetches a user's role from the DB.
func (r *PostgresRoleRepository) GetStaffRole(ctx context.Context, userSub string) (*StaffRole, error) {
	query := `SELECT user_sub, role, assigned_at, COALESCE(assigned_by, '') FROM staff_roles WHERE user_sub = $1`
	var sr StaffRole
	err := r.db.GetExecutor(ctx).QueryRow(ctx, query, userSub).Scan(&sr.UserSub, &sr.Role, &sr.AssignedAt, &sr.AssignedBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found is not an error, will fall back
		}
		return nil, fmt.Errorf("get staff role: %w", err)
	}
	return &sr, nil
}

// SetStaffRole inserts or updates a user's role.
func (r *PostgresRoleRepository) SetStaffRole(ctx context.Context, sr StaffRole) error {
	query := `
		INSERT INTO staff_roles (user_sub, role, assigned_by)
		VALUES ($1, $2, NULLIF($3, ''))
		ON CONFLICT (user_sub) DO UPDATE
		SET role = EXCLUDED.role,
		    assigned_by = EXCLUDED.assigned_by,
		    assigned_at = NOW()
	`
	_, err := r.db.GetExecutor(ctx).Exec(ctx, query, sr.UserSub, sr.Role, sr.AssignedBy)
	if err != nil {
		return fmt.Errorf("set staff role: %w", err)
	}
	return nil
}

// ListStaffRoles lists all staff roles in the system.
func (r *PostgresRoleRepository) ListStaffRoles(ctx context.Context) ([]StaffRole, error) {
	query := `SELECT user_sub, role, assigned_at, COALESCE(assigned_by, '') FROM staff_roles ORDER BY assigned_at DESC`
	rows, err := r.db.GetExecutor(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list staff roles: %w", err)
	}
	defer rows.Close()

	var out []StaffRole
	for rows.Next() {
		var sr StaffRole
		if err := rows.Scan(&sr.UserSub, &sr.Role, &sr.AssignedAt, &sr.AssignedBy); err != nil {
			return nil, fmt.Errorf("scan staff role: %w", err)
		}
		out = append(out, sr)
	}
	return out, rows.Err()
}

// CreateApprovalRequest inserts a new pending permission override request.
func (r *PostgresRoleRepository) CreateApprovalRequest(ctx context.Context, req *PermissionApprovalRequest) error {
	query := `
		INSERT INTO permission_approval_requests (id, user_sub, branch_id, policy_type, details, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	if req.ID == uuid.Nil {
		req.ID = uuid.New()
	}
	req.Status = StatusPending
	_, err := r.db.GetExecutor(ctx).Exec(ctx, query, req.ID, req.UserSub, req.BranchID, req.PolicyType, req.Details, req.Status)
	if err != nil {
		return fmt.Errorf("create approval request: %w", err)
	}
	return nil
}

// GetApprovalRequest retrieves a single approval request.
func (r *PostgresRoleRepository) GetApprovalRequest(ctx context.Context, id uuid.UUID) (*PermissionApprovalRequest, error) {
	query := `
		SELECT id, user_sub, branch_id, policy_type, details, status, requested_at, decided_at, COALESCE(decided_by, ''), COALESCE(rejection_reason, ''), created_at
		FROM permission_approval_requests
		WHERE id = $1
	`
	var req PermissionApprovalRequest
	var decidedBy, rejectionReason string
	err := r.db.GetExecutor(ctx).QueryRow(ctx, query, id).Scan(
		&req.ID, &req.UserSub, &req.BranchID, &req.PolicyType, &req.Details, &req.Status,
		&req.RequestedAt, &req.DecidedAt, &decidedBy, &rejectionReason, &req.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get approval request: %w", err)
	}
	req.DecidedBy = decidedBy
	req.RejectionReason = rejectionReason
	return &req, nil
}

// ListApprovalRequests returns chronological approval requests, optionally filtered by branch and status.
func (r *PostgresRoleRepository) ListApprovalRequests(ctx context.Context, branchID *uuid.UUID, status *string) ([]PermissionApprovalRequest, error) {
	query := `
		SELECT id, user_sub, branch_id, policy_type, details, status, requested_at, decided_at, COALESCE(decided_by, ''), COALESCE(rejection_reason, ''), created_at
		FROM permission_approval_requests
		WHERE ($1::uuid IS NULL OR branch_id = $1)
		  AND ($2::text IS NULL OR status = $2)
		ORDER BY requested_at DESC
	`
	rows, err := r.db.GetExecutor(ctx).Query(ctx, query, branchID, status)
	if err != nil {
		return nil, fmt.Errorf("list approval requests: %w", err)
	}
	defer rows.Close()

	var out []PermissionApprovalRequest
	for rows.Next() {
		var req PermissionApprovalRequest
		var decidedBy, rejectionReason string
		err := rows.Scan(
			&req.ID, &req.UserSub, &req.BranchID, &req.PolicyType, &req.Details, &req.Status,
			&req.RequestedAt, &req.DecidedAt, &decidedBy, &rejectionReason, &req.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan approval request: %w", err)
		}
		req.DecidedBy = decidedBy
		req.RejectionReason = rejectionReason
		out = append(out, req)
	}
	return out, rows.Err()
}

// DecideApprovalRequest approves or rejects a pending request.
func (r *PostgresRoleRepository) DecideApprovalRequest(ctx context.Context, id uuid.UUID, decidedBy string, status ApprovalRequestStatus, reason string) error {
	query := `
		UPDATE permission_approval_requests
		SET status = $1,
		    decided_by = $2,
		    decided_at = NOW(),
		    rejection_reason = NULLIF($3, '')
		WHERE id = $4
	`
	tag, err := r.db.GetExecutor(ctx).Exec(ctx, query, status, decidedBy, reason, id)
	if err != nil {
		return fmt.Errorf("decide approval request: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
