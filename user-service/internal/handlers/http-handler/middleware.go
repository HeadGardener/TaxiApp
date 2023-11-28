package http_handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/HeadGardener/TaxiApp/user-service/internal/lib/auth"
	"net/http"
	"strings"
)

type UserCtx string

const (
	userCtx UserCtx = "userAtr"
)

const (
	headerPartsLen = 2
)

var (
	ErrNotUserAttributes = errors.New("userCtx value is not of type UserAttributes")
)

func (h *Handler) identifyUser(next http.Handler) http.Handler {
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

		token := headerParts[1]
		if token == "" {
			newErrResponse(w, http.StatusUnauthorized, "failed while identifying user",
				errors.New("jwt token is empty"))
			return
		}

		userAttributes, err := auth.ParseToken(token)
		if err != nil {
			newErrResponse(w, http.StatusUnauthorized, "failed while parsing token", err)
			return
		}

		if err = h.identityService.Check(r.Context(), userAttributes.ID, token); err != nil {
			newErrResponse(w, http.StatusUnauthorized, "failed while checking session", err)
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, userAttributes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserID(r *http.Request) (string, error) {
	userCtxValue := r.Context().Value(userCtx)
	userAttributes, ok := userCtxValue.(auth.UserAttributes)
	if !ok {
		return "", ErrNotUserAttributes
	}

	return userAttributes.ID, nil
}
