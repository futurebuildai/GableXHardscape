package purchase_order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// --- Replenishment Settings Repository ---

// GetReplenishmentSetting returns the per-product override, or nil if none exists.
func (r *Repository) GetReplenishmentSetting(ctx context.Context, productID uuid.UUID) (*ReplenishmentSetting, error) {
	const q = `
		SELECT id, product_id, min_safety_stock, velocity_window_days,
		       lead_time_override_days, created_at, updated_at
		FROM replenishment_settings
		WHERE product_id = $1
	`
	var rs ReplenishmentSetting
	err := r.db.GetExecutor(ctx).QueryRow(ctx, q, productID).Scan(
		&rs.ID, &rs.ProductID, &rs.MinSafetyStock, &rs.VelocityWindowDays,
		&rs.LeadTimeOverrideDays, &rs.CreatedAt, &rs.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get replenishment setting: %w", err)
	}
	return &rs, nil
}

// ListReplenishmentSettings returns all per-product overrides.
func (r *Repository) ListReplenishmentSettings(ctx context.Context) ([]ReplenishmentSetting, error) {
	const q = `
		SELECT id, product_id, min_safety_stock, velocity_window_days,
		       lead_time_override_days, created_at, updated_at
		FROM replenishment_settings
		ORDER BY updated_at DESC
	`
	rows, err := r.db.GetExecutor(ctx).Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list replenishment settings: %w", err)
	}
	defer rows.Close()

	var out []ReplenishmentSetting
	for rows.Next() {
		var rs ReplenishmentSetting
		if err := rows.Scan(
			&rs.ID, &rs.ProductID, &rs.MinSafetyStock, &rs.VelocityWindowDays,
			&rs.LeadTimeOverrideDays, &rs.CreatedAt, &rs.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan replenishment setting: %w", err)
		}
		out = append(out, rs)
	}
	return out, nil
}

// UpsertReplenishmentSetting creates or updates a per-product override.
func (r *Repository) UpsertReplenishmentSetting(ctx context.Context, rs *ReplenishmentSetting) error {
	const q = `
		INSERT INTO replenishment_settings (product_id, min_safety_stock, velocity_window_days, lead_time_override_days)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (product_id) DO UPDATE SET
			min_safety_stock = EXCLUDED.min_safety_stock,
			velocity_window_days = EXCLUDED.velocity_window_days,
			lead_time_override_days = EXCLUDED.lead_time_override_days,
			updated_at = now()
		RETURNING id, created_at, updated_at
	`
	return r.db.GetExecutor(ctx).QueryRow(ctx, q,
		rs.ProductID, rs.MinSafetyStock, rs.VelocityWindowDays, rs.LeadTimeOverrideDays,
	).Scan(&rs.ID, &rs.CreatedAt, &rs.UpdatedAt)
}

// --- Replenishment Draft Repository ---

// CreateReplenishmentDraft inserts a new draft row linked to a DRAFT PO.
func (r *Repository) CreateReplenishmentDraft(ctx context.Context, d *ReplenishmentDraft) error {
	const q = `
		INSERT INTO replenishment_drafts (po_id, vendor_id, status, confidence, total_lines, total_est_cost)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, generated_at
	`
	return r.db.GetExecutor(ctx).QueryRow(ctx, q,
		d.POID, d.VendorID, d.Status, d.Confidence, d.TotalLines, d.TotalEstCost,
	).Scan(&d.ID, &d.GeneratedAt)
}

// ListPendingDrafts returns all PENDING_REVIEW drafts with vendor names.
func (r *Repository) ListPendingDrafts(ctx context.Context) ([]ReplenishmentDraft, error) {
	const q = `
		SELECT d.id, d.po_id, d.vendor_id, COALESCE(v.name, ''), d.status,
		       d.generated_at, d.reviewed_by, d.reviewed_at, COALESCE(d.notes, ''),
		       d.confidence, d.total_lines, d.total_est_cost
		FROM replenishment_drafts d
		LEFT JOIN vendors v ON v.id = d.vendor_id
		WHERE d.status = 'PENDING_REVIEW'
		ORDER BY d.generated_at DESC
	`
	rows, err := r.db.GetExecutor(ctx).Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list pending drafts: %w", err)
	}
	defer rows.Close()

	var out []ReplenishmentDraft
	for rows.Next() {
		var d ReplenishmentDraft
		if err := rows.Scan(
			&d.ID, &d.POID, &d.VendorID, &d.VendorName, &d.Status,
			&d.GeneratedAt, &d.ReviewedBy, &d.ReviewedAt, &d.Notes,
			&d.Confidence, &d.TotalLines, &d.TotalEstCost,
		); err != nil {
			return nil, fmt.Errorf("scan replenishment draft: %w", err)
		}
		out = append(out, d)
	}
	return out, nil
}

// UpdateDraftStatus transitions a draft's status and records the reviewer.
func (r *Repository) UpdateDraftStatus(ctx context.Context, draftID uuid.UUID, status ReplenishmentDraftStatus, reviewerID uuid.UUID, notes string) error {
	now := time.Now()
	const q = `
		UPDATE replenishment_drafts
		SET status = $2, reviewed_by = $3, reviewed_at = $4, notes = COALESCE(NULLIF($5, ''), notes)
		WHERE id = $1
	`
	_, err := r.db.GetExecutor(ctx).Exec(ctx, q, draftID, status, reviewerID, now, notes)
	return err
}

// UpdateDraftTotals recalculates the total_lines and total_est_cost after edits.
func (r *Repository) UpdateDraftTotals(ctx context.Context, draftID uuid.UUID, totalLines int, totalEstCost float64) error {
	const q = `UPDATE replenishment_drafts SET total_lines = $2, total_est_cost = $3 WHERE id = $1`
	_, err := r.db.GetExecutor(ctx).Exec(ctx, q, draftID, totalLines, totalEstCost)
	return err
}

// GetDraft returns a single draft by ID.
func (r *Repository) GetDraft(ctx context.Context, draftID uuid.UUID) (*ReplenishmentDraft, error) {
	const q = `
		SELECT d.id, d.po_id, d.vendor_id, COALESCE(v.name, ''), d.status,
		       d.generated_at, d.reviewed_by, d.reviewed_at, COALESCE(d.notes, ''),
		       d.confidence, d.total_lines, d.total_est_cost
		FROM replenishment_drafts d
		LEFT JOIN vendors v ON v.id = d.vendor_id
		WHERE d.id = $1
	`
	var d ReplenishmentDraft
	err := r.db.GetExecutor(ctx).QueryRow(ctx, q, draftID).Scan(
		&d.ID, &d.POID, &d.VendorID, &d.VendorName, &d.Status,
		&d.GeneratedAt, &d.ReviewedBy, &d.ReviewedAt, &d.Notes,
		&d.Confidence, &d.TotalLines, &d.TotalEstCost,
	)
	if err != nil {
		return nil, fmt.Errorf("get draft: %w", err)
	}
	return &d, nil
}
