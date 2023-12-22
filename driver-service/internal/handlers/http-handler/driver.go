package http_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
)

type updateReq struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

func (h *Handler) profile(w http.ResponseWriter, r *http.Request) {
	driverID, err := GetDriverID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting driver id", err)
		return
	}

	driver, err := h.driverService.GetProfile(r.Context(), driverID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while getting profile", err)
		return
	}

	newResponse(w, http.StatusOK, driver)
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	driverID, err := GetDriverID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting driver id", err)
		return
	}

	var req updateReq

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while decoding update req", err)
		return
	}

	if err = req.validate(); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while validating update req", err)
		return
	}

	driverUpdate := &models.Driver{
		Name:    req.Name,
		Surname: req.Surname,
		Email:   req.Email,
		Phone:   req.Phone,
	}

	if err = h.driverService.Update(r.Context(), driverID, driverUpdate); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while updating driver", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"status": "updated",
	})
}

func (h *Handler) changeStatus(w http.ResponseWriter, r *http.Request) {
	driverID, err := GetDriverID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting driver id", err)
		return
	}

	status := r.URL.Query().Get(driverStatusQuery)

	if err = h.driverService.ChangeStatus(r.Context(), driverID, status); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while updating status", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"status": "changed",
	})
}

func (h *Handler) setInactive(w http.ResponseWriter, r *http.Request) {
	driverID, err := GetDriverID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting driver id", err)
		return
	}

	if err = h.driverService.SetInactive(r.Context(), driverID); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while setting inactive", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"status": "set inactive",
	})
}

func (r *updateReq) validate() error {
	if !checkName.MatchString(r.Name) || !checkName.MatchString(r.Surname) {
		return errors.New("invalid name, it must contain only letters")
	}

	if !checkPhone.MatchString(r.Phone) {
		return errors.New("invalid phone number, should be like +**[*] **[*] *******, where (*) is a number")
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return fmt.Errorf("invalid email address: %e", err)
	}

	return nil
}
