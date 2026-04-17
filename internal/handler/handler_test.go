package handler_test

import (
	"Assessment_gift-redemption/internal/handler"
	"Assessment_gift-redemption/internal/model"
	"Assessment_gift-redemption/internal/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mock service for testing
type mockServ struct{}

// checking correct error response is thrown
func (m *mockServ) Redeem(staffPassID string) (model.Redemption, error) {
	switch staffPassID {
	case "valid":
		return model.Redemption{TeamName: "TEAM_A", RedeemedAt: 1700000000000}, nil
	case "unknown":
		return model.Redemption{}, service.ErrStaffNotFound
	case "redeemed":
		return model.Redemption{}, service.ErrAlreadyRedeemed
	default:
		return model.Redemption{}, nil
	}
}

// test helpers:
// newTestServer wires handler with mock manager and a router
func newTestServer() *http.ServeMux {
	h := handler.NewHandler(&mockServ{})
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

// postRedeem helps simplifies sending mock http post request to the /redeem endpoint
func postRedeem(t *testing.T, body any) *httptest.ResponseRecorder {
	//points failure location to calling test function instead of in this func
	t.Helper()

	//convert Go map/struct into JSON string
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal("failed to marshal request body", err)
	}

	//create a mock http post request
	req := httptest.NewRequest(http.MethodPost, "/redeem", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	//create a mock browser to record response
	w := httptest.NewRecorder()

	//send the mock req to server, and record output in 'w'
	newTestServer().ServeHTTP(w, req)

	return w
}

// tests
func TestRedeem_Success(t *testing.T) {
	w := postRedeem(t, map[string]string{"staff_pass_id": "valid"})

	//verify that handler translated success to 201 created status
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, but did not, got %d", w.Code)
	}

	//read JSON response body
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)

	//verify that the input data is correct
	if resp["team_name"] != "TEAM_A" {
		t.Errorf("expected team_name TEAM_A, but did not, got %v", resp["team_name"])
	}

	if resp["redeemed_at"] == nil {
		t.Error("expected redeemed_at in resp, but did not")
	}

	if resp["message"] == nil {
		t.Error("expected message in resp, but did not")
	}
}

func TestRedeem_Unknown(t *testing.T) {
	//send unknown to trigger ErrStaffNotFound
	w := postRedeem(t, map[string]string{"staff_pass_id": "unknown"})

	//verify that handler translate ErrStaffNotFound to 404
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, but did not, got %d", w.Code)
	}

	//check for error in resp
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["error"] == "" {
		t.Error("expected error message in response body, but did not")
	}
}

func TestRedeem_AlreadyRedeemed(t *testing.T) {
	//send redeemed to trigger ErrAlreadyRedeemed
	w := postRedeem(t, map[string]string{"staff_pass_id": "redeemed"})

	//verify that handler translate ErrAlreadyRedeemed to 409
	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, but did not, got %d", w.Code)
	}
}

func TestRedeem_MissingStaffPassID(t *testing.T) {
	//send empty JSON object
	w := postRedeem(t, map[string]string{})

	//verify that handler catches the missing id before sending it to service
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, but did not, got %d", w.Code)
	}
}

func TestRedeem_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/redeem", bytes.NewBufferString("notJSON"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	newTestServer().ServeHTTP(w, req)

	//verify that handler reject broken data and returns 400
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, but did not, got %d", w.Code)
	}
}

// test if server awake
func TestHealth(t *testing.T) {
	//create a simple get request to /health
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	newTestServer().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, but did not, got %d", w.Code)
	}
}

func TestRedeem_JSONContentType(t *testing.T) {
	w := postRedeem(t, map[string]string{"staff_pass_id": "valid"})

	//verify that writeJSON sets headers correctly
	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, but did not, got %s", ct)
	}
}
