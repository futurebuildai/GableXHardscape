package purchase_order

import (
	"encoding/json"
	"net/http"

	"github.com/futurebuildai/gablexhardscape/pkg/httputil"
	"github.com/google/uuid"
)

// --- Procurement Dashboard Handlers ---

// HandleGetProcurementDashboard returns all PENDING_REVIEW drafts grouped by vendor.
// GET /api/v1/purchase-orders/procurement-dashboard
func (h *Handler) HandleGetProcurementDashboard(w http.ResponseWriter, r *http.Request) {
	summary, err := h.service.GetProcurementDashboard(r.Context())
	if err != nil {
		httputil.RespondError(w, r, "failed to load procurement dashboard", http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// HandleGenerateProcurementDrafts triggers on-demand draft generation from recommendations.
// POST /api/v1/purchase-orders/procurement-dashboard/generate
func (h *Handler) HandleGenerateProcurementDrafts(w http.ResponseWriter, r *http.Request) {
	drafts, err := h.service.GenerateProcurementDrafts(r.Context())
	if err != nil {
		httputil.RespondError(w, r, "failed to generate procurement drafts", http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(drafts)
}

// HandleEditProcurementDraft modifies line quantities or notes on a pending draft.
// PATCH /api/v1/purchase-orders/procurement-dashboard/{draft_id}
func (h *Handler) HandleEditProcurementDraft(w http.ResponseWriter, r *http.Request) {
	draftID, err := uuid.Parse(r.PathValue("draft_id"))
	if err != nil {
		httputil.RespondError(w, r, "Invalid draft_id", http.StatusBadRequest, err)
		return
	}

	var req EditDraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondError(w, r, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	draft, err := h.service.EditProcurementDraft(r.Context(), draftID, req)
	if err != nil {
		httputil.RespondError(w, r, "failed to edit procurement draft", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(draft)
}

// HandleApproveProcurementDraft approves a draft and sends the PO.
// POST /api/v1/purchase-orders/procurement-dashboard/{draft_id}/approve
func (h *Handler) HandleApproveProcurementDraft(w http.ResponseWriter, r *http.Request) {
	draftID, err := uuid.Parse(r.PathValue("draft_id"))
	if err != nil {
		httputil.RespondError(w, r, "Invalid draft_id", http.StatusBadRequest, err)
		return
	}

	// TODO: extract reviewer ID from auth context
	reviewerID := uuid.Nil

	if err := h.service.ApproveProcurementDraft(r.Context(), draftID, reviewerID); err != nil {
		httputil.RespondError(w, r, "failed to approve procurement draft", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "approved"})
}

// HandleRejectProcurementDraft rejects a draft without sending the PO.
// POST /api/v1/purchase-orders/procurement-dashboard/{draft_id}/reject
func (h *Handler) HandleRejectProcurementDraft(w http.ResponseWriter, r *http.Request) {
	draftID, err := uuid.Parse(r.PathValue("draft_id"))
	if err != nil {
		httputil.RespondError(w, r, "Invalid draft_id", http.StatusBadRequest, err)
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	// TODO: extract reviewer ID from auth context
	reviewerID := uuid.Nil

	if err := h.service.RejectProcurementDraft(r.Context(), draftID, reviewerID, req.Notes); err != nil {
		httputil.RespondError(w, r, "failed to reject procurement draft", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "rejected"})
}

// --- Replenishment Settings Handlers ---

// HandleListReplenishmentSettings returns all per-product overrides.
// GET /api/v1/purchase-orders/replenishment-settings
func (h *Handler) HandleListReplenishmentSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.service.ListReplenishmentSettings(r.Context())
	if err != nil {
		httputil.RespondError(w, r, "failed to list replenishment settings", http.StatusInternalServerError, err)
		return
	}
	if settings == nil {
		settings = []ReplenishmentSetting{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

// HandleUpsertReplenishmentSetting creates or updates a per-product override.
// PUT /api/v1/purchase-orders/replenishment-settings/{product_id}
func (h *Handler) HandleUpsertReplenishmentSetting(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(r.PathValue("product_id"))
	if err != nil {
		httputil.RespondError(w, r, "Invalid product_id", http.StatusBadRequest, err)
		return
	}

	var setting ReplenishmentSetting
	if err := json.NewDecoder(r.Body).Decode(&setting); err != nil {
		httputil.RespondError(w, r, "Invalid request body", http.StatusBadRequest, err)
		return
	}
	setting.ProductID = productID

	result, err := h.service.UpsertReplenishmentSetting(r.Context(), productID, setting)
	if err != nil {
		httputil.RespondError(w, r, "failed to upsert replenishment setting", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// RegisterProcurementRoutes registers the procurement dashboard and replenishment settings routes.
// Called from the existing RegisterRoutes method.
func (h *Handler) RegisterProcurementRoutes(mux *http.ServeMux, roleGuard ...func(http.Handler) http.Handler) {
	guard := func(handler http.HandlerFunc) http.HandlerFunc {
		if len(roleGuard) > 0 && roleGuard[0] != nil {
			return func(w http.ResponseWriter, r *http.Request) {
				roleGuard[0](handler).ServeHTTP(w, r)
			}
		}
		return handler
	}

	mux.HandleFunc("GET /api/v1/purchase-orders/procurement-dashboard", guard(h.HandleGetProcurementDashboard))
	mux.HandleFunc("POST /api/v1/purchase-orders/procurement-dashboard/generate", guard(h.HandleGenerateProcurementDrafts))
	mux.HandleFunc("PATCH /api/v1/purchase-orders/procurement-dashboard/{draft_id}", guard(h.HandleEditProcurementDraft))
	mux.HandleFunc("POST /api/v1/purchase-orders/procurement-dashboard/{draft_id}/approve", guard(h.HandleApproveProcurementDraft))
	mux.HandleFunc("POST /api/v1/purchase-orders/procurement-dashboard/{draft_id}/reject", guard(h.HandleRejectProcurementDraft))
	mux.HandleFunc("GET /api/v1/purchase-orders/replenishment-settings", guard(h.HandleListReplenishmentSettings))
	mux.HandleFunc("PUT /api/v1/purchase-orders/replenishment-settings/{product_id}", guard(h.HandleUpsertReplenishmentSetting))
}
