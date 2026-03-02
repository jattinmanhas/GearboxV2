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

// GetPaymentByOrderID retrieves the latest payment by order ID
func (h *PaymentHandler) GetPaymentByOrderID(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "orderId")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "Invalid order ID", err)
		return
	}

	response, err := h.paymentService.GetPaymentByOrderID(r.Context(), orderID)
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

	if paymentMethodParam := r.URL.Query().Get("payment_method"); paymentMethodParam != "" {
		filter.PaymentMethod = &paymentMethodParam
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
