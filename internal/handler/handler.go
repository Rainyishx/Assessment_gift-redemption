package handler

import (
	"Assessment_gift-redemption/internal/service"
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc service.RedempService
}

func NewHandler(svc service.RedempService) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes attaches all routes to the given ServeMux
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /redeem", h.Redeem)
	mux.HandleFunc("GET /Health", h.Health)
}

type redeemRequest struct {
	StaffPassID string `json:"staff_pass_id"`
}

type redeemResponse struct {
	TeamName   string `json:"team_name"`
	RedeemedAt int64  `json:"redeemed_at"`
	Message    string `json:"message"`
}

type errorResponse struct {
	Error string `json:"error"`
}

//endpoints

// writing status codes and encoding JSON
func writeJSON(w http.ResponseWriter, status int, body any) {
	//informs the browser to expect JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func (h *Handler) Redeem(w http.ResponseWriter, r *http.Request) {
	var req redeemRequest

	//read incoming json and convert it into Go struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "invalid JSON body",
		})
		return
	}

	//validate input
	if req.StaffPassID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "staff_pass_id is required",
		})
		return
	}

	//pass ID to service and wait for result
	redemp, err := h.svc.Redeem(req.StaffPassID)

	//handle errors by translating into http codes
	if err != nil {
		switch err {
		//if not found, send HTTP 404
		case service.ErrStaffNotFound:
			writeJSON(w, http.StatusNotFound, errorResponse{
				Error: err.Error(),
			})

			//if already redeemed, send HTTP 409
		case service.ErrAlreadyRedeemed:
			writeJSON(w, http.StatusConflict, errorResponse{
				Error: err.Error(),
			})

		//for other errors, send HTTP 500
		default:
			writeJSON(w, http.StatusInternalServerError, errorResponse{
				Error: "internal server error",
			})
		}
		return
	}

	//if there is no error, send HTTP 201 and the data
	writeJSON(w, http.StatusCreated, redeemResponse{
		TeamName:   redemp.TeamName,
		RedeemedAt: redemp.RedeemedAt,
		Message:    "Gift redeemed successfully for team " + redemp.TeamName,
	})

}

// check if server is awake
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
