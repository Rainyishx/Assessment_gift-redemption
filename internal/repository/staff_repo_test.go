package repository

import (
	"testing"
)

func TestSatffRepo(t *testing.T) {
	repo, err := NewStaffRepository("../../data/staffmapping.csv")
	if err != nil {
		t.Fatal("failed to initalise repo:", err)
	}

	//finding existing staff member
	t.Run("Valid Staff ID", func(t *testing.T) {
		teamName, found := repo.FindTeamByStaffPassID("STAFF_001")

		if !found {
			t.Errorf("expected to find STAFF_001, but did not")
		}
		if teamName != "TEAM_A" {
			t.Errorf("expected to find TEAM_A, got '%s'", teamName)
		}
	})

	//finding staff that doesn't exist
	t.Run("Invalid Staff ID", func(t *testing.T) {
		_, found := repo.FindTeamByStaffPassID("unknown")
		if found {
			t.Errorf("not expected to find unknown, but did")
		}

	})
}
