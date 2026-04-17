package repository

import (
	"Assessment_gift-redemption/internal/model"
	"encoding/csv"
	"os"
	"strconv"
	"sync"
)

type RedempRepo interface {
	HasRedeemed(teamName string) bool
	AddRedemp(record model.Redemption) error
}

// implements RedempRepo interface
type redempRepo struct {
	//to ensure only one request at a time
	mu          sync.Mutex
	redemptions map[string]model.Redemption
	filePath    string
}

// initialise repo and load any existing data
func NewRedempRepo(filePath string) (RedempRepo, error) {
	repo := &redempRepo{
		redemptions: make(map[string]model.Redemption),
		filePath:    filePath,
	}

	f, err := os.Open(filePath)
	if err != nil {
		//return empty repo if csv doesn't exist
		if os.IsNotExist(err) {
			return repo, nil
		}
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	//loading any existing data into map, skipping header row[0]
	for _, row := range records[1:] {
		redeemedAt, _ := strconv.ParseInt(row[1], 10, 64)
		repo.redemptions[row[0]] = model.Redemption{
			TeamName:   row[0],
			RedeemedAt: redeemedAt,
		}
	}
	return repo, nil

}

// check if team has redeemed before
func (repo *redempRepo) HasRedeemed(teamName string) bool {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	_, ok := repo.redemptions[teamName]
	return ok
}

// Adding a redemption record to memory and updating the csv
func (repo *redempRepo) AddRedemp(redemp model.Redemption) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.redemptions[redemp.TeamName] = redemp

	//create the file if it doesnt exist
	f, err := os.Create(repo.filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	//write the header
	writer.Write([]string{"team_Name", "redeemed_at"})
	//write the data
	for _, rd := range repo.redemptions {
		writer.Write([]string{rd.TeamName, strconv.FormatInt(rd.RedeemedAt, 10)})
	}
	return nil
}
