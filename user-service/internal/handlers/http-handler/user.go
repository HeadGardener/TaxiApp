package http_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
	"net/http"
	"net/mail"
)

type updateReq struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

func (h *Handler) profile(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
		return
	}

	profile, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while getting profile", err)
		return
	}

	newResponse(w, http.StatusOK, profile)
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
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

	user := &models.User{
		Name:    req.Name,
		Surname: req.Surname,
		Phone:   req.Phone,
		Email:   req.Email,
	}

	if err = h.userService.Update(r.Context(), userID, user); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while updating profile", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"Msg": "updated",
	})
}

func (h *Handler) setInactive(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
		return
	}

	if err = h.userService.SetInactive(r.Context(), userID); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while setting inactive", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"Msg": "deleted",
	})
}

func (u *updateReq) validate() error {
	if (u.Name != "" || u.Surname != "") && (!checkName.MatchString(u.Name) || !checkName.MatchString(u.Surname)) {
		return errors.New("invalid name, it must contain only letters")
	}

	if u.Phone != "" && !checkPhone.MatchString(u.Phone) {
		return errors.New("invalid phone number, must be like +**[*] **[*] *******, where (*) is a number")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf("invalid email address: %e", err)
	}

	return nil
}
