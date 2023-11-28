package http_handler

import "net/http"

func (h *Handler) viewTrips(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting user id", err)
		return
	}

	trips, err := h.tripService.ViewAll(r.Context(), userID)
	if err != nil {
		if errIsCustom(err) {
			newErrResponse(w, http.StatusBadRequest, "failed while getting trips", err)
			return
		}

		newErrResponse(w, http.StatusInternalServerError, "failed while getting trips", err)
		return
	}

	newResponse(w, http.StatusOK, trips)
}
