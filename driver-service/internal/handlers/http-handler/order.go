package http_handler

import (
	"encoding/json"
	"net/http"
)

type processOrderReq struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type completeOrderReq struct {
	OrderID string  `json:"order_id"`
	Status  string  `json:"status"`
	Rating  float64 `json:"rating"`
}

func (h *Handler) getInLine(w http.ResponseWriter, r *http.Request) {
	driverID, err := GetDriverID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting driver id", err)
		return
	}

	taxiType, err := GetTaxiType(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting driver taxi type", err)
		return
	}

	if err = h.orderService.GetInLine(r.Context(), driverID, taxiType); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while getting in line", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"status": "in line",
	})
}

func (h *Handler) processOrder(w http.ResponseWriter, r *http.Request) {
	driverID, err := GetDriverID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting driver id", err)
		return
	}

	var req processOrderReq

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while decoding process req", err)
		return
	}

	if err = req.validate(); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while validating process req", err)
		return
	}

	if err = h.orderService.ProcessOrder(r.Context(), driverID, req.OrderID, req.Status); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while processing order", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"status": "processed",
	})
}

func (h *Handler) completeOrder(w http.ResponseWriter, r *http.Request) {
	driverID, err := GetDriverID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting driver id", err)
		return
	}

	var req completeOrderReq

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while decoding complete order req", err)
		return
	}

	if err = req.validate(); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while validating complete order req", err)
		return
	}

	if err = h.orderService.Complete(r.Context(), driverID, req.OrderID, req.Status, req.Rating); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while completing order", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"status": "completed",
	})
}

func (h *Handler) awaitingOrder(w http.ResponseWriter, r *http.Request) {
	driverID, err := GetDriverID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting driver id", err)
		return
	}

	order, err := h.orderService.CurrentOrder(driverID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while getting current order", err)
		return
	}

	newResponse(w, http.StatusOK, order)
}

func (r *processOrderReq) validate() error {
	return nil
}

func (r *completeOrderReq) validate() error {
	return nil
}
