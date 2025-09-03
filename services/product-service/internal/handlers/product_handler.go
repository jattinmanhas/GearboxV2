package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/product-service/internal/validation"
	"github.com/jattinmanhas/GearboxV2/services/shared/httpx"
)

type IProductHandler interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
	GetProductBySKU(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
	ListProducts(w http.ResponseWriter, r *http.Request)
	GetProductsByCategory(w http.ResponseWriter, r *http.Request)
	SearchProducts(w http.ResponseWriter, r *http.Request)
	UpdateProductQuantity(w http.ResponseWriter, r *http.Request)
	GetProductsByTags(w http.ResponseWriter, r *http.Request)

	// Product Variants
	CreateProductVariant(w http.ResponseWriter, r *http.Request)
	GetProductVariant(w http.ResponseWriter, r *http.Request)
	GetProductVariants(w http.ResponseWriter, r *http.Request)
	UpdateProductVariant(w http.ResponseWriter, r *http.Request)
	DeleteProductVariant(w http.ResponseWriter, r *http.Request)

	// Product Categories
	AddProductToCategory(w http.ResponseWriter, r *http.Request)
	RemoveProductFromCategory(w http.ResponseWriter, r *http.Request)
	GetProductCategories(w http.ResponseWriter, r *http.Request)
	UpdateProductCategories(w http.ResponseWriter, r *http.Request)
}

type productHandler struct {
	productService services.ProductService
}

func NewProductHandler(productService services.ProductService) IProductHandler {
	return &productHandler{
		productService: productService,
	}
}

// CreateProduct handles POST /api/v1/products
func (h *productHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	// Validate Request
	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	product, err := h.productService.CreateProduct(r.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			httpx.Error(w, http.StatusConflict, err.Error(), err)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "failed to create product", err)
		return
	}

	httpx.Created(w, "product created", product)
}

// GetProduct handles GET /api/v1/products/{id}
func (h *productHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid product ID", err)
		return
	}

	product, err := h.productService.GetProductByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httpx.OK(w, "product retrieved", product)
}

// GetProductBySKU handles GET /api/v1/products/sku/{sku}
func (h *productHandler) GetProductBySKU(w http.ResponseWriter, r *http.Request) {
	sku := chi.URLParam(r, "sku")

	if sku == "" {
		httpx.Error(w, http.StatusBadRequest, "SKU is required", nil)
		return
	}

	product, err := h.productService.GetProductBySKU(r.Context(), sku)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	httpx.OK(w, "product retrieved", product)
}

// UpdateProduct handles PUT /api/v1/products/{id}
func (h *productHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid product ID", err)
		return
	}

	var req dto.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	product, err := h.productService.UpdateProduct(r.Context(), id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			httpx.Error(w, http.StatusConflict, err.Error(), err)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httpx.OK(w, "product updated", product)
}

// DeleteProduct handles DELETE /api/v1/products/{id}
func (h *productHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid product ID", err)
		return
	}

	err = h.productService.DeleteProduct(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httpx.OK(w, "product deleted", nil)
}

// ListProducts handles GET /api/v1/products
func (h *productHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	req := &dto.ListProductsRequest{}

	// Parse query parameters
	if categoryIDStr := r.URL.Query().Get("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			req.CategoryID = &categoryID
		}
	}

	if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err == nil {
			req.IsActive = &isActive
		}
	}

	if isDigitalStr := r.URL.Query().Get("is_digital"); isDigitalStr != "" {
		if isDigital, err := strconv.ParseBool(isDigitalStr); err == nil {
			req.IsDigital = &isDigital
		}
	}

	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			req.MinPrice = &minPrice
		}
	}

	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			req.MaxPrice = &maxPrice
		}
	}

	if inStockStr := r.URL.Query().Get("in_stock"); inStockStr != "" {
		if inStock, err := strconv.ParseBool(inStockStr); err == nil {
			req.InStock = &inStock
		}
	}

	req.Search = r.URL.Query().Get("search")
	req.SortBy = r.URL.Query().Get("sort_by")
	req.SortOrder = r.URL.Query().Get("sort_order")

	// Parse tags (comma-separated)
	if tagsStr := r.URL.Query().Get("tags"); tagsStr != "" {
		req.Tags = strings.Split(tagsStr, ",")
		for i, tag := range req.Tags {
			req.Tags[i] = strings.TrimSpace(tag)
		}
	}

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

	response, err := h.productService.ListProducts(r.Context(), req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to list products", err)
		return
	}

	httpx.OK(w, "products retrieved", response)
}

// GetProductsByCategory handles GET /api/v1/categories/{id}/products
func (h *productHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := chi.URLParam(r, "id")

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid category ID", err)
		return
	}

	// Parse pagination parameters
	page := 1
	limit := 20

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	response, err := h.productService.GetProductsByCategory(r.Context(), categoryID, page, limit)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get products by category", err)
		return
	}

	httpx.OK(w, "products retrieved", response)
}

// SearchProducts handles GET /api/v1/products/search
func (h *productHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		httpx.Error(w, http.StatusBadRequest, "search query is required", nil)
		return
	}

	// Parse pagination parameters
	page := 1
	limit := 20

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	response, err := h.productService.SearchProducts(r.Context(), query, page, limit)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to search products", err)
		return
	}

	httpx.OK(w, "search results retrieved", response)
}

// UpdateProductQuantity handles PATCH /api/v1/products/{id}/quantity
func (h *productHandler) UpdateProductQuantity(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid product ID", err)
		return
	}

	var req struct {
		Quantity int `json:"quantity" validate:"required,min=0"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	err = h.productService.UpdateProductQuantity(r.Context(), id, req.Quantity)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httpx.OK(w, "product quantity updated", map[string]interface{}{
		"product_id": id,
		"quantity":   req.Quantity,
	})
}

// GetProductsByTags handles GET /api/v1/products/tags
func (h *productHandler) GetProductsByTags(w http.ResponseWriter, r *http.Request) {
	tagsStr := r.URL.Query().Get("tags")
	if tagsStr == "" {
		httpx.Error(w, http.StatusBadRequest, "tags parameter is required", nil)
		return
	}

	tags := strings.Split(tagsStr, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	// Parse pagination parameters
	page := 1
	limit := 20

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	response, err := h.productService.GetProductsByTags(r.Context(), tags, page, limit)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get products by tags", err)
		return
	}

	httpx.OK(w, "products retrieved", response)
}

// Product Variant handlers

// CreateProductVariant handles POST /api/v1/products/{id}/variants
func (h *productHandler) CreateProductVariant(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")

	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid product ID", err)
		return
	}

	var req dto.CreateProductVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	req.ProductID = productID

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	variant, err := h.productService.CreateProductVariant(r.Context(), &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to create product variant", err)
		return
	}

	httpx.Created(w, "product variant created", variant)
}

// GetProductVariant handles GET /api/v1/products/variants/{id}
func (h *productHandler) GetProductVariant(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid variant ID", err)
		return
	}

	variant, err := h.productService.GetProductVariantByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httpx.OK(w, "product variant retrieved", variant)
}

// GetProductVariants handles GET /api/v1/products/{id}/variants
func (h *productHandler) GetProductVariants(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")

	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid product ID", err)
		return
	}

	variants, err := h.productService.GetProductVariantsByProductID(r.Context(), productID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get product variants", err)
		return
	}

	httpx.OK(w, "product variants retrieved", variants)
}

// UpdateProductVariant handles PUT /api/v1/products/variants/{id}
func (h *productHandler) UpdateProductVariant(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid variant ID", err)
		return
	}

	var req dto.UpdateProductVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	variant, err := h.productService.UpdateProductVariant(r.Context(), id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httpx.OK(w, "product variant updated", variant)
}

// DeleteProductVariant handles DELETE /api/v1/products/variants/{id}
func (h *productHandler) DeleteProductVariant(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid variant ID", err)
		return
	}

	err = h.productService.DeleteProductVariant(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			httpx.Error(w, http.StatusNotFound, err.Error(), nil)
			return
		}
		httpx.Error(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	httpx.OK(w, "product variant deleted", nil)
}

// Product Category handlers

// AddProductToCategory handles POST /api/v1/products/{id}/categories
func (h *productHandler) AddProductToCategory(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")

	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid product ID", err)
		return
	}

	var req struct {
		CategoryID int64 `json:"category_id" validate:"required"`
		IsPrimary  bool  `json:"is_primary"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	err = h.productService.AddProductToCategory(r.Context(), productID, req.CategoryID, req.IsPrimary)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to add product to category", err)
		return
	}

	httpx.OK(w, "product added to category", map[string]interface{}{
		"product_id":  productID,
		"category_id": req.CategoryID,
		"is_primary":  req.IsPrimary,
	})
}

// RemoveProductFromCategory handles DELETE /api/v1/products/{id}/categories/{category_id}
func (h *productHandler) RemoveProductFromCategory(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")
	categoryIDStr := chi.URLParam(r, "category_id")

	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid product ID", err)
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid category ID", err)
		return
	}

	err = h.productService.RemoveProductFromCategory(r.Context(), productID, categoryID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to remove product from category", err)
		return
	}

	httpx.OK(w, "product removed from category", map[string]interface{}{
		"product_id":  productID,
		"category_id": categoryID,
	})
}

// GetProductCategories handles GET /api/v1/products/{id}/categories
func (h *productHandler) GetProductCategories(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")

	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid product ID", err)
		return
	}

	categories, err := h.productService.GetProductCategories(r.Context(), productID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to get product categories", err)
		return
	}

	httpx.OK(w, "product categories retrieved", categories)
}

// UpdateProductCategories handles PUT /api/v1/products/{id}/categories
func (h *productHandler) UpdateProductCategories(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")

	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid product ID", err)
		return
	}

	var req struct {
		CategoryIDs []int64 `json:"category_ids" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if validationErrors := validation.ValidateStruct(req); len(validationErrors) > 0 {
		httpx.Error(w, http.StatusBadRequest, validationErrors.Error(), validationErrors)
		return
	}

	err = h.productService.UpdateProductCategories(r.Context(), productID, req.CategoryIDs)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "failed to update product categories", err)
		return
	}

	httpx.OK(w, "product categories updated", map[string]interface{}{
		"product_id":   productID,
		"category_ids": req.CategoryIDs,
	})
}
