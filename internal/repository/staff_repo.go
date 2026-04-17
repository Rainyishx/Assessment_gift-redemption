package repository

import (
	"Assessment_gift-redemption/internal/model"
	"encoding/csv"
	"os"
	"strconv"
)

type StaffRepository interface {
	FindTeamByStaffPassID(staffPassID string) (string, bool)
}

// using an in-memory map, implements StaffRepository interface
type staffRepository struct {
	mappings map[string]model.StaffMapping
}

// intialise repo and load data
func NewStaffRepository(filePath string) (StaffRepository, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	mappings := make(map[string]model.StaffMapping)

	//loading data into map, skipping header row[0]
	for _, row := range records[1:] {
		createdAt, _ := strconv.ParseInt(row[2], 10, 64)
		mappings[row[0]] = model.StaffMapping{
			StaffPassID: row[0],
			TeamName:    row[1],
			CreatedAt:   createdAt,
		}
	}
	return &staffRepository{mappings: mappings}, nil
}

// map lookup to find team with the staffpassID
func (s *staffRepository) FindTeamByStaffPassID(staffPassID string) (string, bool) {
	m, ok := s.mappings[staffPassID]
	if !ok {
		return "", false
	}
	return m.TeamName, true

}
