// Types for the A-12 Suggested-PO / Replenishment Procurement Dashboard

export type ReplenishmentDraftStatus = 'PENDING_REVIEW' | 'APPROVED' | 'REJECTED' | 'EXPIRED';

export interface ReplenishmentDraft {
    id: string;
    po_id: string;
    vendor_id: string;
    vendor_name?: string;
    status: ReplenishmentDraftStatus;
    generated_at: string;
    reviewed_by?: string;
    reviewed_at?: string;
    notes?: string;
    confidence: number;        // 0–100
    total_lines: number;
    total_est_cost: number;
    po?: import('./purchaseOrder').PurchaseOrder;
    recommendations?: import('./purchaseOrder').PurchaseRecommendation[];
}

export interface DashboardSummary {
    pending_count: number;
    total_est_cost: number;
    vendor_groups: VendorDraftGroup[];
}

export interface VendorDraftGroup {
    vendor_id: string;
    vendor_name: string;
    lead_time_days: number;
    drafts: ReplenishmentDraft[];
    total_cost: number;
}

export interface EditDraftRequest {
    lines: EditDraftLine[];
    notes?: string;
}

export interface EditDraftLine {
    line_id: string;
    quantity?: number;         // null/undefined = remove line
}

export interface ReplenishmentSetting {
    id: string;
    product_id: string;
    min_safety_stock: number;
    velocity_window_days: number;
    lead_time_override_days?: number;
    created_at: string;
    updated_at: string;
}
