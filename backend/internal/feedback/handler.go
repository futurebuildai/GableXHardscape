package feedback

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/futurebuildai/gablexhardscape/pkg/httputil"
	"github.com/futurebuildai/gablexhardscape/pkg/middleware"
	"github.com/google/uuid"
)

// maxBodySize limits request body to 1MB.
const maxBodySize = 1 << 20

// Handler provides HTTP handlers for the feedback module.
type Handler struct {
	svc *Service
}

// NewHandler creates a new feedback handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterERPRoutes registers ERP-side feedback endpoints (JWT auth).
func (h *Handler) RegisterERPRoutes(mux *http.ServeMux, authMw func(http.Handler) http.Handler) {
	mux.Handle("POST /api/v1/feedback", authMw(http.HandlerFunc(h.HandleSubmitERP)))
	mux.Handle("GET /api/v1/feedback", authMw(http.HandlerFunc(h.HandleList)))
	mux.Handle("GET /api/v1/feedback/{id}", authMw(http.HandlerFunc(h.HandleGet)))
	mux.Handle("PUT /api/v1/feedback/{id}", authMw(http.HandlerFunc(h.HandleUpdate)))
}

// RegisterPortalRoutes registers portal-side feedback endpoints (portal cookie auth).
func (h *Handler) RegisterPortalRoutes(mux *http.ServeMux, authMw func(http.Handler) http.Handler) {
	mux.Handle("POST /api/portal/v1/feedback", authMw(http.HandlerFunc(h.HandleSubmitPortal)))
	mux.Handle("GET /api/portal/v1/feedback", authMw(http.HandlerFunc(h.HandleList)))
	mux.Handle("GET /api/portal/v1/feedback/{id}", authMw(http.HandlerFunc(h.HandleGet)))
	mux.Handle("PUT /api/portal/v1/feedback/{id}", authMw(http.HandlerFunc(h.HandleUpdate)))
}

// HandleSubmitERP handles feedback submission from the ERP UI.
func (h *Handler) HandleSubmitERP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	var req CreateFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondError(w, r, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	// Extract user info from ERP JWT claims.
	var userID *uuid.UUID
	if claims, ok := r.Context().Value(middleware.UserContextKey).(*middleware.UserClaims); ok && claims != nil {
		if req.Email == "" {
			req.Email = claims.Email
		}
	}

	created, err := h.svc.SubmitFeedback(r.Context(), req, "ERP", userID)
	if err != nil {
		httputil.RespondError(w, r, "Failed to submit feedback", http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, created)
}

// HandleSubmitPortal handles feedback submission from the Partner Portal.
func (h *Handler) HandleSubmitPortal(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	var req CreateFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondError(w, r, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	// Extract user info from portal JWT claims.
	var userID *uuid.UUID
	if claims, ok := r.Context().Value(middleware.PortalClaimsKey).(*middleware.PortalClaims); ok && claims != nil {
		uid := claims.CustomerUserID
		userID = &uid
		if req.Name == "" {
			req.Name = claims.Name
		}
		if req.Email == "" {
			req.Email = claims.Email
		}
	}

	created, err := h.svc.SubmitFeedback(r.Context(), req, "PORTAL", userID)
	if err != nil {
		httputil.RespondError(w, r, "Failed to submit feedback", http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondJSON(w, created)
}

// HandleList returns a paginated list of feedback items.
func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(q.Get("limit"))

	filter := FeedbackListFilter{
		Status:   q.Get("status"),
		Category: q.Get("category"),
		Source:   q.Get("source"),
		Search:   q.Get("search"),
		Page:     page,
		Limit:    limit,
	}

	result, err := h.svc.ListFeedback(r.Context(), filter)
	if err != nil {
		httputil.RespondError(w, r, "Failed to list feedback", http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, result)
}

// HandleGet returns a single feedback item by ID.
func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.RespondError(w, r, "Invalid feedback ID", http.StatusBadRequest, err)
		return
	}

	fb, err := h.svc.GetFeedback(r.Context(), id)
	if err != nil {
		httputil.RespondError(w, r, "Feedback not found", http.StatusNotFound, err)
		return
	}

	respondJSON(w, fb)
}

// HandleUpdate allows admins to update status, priority, or admin notes.
func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httputil.RespondError(w, r, "Invalid feedback ID", http.StatusBadRequest, err)
		return
	}

	var req UpdateFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondError(w, r, "Invalid request body", http.StatusBadRequest, err)
		return
	}

	fb, err := h.svc.UpdateFeedback(r.Context(), id, req)
	if err != nil {
		httputil.RespondError(w, r, "Failed to update feedback", http.StatusBadRequest, err)
		return
	}

	respondJSON(w, fb)
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
