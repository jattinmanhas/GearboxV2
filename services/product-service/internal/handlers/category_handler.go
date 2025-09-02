package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/services"
)

type ICategoryHandler interface {
	CreateCategory(w http.ResponseWriter, r *http.Request)
	GetCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryBySlug(w http.ResponseWriter, r *http.Request)
	UpdateCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	ListCategories(w http.ResponseWriter, r *http.Request)
	GetCategoryHierarchy(w http.ResponseWriter, r *http.Request)
	GetCategoryChildren(w http.ResponseWriter, r *http.Request)
}

type categoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) ICategoryHandler {
	return &categoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory handles POST /api/v1/categories
func (h *categoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req services.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validateCreateCategoryRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := h.categoryService.CreateCategory(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// GetCategory handles GET /api/v1/categories/{id}
func (h *categoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	category, err := h.categoryService.GetCategory(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// GetCategoryBySlug handles GET /api/v1/categories/slug/{slug}
func (h *categoryHandler) GetCategoryBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	if slug == "" {
		http.Error(w, "Slug is required", http.StatusBadRequest)
		return
	}

	category, err := h.categoryService.GetCategoryBySlug(r.Context(), slug)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// UpdateCategory handles PUT /api/v1/categories/{id}
func (h *categoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var req services.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := validateUpdateCategoryRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := h.categoryService.UpdateCategory(r.Context(), id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// DeleteCategory handles DELETE /api/v1/categories/{id}
func (h *categoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	err = h.categoryService.DeleteCategory(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListCategories handles GET /api/v1/categories
func (h *categoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	req := &services.ListCategoriesRequest{}

	// Parse query parameters
	if parentIDStr := r.URL.Query().Get("parent_id"); parentIDStr != "" {
		if parentID, err := strconv.ParseInt(parentIDStr, 10, 64); err == nil {
			req.ParentID = &parentID
		}
	}

	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			req.IsActive = &isActive
		}
	}

	req.Search = r.URL.Query().Get("search")

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			req.Limit = limit
		}
	}

	response, err := h.categoryService.ListCategories(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCategoryHierarchy handles GET /api/v1/categories/hierarchy
func (h *categoryHandler) GetCategoryHierarchy(w http.ResponseWriter, r *http.Request) {
	hierarchy, err := h.categoryService.GetCategoryHierarchy(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hierarchy)
}

// GetCategoryChildren handles GET /api/v1/categories/{id}/children
func (h *categoryHandler) GetCategoryChildren(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	children, err := h.categoryService.GetCategoryChildren(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(children)
}

// Validation functions
func validateCreateCategoryRequest(req *services.CreateCategoryRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if len(req.Name) > 255 {
		return fmt.Errorf("name must be less than 255 characters")
	}
	if strings.TrimSpace(req.Slug) == "" {
		return fmt.Errorf("slug is required")
	}
	if len(req.Slug) > 255 {
		return fmt.Errorf("slug must be less than 255 characters")
	}
	return nil
}

func validateUpdateCategoryRequest(req *services.UpdateCategoryRequest) error {
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return fmt.Errorf("name cannot be empty")
		}
		if len(*req.Name) > 255 {
			return fmt.Errorf("name must be less than 255 characters")
		}
	}
	if req.Slug != nil {
		if strings.TrimSpace(*req.Slug) == "" {
			return fmt.Errorf("slug cannot be empty")
		}
		if len(*req.Slug) > 255 {
			return fmt.Errorf("slug must be less than 255 characters")
		}
	}
	return nil
}
