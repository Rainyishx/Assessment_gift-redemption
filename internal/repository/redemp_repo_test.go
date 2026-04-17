package repository

import (
	"Assessment_gift-redemption/internal/model"
	"os"
	"sync"
	"testing"
	"time"
)

// generate a temp file path with nonexistent file for testing creation of redemp file
func getTempFilePath(t *testing.T) string {
	file, err := os.CreateTemp("", "test_redemp_*.csv")
	if err != nil {
		t.Fatal("failed to create temp file", err)
	}
	file.Close()
	os.Remove(file.Name())
	return file.Name()
}

func TestRedempRepo(t *testing.T) {
	tempPath := getTempFilePath(t)
	//clean up after test ends
	defer os.Remove(tempPath)

	//testing file not found
	repo, err := NewRedempRepo(tempPath)
	if err != nil {
		t.Fatal("failed to initialise repo", err)
	}

	//verify team hasnt redeemed yet
	if repo.HasRedeemed("TEAM_A") {
		t.Errorf("expected 'TEAM_A' to not have redeemed yet")
	}

	//testing add redemption
	now := time.Now().UnixMilli()
	record := model.Redemption{
		TeamName:   "TEAM_A",
		RedeemedAt: now,
	}

	if err := repo.AddRedemp(record); err != nil {
		t.Fatal("failed to add redemption", err)
	}

	//verify team has redeemed
	if !repo.HasRedeemed("TEAM_A") {
		t.Errorf("expected TEAM_A to have redeemed")
	}

	//initialise new repo to test writing to csv succeeded
	repo2, err := NewRedempRepo(tempPath)
	if err != nil {
		t.Fatal("failed to initialise second repo", err)
	}
	if !repo2.HasRedeemed("TEAM_A") {
		t.Errorf("expected TEAM_A to have redeemed after loading data from csv, but it was not")
	}
}

func TestRedempRepo_Concurrency(t *testing.T) {
	tempPath := getTempFilePath(t)
	defer os.Remove(tempPath)

	repo, err := NewRedempRepo(tempPath)
	if err != nil {
		t.Fatal("failed to initialise repo", err)
	}

	//use waitgroup to ensure the test waits for all parallel requests to finish
	var wg sync.WaitGroup
	teams := []string{"TEAM_A", "TEAM_B", "TEAM_C", "TEAM_D", "TEAM_E"}

	//simulate simultaneos requests
	for _, team := range teams {
		wg.Add(1)
		go func(tName string) {
			defer wg.Done()
			//checking mutex is working
			repo.AddRedemp(model.Redemption{
				TeamName:   tName,
				RedeemedAt: time.Now().UnixMilli(),
			})
		}(team)
	}
	wg.Wait()

	//verify that teams are successfully recorded
	for _, team := range teams {
		if !repo.HasRedeemed(team) {
			t.Errorf("expected %s to be recorded, but it did not", team)
		}
	}
}
