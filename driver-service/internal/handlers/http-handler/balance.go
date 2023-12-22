package http_handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
)

var checkCard = regexp.MustCompile(`^\d{4}-\d{4}-\d{4}-\d{4}$`)

type credsReq struct {
	Card  string  `json:"card"`
	Money float64 `json:"money"`
}

func (h *Handler) addBalance(w http.ResponseWriter, r *http.Request) {
	driverID, err := GetDriverID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting driver id", err)
		return
	}

	money, err := strconv.ParseFloat(r.URL.Query().Get(moneyQuery), 64)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while parsing money query param", err)
		return
	}

	if err = h.balanceService.Add(r.Context(), driverID, money); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed adding balance", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"status": "added",
	})
}

func (h *Handler) cashOut(w http.ResponseWriter, r *http.Request) {
	driverID, err := GetDriverID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting driver id", err)
		return
	}

	var req credsReq

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while decoding sign in req", err)
		return
	}

	if err = req.validate(); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while validating sign in req", err)
		return
	}

	creds := models.Credentials{
		Card:  req.Card,
		Money: req.Money,
	}

	if err = h.balanceService.CashOut(r.Context(), driverID, creds); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while cashing out", err)
		return
	}
}

func (r *credsReq) validate() error {
	if !checkCard.MatchString(r.Card) {
		return errors.New("invalid card format: must be like ****-****-****-****")
	}

	if r.Money <= 0 {
		return errors.New("money must be greater than 0")
	}

	return nil
}
