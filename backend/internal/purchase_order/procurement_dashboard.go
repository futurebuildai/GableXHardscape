package purchase_order

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// --- Models ---

// ReplenishmentSetting holds per-product overrides for the dynamic reorder formula.
// When present, these values take precedence over the global RecommendationConfig.
type ReplenishmentSetting struct {
	ID                   uuid.UUID `json:"id"`
	ProductID            uuid.UUID `json:"product_id"`
	MinSafetyStock       float64   `json:"min_safety_stock"`
	VelocityWindowDays   int       `json:"velocity_window_days"`
	LeadTimeOverrideDays *float64  `json:"lead_time_override_days,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// ReplenishmentDraftStatus enumerates the lifecycle states of a procurement draft.
type ReplenishmentDraftStatus string

const (
	DraftStatusPendingReview ReplenishmentDraftStatus = "PENDING_REVIEW"
	DraftStatusApproved      ReplenishmentDraftStatus = "APPROVED"
	DraftStatusRejected      ReplenishmentDraftStatus = "REJECTED"
	DraftStatusExpired       ReplenishmentDraftStatus = "EXPIRED"
)

// ReplenishmentDraft wraps a DRAFT PO with procurement review metadata.
type ReplenishmentDraft struct {
	ID           uuid.UUID                `json:"id"`
	POID         uuid.UUID                `json:"po_id"`
	VendorID     uuid.UUID                `json:"vendor_id"`
	VendorName   string                   `json:"vendor_name,omitempty"`
	Status       ReplenishmentDraftStatus `json:"status"`
	GeneratedAt  time.Time                `json:"generated_at"`
	ReviewedBy   *uuid.UUID               `json:"reviewed_by,omitempty"`
	ReviewedAt   *time.Time               `json:"reviewed_at,omitempty"`
	Notes        string                   `json:"notes,omitempty"`
	Confidence   float64                  `json:"confidence"`
	TotalLines   int                      `json:"total_lines"`
	TotalEstCost float64                  `json:"total_est_cost"`

	// Joined data — populated when loading for the dashboard
	PO              *PurchaseOrder         `json:"po,omitempty"`
	Recommendations []PurchaseRecommendation `json:"recommendations,omitempty"`
}

// DashboardSummary is the top-level response for GET /procurement-dashboard.
type DashboardSummary struct {
	PendingCount   int                  `json:"pending_count"`
	TotalEstCost   float64              `json:"total_est_cost"`
	VendorGroups   []VendorDraftGroup   `json:"vendor_groups"`
}

// VendorDraftGroup clusters drafts by vendor for the dashboard view.
type VendorDraftGroup struct {
	VendorID     uuid.UUID            `json:"vendor_id"`
	VendorName   string               `json:"vendor_name"`
	LeadTimeDays float64              `json:"lead_time_days"`
	Drafts       []ReplenishmentDraft `json:"drafts"`
	TotalCost    float64              `json:"total_cost"`
}

// EditDraftRequest is the PATCH body for modifying a draft's lines.
type EditDraftRequest struct {
	Lines []EditDraftLine `json:"lines"`
	Notes *string         `json:"notes,omitempty"`
}

// EditDraftLine represents a quantity edit or removal on a draft PO line.
type EditDraftLine struct {
	LineID   uuid.UUID `json:"line_id"`
	Quantity *float64  `json:"quantity,omitempty"` // nil = remove line
}

// --- Service Methods (stubs) ---

// GenerateProcurementDrafts runs the recommendation engine, groups results by
// vendor, creates DRAFT POs + replenishment_drafts rows, and returns the
// generated drafts for the dashboard. This is the human-in-the-loop entry
// point — no POs are sent without explicit approval.
func (s *Service) GenerateProcurementDrafts(ctx context.Context) ([]ReplenishmentDraft, error) {
	// TODO: implement
	// 1. Call RecommendationService.GenerateRecommendations()
	// 2. Group recommendations by vendor
	// 3. For each vendor group:
	//    a. Create a DRAFT PO via CreateManualPOFromHandler (source=SUGGESTED)
	//    b. Insert a replenishment_drafts row with PENDING_REVIEW status
	//    c. Compute confidence score based on data quality (velocity coverage, lead time known, etc.)
	// 4. Return the created drafts
	return nil, nil
}

// GetProcurementDashboard loads all PENDING_REVIEW drafts grouped by vendor
// with their underlying PO lines and recommendation detail.
func (s *Service) GetProcurementDashboard(ctx context.Context) (*DashboardSummary, error) {
	// TODO: implement
	// 1. Query replenishment_drafts WHERE status = 'PENDING_REVIEW'
	// 2. Join with purchase_orders + purchase_order_lines
	// 3. Join with vendors for name + lead time
	// 4. Group by vendor_id
	// 5. Return DashboardSummary
	return nil, nil
}

// EditProcurementDraft modifies line quantities or notes on a pending draft.
func (s *Service) EditProcurementDraft(ctx context.Context, draftID uuid.UUID, req EditDraftRequest) (*ReplenishmentDraft, error) {
	// TODO: implement
	// 1. Verify draft status == PENDING_REVIEW
	// 2. For each line edit: update quantity or remove line from PO
	// 3. Recalculate total_est_cost and total_lines on the draft
	// 4. Update notes if provided
	// 5. Return updated draft
	return nil, nil
}

// ApproveProcurementDraft approves a draft, transitioning the underlying PO
// from DRAFT to SENT. This is the "Say Yes" confirmation gate.
func (s *Service) ApproveProcurementDraft(ctx context.Context, draftID uuid.UUID, reviewerID uuid.UUID) error {
	// TODO: implement
	// 1. Verify draft status == PENDING_REVIEW
	// 2. Set draft status = APPROVED, reviewed_by, reviewed_at
	// 3. Call SubmitPO on the underlying PO (transitions DRAFT → SENT, sends EDI)
	return nil
}

// RejectProcurementDraft marks a draft as rejected without sending the PO.
func (s *Service) RejectProcurementDraft(ctx context.Context, draftID uuid.UUID, reviewerID uuid.UUID, notes string) error {
	// TODO: implement
	// 1. Verify draft status == PENDING_REVIEW
	// 2. Set draft status = REJECTED, reviewed_by, reviewed_at, notes
	// 3. Keep underlying PO as DRAFT for potential manual handling
	return nil
}

// --- Replenishment Settings ---

// ListReplenishmentSettings returns all per-product overrides.
func (s *Service) ListReplenishmentSettings(ctx context.Context) ([]ReplenishmentSetting, error) {
	// TODO: implement
	return nil, nil
}

// UpsertReplenishmentSetting creates or updates per-product replenishment overrides.
func (s *Service) UpsertReplenishmentSetting(ctx context.Context, productID uuid.UUID, setting ReplenishmentSetting) (*ReplenishmentSetting, error) {
	// TODO: implement
	return nil, nil
}
