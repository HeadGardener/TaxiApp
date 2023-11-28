package http_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type credsReq struct {
	WalletID string  `json:"wallet_id"`
	Money    float64 `json:"spent"`
}

func (h *Handler) viewTransactions(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
		return
	}

	walletID := chi.URLParam(r, walletIDParam)
	if _, err = uuid.Parse(walletID); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while parsing wallet id", err)
		return
	}

	walletType := r.URL.Query().Get(walletTypeQuery)
	if _, ok := models.WalletTypesStr[walletType]; !ok {
		newErrResponse(w, http.StatusBadRequest, "failed while checking wallet type",
			models.ErrInvalidWalletType)
		return
	}

	transactions, err := h.transactionService.ViewAll(r.Context(), userID, walletID, models.WalletTypesStr[walletType])
	if err != nil {
		if errIsCustom(err) {
			newErrResponse(w, http.StatusBadRequest, "failed while getting transactions", err)
			return
		}

		newErrResponse(w, http.StatusInternalServerError, "failed while getting transactions", err)
		return
	}

	newResponse(w, http.StatusOK, transactions)
}

func (h *Handler) processTransaction(w http.ResponseWriter, r *http.Request) {
	transactionID, err := strconv.Atoi(chi.URLParam(r, transactionIDParam))
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while converting transaction id", err)
		return
	}

	walletType := r.URL.Query().Get(walletTypeQuery)
	if _, ok := models.WalletTypesStr[walletType]; !ok {
		newErrResponse(w, http.StatusBadRequest, "failed while checking wallet type",
			models.ErrInvalidWalletType)
		return
	}

	status := r.URL.Query().Get(transactionStatusQuery)

	switch status {
	case models.Success:
		if err = h.transactionService.Confirm(r.Context(),
			models.WalletTypesStr[walletType],
			transactionID); err != nil {
			if errIsCustom(err) {
				newErrResponse(w, http.StatusBadRequest, "failed while confirming transaction", err)
				return
			}

			newErrResponse(w, http.StatusInternalServerError, "failed while confirming transaction", err)
			return
		}

	case models.Canceled:
		var req credsReq

		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			newErrResponse(w, http.StatusBadRequest, "failed while decoding creds req", err)
			return
		}

		if err = req.validate(); err != nil {
			newErrResponse(w, http.StatusBadRequest, "failed while validating creds req", err)
			return
		}

		credentials := models.Credentials{
			WalletID: req.WalletID,
			Money:    req.Money,
		}

		if err = h.transactionService.Cancel(r.Context(),
			models.WalletTypesStr[walletType],
			transactionID, credentials); err != nil {
			if errIsCustom(err) {
				newErrResponse(w, http.StatusBadRequest, "failed while canceling transaction", err)
				return
			}

			newErrResponse(w, http.StatusInternalServerError, "failed while canceling transaction", err)
			return
		}

	default:
		newErrResponse(w, http.StatusBadRequest, "failed while processing transaction",
			errors.New("invalid transaction status"))
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"Msg": "processed",
	})
}

func (c *credsReq) validate() error {
	if _, err := uuid.Parse(c.WalletID); err != nil {
		return fmt.Errorf("invalid walletID: %e", err)
	}

	if c.Money <= 0.0 {
		return errors.New("invalid spent value, you can't spent negative amount of money")
	}

	return nil
}
