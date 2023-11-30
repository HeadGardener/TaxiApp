package http_handler

import (
	"encoding/json"
	"errors"
	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var checkCard = regexp.MustCompile(`^\d{4}-\d{4}-\d{4}-\d{4}$`)

type createWalletReq struct {
	Card string `json:"card_number"`
}

// personal route

func (h *Handler) createWallet(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
		return
	}

	var req createWalletReq

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while decoding create wallet req", err)
		return
	}

	if err = req.validate(); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while validating create wallet req", err)
		return
	}

	id, err := h.walletService.Create(r.Context(), userID, req.Card)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while creating wallet", err)
		return
	}

	newResponse(w, http.StatusCreated, map[string]any{
		"id": id,
	})
}

func (h *Handler) viewWallet(w http.ResponseWriter, r *http.Request) {
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

	wallet, err := h.walletService.GetByID(r.Context(), userID, walletID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while getting wallet", err)
		return
	}

	newResponse(w, http.StatusOK, wallet)
}

func (h *Handler) viewWallets(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
		return
	}

	wallets, err := h.walletService.ViewAll(r.Context(), userID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while getting wallets", err)
		return
	}

	newResponse(w, http.StatusOK, wallets)
}

func (h *Handler) topUp(w http.ResponseWriter, r *http.Request) {
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

	money, err := strconv.ParseFloat(r.URL.Query().Get(moneyQuery), 64)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while parsing money", err)
		return
	}

	id, err := h.walletService.TopUp(r.Context(), userID, walletID, money)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while refilling balance", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"transaction_id": id,
	})
}

// family route

//nolint:dupl
func (h *Handler) createFamilyWallet(w http.ResponseWriter, r *http.Request) {
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

	id, err := h.walletService.CreateFamilyWallet(r.Context(), userID, walletID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while creating family wallet", err)
		return
	}

	newResponse(w, http.StatusCreated, map[string]any{
		"id": id,
	})
}

func (h *Handler) viewFamilyWallet(w http.ResponseWriter, r *http.Request) {
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

	wallet, err := h.walletService.GetFamilyWalletByID(r.Context(), userID, walletID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while getting family wallet", err)
		return
	}

	newResponse(w, http.StatusOK, wallet)
}

func (h *Handler) viewFamilyWallets(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
		return
	}

	wallets, err := h.walletService.ViewAllFamily(r.Context(), userID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while getting wallets", err)
		return
	}

	newResponse(w, http.StatusOK, wallets)
}

func (h *Handler) viewMemberships(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
		return
	}

	wallets, err := h.walletService.ViewMemberships(r.Context(), userID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while getting memberships", err)
		return
	}

	newResponse(w, http.StatusOK, wallets)
}

func (h *Handler) deleteFamilyWallet(w http.ResponseWriter, r *http.Request) {
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

	if err = h.walletService.DeleteFamilyWallet(r.Context(), userID, walletID); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while deleting family wallet", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"status": "deleted",
	})
}

func (h *Handler) topUpFamily(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
		return
	}

	famWalletID := chi.URLParam(r, walletIDParam)
	if _, err = uuid.Parse(famWalletID); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while parsing fam wallet id", err)
		return
	}

	walletID := r.URL.Query().Get(walletIDParam)
	if _, err = uuid.Parse(walletID); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while parsing wallet id", err)
		return
	}

	money, err := strconv.ParseFloat(r.URL.Query().Get(moneyQuery), 64)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while parsing money", err)
		return
	}

	if err = h.walletService.AddFamilyBalance(r.Context(), userID, walletID, famWalletID, money); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while refilling balance", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"status": "success",
	})
}

func (h *Handler) setFixedBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
		return
	}

	walletID := r.URL.Query().Get(walletIDParam)
	if _, err = uuid.Parse(walletID); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while parsing wallet id", err)
		return
	}

	money, err := strconv.ParseFloat(r.URL.Query().Get(moneyQuery), 64)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while parsing money", err)
		return
	}

	if err = h.walletService.SetFixedBalance(r.Context(), userID, walletID, money); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while setting fixed balance", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"status": "set",
	})
}

func (h *Handler) addMember(w http.ResponseWriter, r *http.Request) {
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

	phone := r.URL.Query().Get(phoneQuery)
	if !checkPhone.MatchString(phone) {
		newErrResponse(w, http.StatusBadRequest, "failed while matching phone", errors.New("invalid phone"))
		return
	}

	if err = h.walletService.AddMember(r.Context(), userID, walletID, phone); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while adding member", err)
		return
	}

	newResponse(w, http.StatusCreated, map[string]any{
		"status": "added",
	})
}

//nolint:dupl
func (h *Handler) viewMembers(w http.ResponseWriter, r *http.Request) {
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

	members, err := h.walletService.ViewMembers(r.Context(), userID, walletID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while getting members", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"members": members,
	})
}

func (h *Handler) deleteMember(w http.ResponseWriter, r *http.Request) {
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

	phone := r.URL.Query().Get(phoneQuery)
	if !checkPhone.MatchString(phone) {
		newErrResponse(w, http.StatusBadRequest, "failed while matching phone", errors.New("invalid phone"))
		return
	}

	if err = h.walletService.DeleteMember(r.Context(), userID, walletID, phone); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while deleting member", err)
		return
	}

	newResponse(w, http.StatusCreated, map[string]any{
		"status": "deleted",
	})
}

// pay route

func (h *Handler) pay(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
		return
	}

	walletType := r.URL.Query().Get(walletTypeQuery)
	if _, ok := models.WalletTypesStr[walletType]; !ok {
		newErrResponse(w, http.StatusBadRequest, "failed while checking wallet type",
			models.ErrInvalidWalletType)
		return
	}

	transactionID, err := strconv.Atoi(chi.URLParam(r, transactionIDParam))
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while converting transaction id", err)
		return
	}

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

	if err = h.walletService.PickWalletAndPay(r.Context(),
		userID,
		models.WalletTypesStr[walletType],
		transactionID,
		credentials); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while paying", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"status": "success",
	})
}

func (w *createWalletReq) validate() error {
	if !checkCard.MatchString(w.Card) {
		return errors.New("invalid card number, must be like ****-****-****-****, where (*) is a number")
	}

	return nil
}
