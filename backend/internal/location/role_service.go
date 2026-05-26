package location

import (
	"context"
	"fmt"
	"time"

	"github.com/futurebuildai/gablexhardscape/pkg/middleware"
	"github.com/google/uuid"
)

type RoleService struct {
	repo     RoleRepository
	userRepo UserRepository
}

func NewRoleService(repo RoleRepository, userRepo UserRepository) *RoleService {
	return &RoleService{
		repo:     repo,
		userRepo: userRepo,
	}
}

// GetOrResolveStaffRole fetches a user's role, auto-bootstrapping "General Manager" if they are admin/owner in claims.
func (s *RoleService) GetOrResolveStaffRole(ctx context.Context, userSub string, claims *middleware.UserClaims) (*StaffRole, error) {
	sr, err := s.repo.GetStaffRole(ctx, userSub)
	if err != nil {
		return nil, err
	}
	if sr != nil {
		return sr, nil
	}

	role := RoleInsideSales
	if claims != nil {
		isAdmin := claims.Role == "admin" || claims.Role == "owner"
		if !isAdmin {
			for _, r := range claims.Roles {
				if r == "admin" || r == "owner" {
					isAdmin = true
					break
				}
			}
		}
		if isAdmin {
			role = RoleGeneralManager
			newSr := StaffRole{
				UserSub:    userSub,
				Role:       RoleGeneralManager,
				AssignedAt: time.Now(),
				AssignedBy: "system:bootstrap",
			}
			if err := s.repo.SetStaffRole(ctx, newSr); err != nil {
				return nil, fmt.Errorf("bootstrap staff role: %w", err)
			}
			return &newSr, nil
		}
	}

	return &StaffRole{
		UserSub:    userSub,
		Role:       role,
		AssignedAt: time.Now(),
		AssignedBy: "system:default",
	}, nil
}

// DecideApprovalRequest decides on a pending override request, enforcing strict role constraints.
func (s *RoleService) DecideApprovalRequest(ctx context.Context, id uuid.UUID, decidedBy string, status ApprovalRequestStatus, reason string, deciderClaims *middleware.UserClaims) error {
	req, err := s.repo.GetApprovalRequest(ctx, id)
	if err != nil {
		return err
	}

	deciderRole, err := s.GetOrResolveStaffRole(ctx, decidedBy, deciderClaims)
	if err != nil {
		return fmt.Errorf("resolve decider role: %w", err)
	}

	// Retrieve the scoping tier for the decider's role
	tier := deciderRole.Role.GetPermissionTier()

	// Only Admin Tier (Super Admins) or Manager Tier (Scoped Branch Managers) can decide.
	if tier == TierAdmin {
		// Admin Tier bypasses all checks
	} else if tier == TierManager {
		// Manager Tier must belong to the branch of the request
		ok, err := s.userRepo.UserHasBranch(ctx, decidedBy, req.BranchID)
		if err != nil {
			return fmt.Errorf("verify decider branch access: %w", err)
		}
		if !ok {
			return fmt.Errorf("scoped manager does not belong to branch %s", req.BranchID)
		}
	} else {
		return fmt.Errorf("role '%s' (Staff Tier) is not authorized to decide approval requests", deciderRole.Role)
	}

	return s.repo.DecideApprovalRequest(ctx, id, decidedBy, status, reason)
}
