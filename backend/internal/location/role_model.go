package location

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// StaffRoleType represents the granular staff roles from the client Permissions Hierarchy.
type StaffRoleType string

const (
	RoleGeneralManager      StaffRoleType = "General Manager"
	RoleFinancialController StaffRoleType = "Financial Controller"
	RoleBranchManager       StaffRoleType = "Branch Manager"
	RoleProcurementManager  StaffRoleType = "Procurement Manager"
	RoleSalesManager        StaffRoleType = "Sales Manager"
	RoleInsideSales         StaffRoleType = "Inside Sales"
	RoleOutsideSales        StaffRoleType = "Outside Sales"
	RoleYardManager         StaffRoleType = "Yard Manager"
	RoleYardTeam            StaffRoleType = "Yard Team"
	RoleLogisticsManager    StaffRoleType = "Logistics Manager"
	RoleDrivers             StaffRoleType = "Drivers"
	RoleHR                  StaffRoleType = "HR"
	RolePayablesReceivables StaffRoleType = "Payables/Receivables"
)

// PermissionTier represents the 3 structural permission tiers for database and middleware scoping.
type PermissionTier string

const (
	TierAdmin   PermissionTier = "ADMIN"
	TierManager PermissionTier = "MANAGER"
	TierStaff   PermissionTier = "STAFF"
)

// GetPermissionTier maps a granular staff role to its structural permission tier.
func (r StaffRoleType) GetPermissionTier() PermissionTier {
	switch r {
	case RoleGeneralManager, RoleFinancialController:
		return TierAdmin
	case RoleBranchManager, RoleProcurementManager, RoleSalesManager, RoleYardManager, RoleLogisticsManager, RoleHR:
		return TierManager
	default:
		return TierStaff
	}
}

// StaffRole represents the assigned role of a staff user in the ERP.
type StaffRole struct {
	UserSub    string        `json:"user_sub"`
	Role       StaffRoleType `json:"role"`
	AssignedAt time.Time     `json:"assigned_at"`
	AssignedBy string        `json:"assigned_by,omitempty"`
}

// ApprovalRequestStatus defines the state transitions of a permission override request.
type ApprovalRequestStatus string

const (
	StatusPending  ApprovalRequestStatus = "PENDING"
	StatusApproved ApprovalRequestStatus = "APPROVED"
	StatusRejected ApprovalRequestStatus = "REJECTED"
)

// PermissionApprovalRequest represents an override request when blocked by a policy.
type PermissionApprovalRequest struct {
	ID              uuid.UUID             `json:"id"`
	UserSub         string                `json:"user_sub"`
	BranchID        uuid.UUID             `json:"branch_id"`
	PolicyType      string                `json:"policy_type"` // 'MIN_MARGIN', 'CREDIT_LIMIT', etc.
	Details         json.RawMessage       `json:"details"`
	Status          ApprovalRequestStatus `json:"status"`
	RequestedAt     time.Time             `json:"requested_at"`
	DecidedAt       *time.Time            `json:"decided_at,omitempty"`
	DecidedBy       string                `json:"decided_by,omitempty"`
	RejectionReason string                `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time             `json:"created_at"`
}

// RoleRepository handles staff roles and approval request persistence.
type RoleRepository interface {
	GetStaffRole(ctx context.Context, userSub string) (*StaffRole, error)
	SetStaffRole(ctx context.Context, sr StaffRole) error
	ListStaffRoles(ctx context.Context) ([]StaffRole, error)

	CreateApprovalRequest(ctx context.Context, req *PermissionApprovalRequest) error
	GetApprovalRequest(ctx context.Context, id uuid.UUID) (*PermissionApprovalRequest, error)
	ListApprovalRequests(ctx context.Context, branchID *uuid.UUID, status *string) ([]PermissionApprovalRequest, error)
	DecideApprovalRequest(ctx context.Context, id uuid.UUID, decidedBy string, status ApprovalRequestStatus, reason string) error
}
