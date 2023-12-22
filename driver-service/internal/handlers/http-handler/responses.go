package http_handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/HeadGardener/TaxiApp/driver-service/internal/services"
)

type response struct {
	Msg   string `json:"Msg"`
	Error string `json:"Error"`
}

func newErrResponse(w http.ResponseWriter, code int, msg string, err error) {
	log.Printf("[ERROR] %s: %s", msg, err.Error())
	if !errIsCustom(err) && code >= http.StatusInternalServerError {
		newResponse(w, code, response{
			Msg:   msg,
			Error: "unexpected error",
		})
		return
	}

	newResponse(w, code, response{
		Msg:   msg,
		Error: err.Error(),
	})
}

func newResponse(w http.ResponseWriter, code int, data any) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

func errIsCustom(err error) bool {
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}

	if errors.Is(err, services.ErrNotEnoughBalance) {
		return true
	}

	if errors.Is(err, services.ErrNotActive) {
		return true
	}

	if errors.Is(err, services.ErrInvalidPassword) {
		return true
	}

	if errors.Is(err, ErrNotUserAttributes) {
		return true
	}

	return false
}
