package http_handler

import "net/http"

func (h *Handler) viewTrips(w http.ResponseWriter, r *http.Request) {
	driverID, err := GetDriverID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting driver id", err)
		return
	}

	trips, err := h.tripService.ViewAll(r.Context(), driverID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while getting trips", err)
		return
	}

	newResponse(w, http.StatusOK, trips)
}
