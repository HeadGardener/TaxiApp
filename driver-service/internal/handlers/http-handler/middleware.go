package http_handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/lib/auth"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
	"net/http"
	"strings"
)

type DriverCtx string

const (
	driverCtx DriverCtx = "userAtr"
)

const (
	headerPartsLen = 2
)

var (
	ErrNotUserAttributes = errors.New("driverCtx value is not of type DriverAttributes")
)

func (h *Handler) identifyDriver(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" {
			newErrResponse(w, http.StatusUnauthorized, "failed while identifying user",
				errors.New("empty auth header"))
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != headerPartsLen {
			newErrResponse(w, http.StatusUnauthorized, "failed while identifying user",
				errors.New("invalid auth header, must be like `Bearer token`"))
			return
		}

		if headerParts[0] != "Bearer" {
			newErrResponse(w, http.StatusUnauthorized, "failed while identifying user",
				fmt.Errorf("invalid auth header %s, must be Bearer", headerParts[0]))
			return
		}

		if headerParts[1] == "" {
			newErrResponse(w, http.StatusUnauthorized, "failed while identifying user",
				errors.New("jwt token is empty"))
			return
		}

		token := headerParts[1]
		driverAttributes, err := auth.ParseToken(token)
		if err != nil {
			newErrResponse(w, http.StatusUnauthorized, "failed while parsing token", err)
			return
		}

		if err = h.identityService.Check(r.Context(), driverAttributes.ID, token); err != nil {
			newErrResponse(w, http.StatusUnauthorized, "failed while checking session", err)
			return
		}

		ctx := context.WithValue(r.Context(), driverCtx, driverAttributes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetDriverID(r *http.Request) (string, error) {
	driverCtxValue := r.Context().Value(driverCtx)
	driverAttributes, ok := driverCtxValue.(auth.DriverAttributes)
	if !ok {
		return "", ErrNotUserAttributes
	}

	return driverAttributes.ID, nil
}

func GetTaxiType(r *http.Request) (models.TaxiType, error) {
	driverCtxValue := r.Context().Value(driverCtx)
	driverAttributes, ok := driverCtxValue.(auth.DriverAttributes)
	if !ok {
		return -1, ErrNotUserAttributes
	}

	return driverAttributes.TaxiType, nil
}
