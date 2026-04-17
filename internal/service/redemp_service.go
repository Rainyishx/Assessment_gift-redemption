package service

import (
	"Assessment_gift-redemption/internal/model"
	"Assessment_gift-redemption/internal/repository"
	"errors"
	"time"
)

// sentinel errors
var (
	ErrStaffNotFound   = errors.New("staff pass ID not found")
	ErrAlreadyRedeemed = errors.New("team ahs already redeemed")
)

type RedempService interface {
	Redeem(staffPassID string) (model.Redemption, error)
}

// implement Redemp Service interface
type redempService struct {
	staffRepo  repository.StaffRepository
	redempRepo repository.RedempRepo
}

// create new redemp service instance
func NewRedempService(
	staffRepo repository.StaffRepository,
	redempRepo repository.RedempRepo,
) RedempService {
	return &redempService{
		staffRepo:  staffRepo,
		redempRepo: redempRepo,
	}
}

// process gift redeem request
func (s *redempService) Redeem(staffPassID string) (model.Redemption, error) {
	//look up staff pass ID
	teamName, found := s.staffRepo.FindTeamByStaffPassID(staffPassID)
	if !found {
		return model.Redemption{}, ErrStaffNotFound
	}

	//check if team has already redeemed
	if s.redempRepo.HasRedeemed(teamName) {
		return model.Redemption{}, ErrAlreadyRedeemed
	}

	//record redemp
	redemp := model.Redemption{
		TeamName:   teamName,
		RedeemedAt: time.Now().UnixMilli(),
	}

	if err := s.redempRepo.AddRedemp(redemp); err != nil {
		return model.Redemption{}, err
	}
	return redemp, nil
}
