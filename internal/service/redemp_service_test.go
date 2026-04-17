package service_test

import (
	"Assessment_gift-redemption/internal/model"
	"Assessment_gift-redemption/internal/service"
	"testing"
)

// mock staff repo acts as the staff csv using a map
type mockStaffRepo struct {
	//maps StaffPassID -> TeamName
	data map[string]string
}

// implement the staff repo interface
func (m *mockStaffRepo) FindTeamByStaffPassID(id string) (string, bool) {
	teamName, found := m.data[id]
	return teamName, found
}

// mock redemp repo acts as the redemption csv using a map
type mockRedempRepo struct {
	//maps TeamName -> true if redeemed
	redeemed map[string]bool
}

// implement the redemp repo interface
func (m *mockRedempRepo) HasRedeemed(teamName string) bool {
	return m.redeemed[teamName]
}

func (m *mockRedempRepo) AddRedemp(record model.Redemption) error {
	m.redeemed[record.TeamName] = true
	return nil
}

// test cases
func TestRedempService_Success(t *testing.T) {
	//create mock databases
	mockStaff := &mockStaffRepo{
		data: map[string]string{"STAFF_001": "TEAM_A"},
	}

	mockRedemp := &mockRedempRepo{
		//empty since no one has redeemed anything
		redeemed: make(map[string]bool),
	}

	//put the mock database into the service
	svc := service.NewRedempService(mockStaff, mockRedemp)

	//try to redeem
	redemption, err := svc.Redeem("STAFF_001")
	if err != nil {
		t.Fatal("expected success, but did not", err)
	}

	if redemption.TeamName != "TEAM_A" {
		t.Errorf("expected TEAM_A, but got '%s'", redemption.TeamName)
	}
}

func TestRedempService_StaffNotFound(t *testing.T) {
	//staff database is empty
	mockStaff := &mockStaffRepo{
		data: make(map[string]string),
	}

	mockRedemp := &mockRedempRepo{
		redeemed: make(map[string]bool),
	}

	svc := service.NewRedempService(mockStaff, mockRedemp)

	//try to redeem with invalid staff pass ID
	_, err := svc.Redeem("unknown")

	//verify it returns the sentinel error defined previously
	if err != service.ErrStaffNotFound {
		t.Error("expected ErrStaffNotFound, but did not", err)
	}
}

func TestRedempService_AlreadyRedeemed(t *testing.T) {
	//staff pass id is valid and exists, but already redeemed
	mockStaff := mockStaffRepo{
		data: map[string]string{"STAFF_001": "TEAM_A"},
	}

	mockRedemp := &mockRedempRepo{
		redeemed: map[string]bool{"TEAM_A": true},
	}

	svc := service.NewRedempService(&mockStaff, mockRedemp)

	//try to redeem
	_, err := svc.Redeem("STAFF_001")

	//verify sentinel error
	if err != service.ErrAlreadyRedeemed {
		t.Error("expected ErrAlreadyRedeemed, but did not", err)
	}
}
