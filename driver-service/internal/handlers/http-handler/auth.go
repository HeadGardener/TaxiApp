package http_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"

	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
)

var (
	checkName     = regexp.MustCompile(`[A-z]$`)
	checkPhone    = regexp.MustCompile(`^[+]?\d{2,3} \d{2,3} \d{7}$`)
	checkPassword = regexp.MustCompile(`[0-9A-z]{8,16}$`)
)

type signUpReq struct {
	Name     string          `json:"name"`
	Surname  string          `json:"surname"`
	Phone    string          `json:"phone"`
	Email    string          `json:"email"`
	TaxiType models.TaxiType `json:"taxi_type"`
	Password string          `json:"password"`
}

type signInReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var req signUpReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while decoding sign up req", err)
		return
	}

	if err := req.validate(); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while validating sign up req", err)
		return
	}

	driver := &models.Driver{
		Name:     req.Name,
		Surname:  req.Surname,
		Phone:    req.Phone,
		Email:    req.Email,
		TaxiType: req.TaxiType,
		Password: req.Password,
	}

	id, err := h.identityService.SignUp(r.Context(), driver)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while signing up", err)
		return
	}

	newResponse(w, http.StatusCreated, map[string]any{
		"id": id,
	})
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var req signInReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while decoding sign in req", err)
		return
	}

	if err := req.validate(); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while validating sign in req", err)
		return
	}

	token, err := h.identityService.SignIn(r.Context(), req.Phone, req.Password)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while signing in", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"token": token,
	})
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	driverID, err := GetDriverID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting driver id", err)
		return
	}

	if err = h.identityService.LogOut(r.Context(), driverID); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while logging out", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]any{
		"msg": "logged out",
	})
}

func (r *signUpReq) validate() error {
	if !checkName.MatchString(r.Name) || !checkName.MatchString(r.Surname) {
		return errors.New("invalid name, it must contain only letters")
	}

	if !checkPhone.MatchString(r.Phone) {
		return errors.New("invalid phone number, should be like +**[*] **[*] *******, where (*) is a number")
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return fmt.Errorf("invalid email address: %e", err)
	}

	if _, ok := models.TaxiTypes[r.TaxiType]; !ok {
		return models.ErrInvalidTaxiType
	}

	if !checkPassword.MatchString(r.Password) {
		return errors.New("invalid password, it must contain only letter and number and 8-16 symbols")
	}

	return nil
}

func (r *signInReq) validate() error {
	if !checkPhone.MatchString(r.Phone) {
		return errors.New("invalid phone number, should be like +**[*] **[*] *******, where (*) is a number")
	}

	if !checkPassword.MatchString(r.Password) {
		return errors.New("invalid password, it must contain only letter and number and 8-16 symbols")
	}

	return nil
}
