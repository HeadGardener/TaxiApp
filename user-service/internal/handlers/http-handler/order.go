package http_handler

import (
	"encoding/json"
	"errors"
	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

const (
	moneyPlaceholder = 5.0
)

type sendOrderReq struct {
	TaxiType models.TaxiType `json:"taxi_type"`
	From     string          `json:"from"`
	To       string          `json:"to"`
}

type sendCommentReq struct {
	OrderID string `json:"order_id"`
	Comment string `json:"comment"`
}

func (h *Handler) sendOrder(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
		return
	}

	var orderReq sendOrderReq
	if err = json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while decoding order req", err)
		return
	}

	if err = orderReq.validate(); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while validating order req", err)
		return
	}

	var order = &models.Order{
		TaxiType: orderReq.TaxiType,
		From:     orderReq.From,
		To:       orderReq.To,
	}

	orderID, err := h.orderService.SendOrder(r.Context(), userID, order)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while sending order", err)
		return
	}

	transactionID, err := h.transactionService.Create(r.Context(), moneyPlaceholder)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while preparing transaction", err)
		return
	}

	// on front make a pair of order_id and transaction_id, so we can process transaction by order_id
	newResponse(w, http.StatusCreated, map[string]any{
		"order_id":       orderID,
		"transaction_id": transactionID,
	})
}

func (h *Handler) sendComment(w http.ResponseWriter, r *http.Request) {
	var commentReq sendCommentReq
	if err := json.NewDecoder(r.Body).Decode(&commentReq); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while validating comment req", err)
		return
	}

	if err := h.orderService.SendComment(r.Context(), commentReq.OrderID, commentReq.Comment); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while sending comment", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"status": "added",
	})
}

// awaitingOrder - method to check user orders
func (h *Handler) awaitingOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
		return
	}

	orderID := chi.URLParam(r, orderIDParam)

	order, err := h.orderService.Get(orderID, userID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while getting order", err)
	}

	newResponse(w, http.StatusOK, order)

	if order.Status != models.AcceptStatus && order.Status != models.ConsumedStatus {
		if err = h.orderService.Delete(orderID, userID); err != nil {
			log.Printf("[ERROR] failed while deleting order: %s", err.Error())
		}
	}
}

func (r *sendOrderReq) validate() error {
	if _, ok := models.TaxiTypes[r.TaxiType]; !ok {
		return errors.New("invalid taxi type: only economy, comfort and business are available")
	}

	return nil
}
