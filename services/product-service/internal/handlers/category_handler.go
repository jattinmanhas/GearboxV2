package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/validation"
	"github.com/jattinmanhas/GearboxV2/services/shared/httpx"
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
	var req dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Validate Request
	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	cat := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
		Slug:        req.Slug,
		ParentID:    req.ParentID,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
		ImageURL:    req.ImageURL,
		MetaTitle:   req.MetaTitle,
		MetaDesc:    req.MetaDescription,
		CreatedAt:   time.Now(),
	}

	category, err := h.categoryService.CreateCategory(r.Context(), cat)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to create category", err)
		return
	}

	httpx.Created(w, "category created", category)
}

// GetCategory handles GET /api/v1/categories/{id}
func (h *categoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid category ID", err)
		return
	}

	category, err := h.categoryService.GetCategory(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httpx.OK(w, "category retrieved", category)
}

// GetCategoryBySlug handles GET /api/v1/categories/slug/{slug}
func (h *categoryHandler) GetCategoryBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	if err := validation.ValidateSlug(slug); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid slug", err)
		return
	}

	category, err := h.categoryService.GetCategoryBySlug(r.Context(), slug)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	httpx.OK(w, "category retrieved", category)
}

// UpdateCategory handles PUT /api/v1/categories/{id}
func (h *categoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid category ID", err)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	category, err := h.categoryService.UpdateCategory(r.Context(), id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httpx.OK(w, "category updated", category)
}

// DeleteCategory handles DELETE /api/v1/categories/{id}
func (h *categoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid category ID", err)
		return
	}

	err = h.categoryService.DeleteCategory(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httpx.OK(w, "category deleted", nil)
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
		httpx.Error(w, http.StatusInternalServerError, "failed to list categories", err)
		return
	}

	httpx.OK(w, "categories retrieved", response)
}

// GetCategoryHierarchy handles GET /api/v1/categories/hierarchy
func (h *categoryHandler) GetCategoryHierarchy(w http.ResponseWriter, r *http.Request) {
	hierarchy, err := h.categoryService.GetCategoryHierarchy(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get category hierarchy", err)
		return
	}

	httpx.OK(w, "category hierarchy retrieved", hierarchy)
}

// GetCategoryChildren handles GET /api/v1/categories/{id}/children
func (h *categoryHandler) GetCategoryChildren(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid category ID", err)
		return
	}

	children, err := h.categoryService.GetCategoryChildren(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get category children", err)
		return
	}

	httpx.OK(w, "category children retrieved", children)
}
