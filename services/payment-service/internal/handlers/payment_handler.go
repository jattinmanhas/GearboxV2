package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/domain"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/dto"
	"github.com/jattinmanhas/GearboxV2/services/payment-service/internal/services"
	"github.com/jattinmanhas/GearboxV2/services/shared/httpx"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	paymentService *services.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// Payment Methods

// CreatePaymentMethod creates a new payment method
func (h *PaymentHandler) CreatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	var req dto.PaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	response, err := h.paymentService.CreatePaymentMethod(r.Context(), &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to create payment method", err)
		return
	}

	httpx.Created(w, "Payment method created successfully", response)
}

// GetPaymentMethod retrieves a payment method by ID
func (h *PaymentHandler) GetPaymentMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid payment method ID", err)
		return
	}

	response, err := h.paymentService.GetPaymentMethod(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get payment method", err)
		return
	}

	httpx.OK(w, "Payment method retrieved successfully", response)
}

// ListPaymentMethods retrieves payment methods with filtering
func (h *PaymentHandler) ListPaymentMethods(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filter := &domain.PaymentMethodFilter{}

	if typeParam := r.URL.Query().Get("type"); typeParam != "" {
		filter.Type = &typeParam
	}

	if activeParam := r.URL.Query().Get("is_active"); activeParam != "" {
		active, err := strconv.ParseBool(activeParam)
		if err == nil {
			filter.IsActive = &active
		}
	}

	filter.Search = r.URL.Query().Get("search")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	response, err := h.paymentService.ListPaymentMethods(r.Context(), filter, page, limit)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to list payment methods", err)
		return
	}

	httpx.OK(w, "Payment methods retrieved successfully", response)
}

// UpdatePaymentMethod updates a payment method
func (h *PaymentHandler) UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid payment method ID", err)
		return
	}

	var req dto.PaymentMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	response, err := h.paymentService.UpdatePaymentMethod(r.Context(), id, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to update payment method", err)
		return
	}

	httpx.OK(w, "Payment method updated successfully", response)
}

// DeletePaymentMethod deletes a payment method
func (h *PaymentHandler) DeletePaymentMethod(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid payment method ID", err)
		return
	}

	if err := h.paymentService.DeletePaymentMethod(r.Context(), id); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to delete payment method", err)
		return
	}

	httpx.OK(w, "Payment method deleted successfully", nil)
}

// Payment Gateways

// CreatePaymentGateway creates a new payment gateway
func (h *PaymentHandler) CreatePaymentGateway(w http.ResponseWriter, r *http.Request) {
	var req dto.PaymentGatewayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	response, err := h.paymentService.CreatePaymentGateway(r.Context(), &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to create payment gateway", err)
		return
	}

	httpx.Created(w, "Payment gateway created successfully", response)
}

// GetPaymentGateway retrieves a payment gateway by ID
func (h *PaymentHandler) GetPaymentGateway(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid payment gateway ID", err)
		return
	}

	response, err := h.paymentService.GetPaymentGateway(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get payment gateway", err)
		return
	}

	httpx.OK(w, "Payment gateway retrieved successfully", response)
}

// ListPaymentGateways retrieves payment gateways with filtering
func (h *PaymentHandler) ListPaymentGateways(w http.ResponseWriter, r *http.Request) {
	filter := &domain.PaymentGatewayFilter{}

	if codeParam := r.URL.Query().Get("code"); codeParam != "" {
		filter.Code = &codeParam
	}

	if activeParam := r.URL.Query().Get("is_active"); activeParam != "" {
		active, err := strconv.ParseBool(activeParam)
		if err == nil {
			filter.IsActive = &active
		}
	}

	if testModeParam := r.URL.Query().Get("is_test_mode"); testModeParam != "" {
		testMode, err := strconv.ParseBool(testModeParam)
		if err == nil {
			filter.IsTestMode = &testMode
		}
	}

	filter.Search = r.URL.Query().Get("search")

	response, err := h.paymentService.ListPaymentGateways(r.Context(), filter)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to list payment gateways", err)
		return
	}

	httpx.OK(w, "Payment gateways retrieved successfully", response)
}

// UpdatePaymentGateway updates a payment gateway
func (h *PaymentHandler) UpdatePaymentGateway(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid payment gateway ID", err)
		return
	}

	var req dto.PaymentGatewayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	response, err := h.paymentService.UpdatePaymentGateway(r.Context(), id, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to update payment gateway", err)
		return
	}

	httpx.OK(w, "Payment gateway updated successfully", response)
}

// DeletePaymentGateway deletes a payment gateway
func (h *PaymentHandler) DeletePaymentGateway(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid payment gateway ID", err)
		return
	}

	if err := h.paymentService.DeletePaymentGateway(r.Context(), id); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to delete payment gateway", err)
		return
	}

	httpx.OK(w, "Payment gateway deleted successfully", nil)
}

// Payments

// CreatePayment creates a new payment
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	response, err := h.paymentService.CreatePayment(r.Context(), &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to create payment", err)
		return
	}

	httpx.Created(w, "Payment created successfully", response)
}

// ProcessPayment processes a payment through the gateway
func (h *PaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	var req dto.ProcessPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	response, err := h.paymentService.ProcessPayment(r.Context(), &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to process payment", err)
		return
	}

	httpx.OK(w, "Payment processed successfully", response)
}

// GetPayment retrieves a payment by ID
func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid payment ID", err)
		return
	}

	response, err := h.paymentService.GetPayment(r.Context(), id)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get payment", err)
		return
	}

	httpx.OK(w, "Payment retrieved successfully", response)
}

// GetPaymentByTransactionID retrieves a payment by transaction ID
func (h *PaymentHandler) GetPaymentByTransactionID(w http.ResponseWriter, r *http.Request) {
	transactionID := chi.URLParam(r, "transactionId")
	if transactionID == "" {
		httpx.Error(w, http.StatusBadRequest, "Transaction ID is required", nil)
		return
	}

	response, err := h.paymentService.GetPaymentByTransactionID(r.Context(), transactionID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get payment", err)
		return
	}

	httpx.OK(w, "Payment retrieved successfully", response)
}

// ListPayments retrieves payments with filtering and pagination
func (h *PaymentHandler) ListPayments(w http.ResponseWriter, r *http.Request) {
	filter := &domain.PaymentFilter{}

	if orderIDParam := r.URL.Query().Get("order_id"); orderIDParam != "" {
		orderID, err := strconv.ParseInt(orderIDParam, 10, 64)
		if err == nil {
			filter.OrderID = &orderID
		}
	}

	if paymentMethodIDParam := r.URL.Query().Get("payment_method_id"); paymentMethodIDParam != "" {
		paymentMethodID, err := strconv.ParseInt(paymentMethodIDParam, 10, 64)
		if err == nil {
			filter.PaymentMethodID = &paymentMethodID
		}
	}

	if statusParam := r.URL.Query().Get("status"); statusParam != "" {
		filter.Status = &statusParam
	}

	if gatewayIDParam := r.URL.Query().Get("gateway_id"); gatewayIDParam != "" {
		filter.GatewayID = &gatewayIDParam
	}

	if currencyParam := r.URL.Query().Get("currency"); currencyParam != "" {
		filter.Currency = &currencyParam
	}

	if dateFromParam := r.URL.Query().Get("date_from"); dateFromParam != "" {
		dateFrom, err := time.Parse("2006-01-02", dateFromParam)
		if err == nil {
			filter.DateFrom = &dateFrom
		}
	}

	if dateToParam := r.URL.Query().Get("date_to"); dateToParam != "" {
		dateTo, err := time.Parse("2006-01-02", dateToParam)
		if err == nil {
			filter.DateTo = &dateTo
		}
	}

	filter.Search = r.URL.Query().Get("search")
	filter.SortBy = r.URL.Query().Get("sort_by")
	filter.SortOrder = r.URL.Query().Get("sort_order")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	response, err := h.paymentService.ListPayments(r.Context(), filter, page, limit)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to list payments", err)
		return
	}

	httpx.OK(w, "Payments retrieved successfully", response)
}

// UpdatePaymentStatus updates payment status
func (h *PaymentHandler) UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid payment ID", err)
		return
	}

	var req dto.UpdatePaymentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	response, err := h.paymentService.UpdatePaymentStatus(r.Context(), id, &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to update payment status", err)
		return
	}

	httpx.OK(w, "Payment status updated successfully", response)
}

// RefundPayment refunds a payment
func (h *PaymentHandler) RefundPayment(w http.ResponseWriter, r *http.Request) {
	var req dto.RefundPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	response, err := h.paymentService.RefundPayment(r.Context(), &req)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to refund payment", err)
		return
	}

	httpx.OK(w, "Payment refunded successfully", response)
}

// GetPaymentSummary retrieves payment summary statistics
func (h *PaymentHandler) GetPaymentSummary(w http.ResponseWriter, r *http.Request) {
	var dateFrom, dateTo *time.Time

	if dateFromParam := r.URL.Query().Get("date_from"); dateFromParam != "" {
		date, err := time.Parse("2006-01-02", dateFromParam)
		if err == nil {
			dateFrom = &date
		}
	}

	if dateToParam := r.URL.Query().Get("date_to"); dateToParam != "" {
		date, err := time.Parse("2006-01-02", dateToParam)
		if err == nil {
			dateTo = &date
		}
	}

	response, err := h.paymentService.GetPaymentSummary(r.Context(), dateFrom, dateTo)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "Failed to get payment summary", err)
		return
	}

	httpx.OK(w, "Payment summary retrieved successfully", response)
}

// Webhooks

// HandleWebhook handles webhook events from payment gateways
func (h *PaymentHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	gatewayID := chi.URLParam(r, "gatewayId")
	if gatewayID == "" {
		httpx.Error(w, http.StatusBadRequest, "Gateway ID is required", nil)
		return
	}

	var req dto.WebhookEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// TODO: Implement webhook processing
	// This would validate the webhook signature and process the event

	response := dto.WebhookEventResponse{
		Success: true,
		Message: "Webhook processed successfully",
	}

	httpx.OK(w, "Webhook processed successfully", response)
}
