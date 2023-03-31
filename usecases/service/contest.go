//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package service

import (
	"context"
	"errors"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type ContestService interface {
	GetContests(ctx context.Context) ([]*domain.Contest, error)
	GetContest(ctx context.Context, contestID uuid.UUID) (*domain.ContestDetail, error)
	CreateContest(ctx context.Context, args *repository.CreateContestArgs) (*domain.ContestDetail, error)
	UpdateContest(ctx context.Context, contestID uuid.UUID, args *repository.UpdateContestArgs) error
	DeleteContest(ctx context.Context, contestID uuid.UUID) error
	GetContestTeams(ctx context.Context, contestID uuid.UUID) ([]*domain.ContestTeam, error)
	GetContestTeam(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) (*domain.ContestTeamDetail, error)
	CreateContestTeam(ctx context.Context, contestID uuid.UUID, args *repository.CreateContestTeamArgs) (*domain.ContestTeamDetail, error)
	UpdateContestTeam(ctx context.Context, teamID uuid.UUID, args *repository.UpdateContestTeamArgs) error
	DeleteContestTeam(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) error
	GetContestTeamMembers(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) ([]*domain.User, error)
	AddContestTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error
	EditContestTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error
}

type contestService struct {
	repo repository.ContestRepository
}

func NewContestService(repo repository.ContestRepository) ContestService {
	return &contestService{
		repo,
	}
}

func (s *contestService) GetContests(ctx context.Context) ([]*domain.Contest, error) {
	contest, err := s.repo.GetContests(ctx)
	if err != nil {
		return nil, err
	}

	return contest, nil
}

func (s *contestService) GetContest(ctx context.Context, contestID uuid.UUID) (*domain.ContestDetail, error) {
	contest, err := s.repo.GetContest(ctx, contestID)
	if err != nil {
		return nil, err
	}

	teams, err := s.repo.GetContestTeams(ctx, contestID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	contest.ContestTeams = teams // TODO: repositoryで行うべきな気がする

	return contest, nil
}

func (s *contestService) CreateContest(ctx context.Context, args *repository.CreateContestArgs) (*domain.ContestDetail, error) {
	contest, err := s.repo.CreateContest(ctx, args)
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (s *contestService) UpdateContest(ctx context.Context, contestID uuid.UUID, args *repository.UpdateContestArgs) error {
	return s.repo.UpdateContest(ctx, contestID, args)
}

func (s *contestService) DeleteContest(ctx context.Context, contestID uuid.UUID) error {
	return s.repo.DeleteContest(ctx, contestID)
}

func (s *contestService) GetContestTeams(ctx context.Context, contestID uuid.UUID) ([]*domain.ContestTeam, error) {
	contestTeams, err := s.repo.GetContestTeams(ctx, contestID)
	if err != nil {
		return nil, err
	}
	return contestTeams, nil
}

func (s *contestService) GetContestTeam(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) (*domain.ContestTeamDetail, error) {
	contestTeam, err := s.repo.GetContestTeam(ctx, contestID, teamID)
	if err != nil {
		return nil, err
	}

	members, err := s.repo.GetContestTeamMembers(ctx, contestID, teamID)
	if err != nil {
		return nil, err
	}

	contestTeam.Members = members

	return contestTeam, nil
}

func (s *contestService) CreateContestTeam(ctx context.Context, contestID uuid.UUID, args *repository.CreateContestTeamArgs) (*domain.ContestTeamDetail, error) {
	contestTeam, err := s.repo.CreateContestTeam(ctx, contestID, args)
	if err != nil {
		return nil, err
	}
	return contestTeam, nil
}

func (s *contestService) UpdateContestTeam(ctx context.Context, teamID uuid.UUID, args *repository.UpdateContestTeamArgs) error {
	return s.repo.UpdateContestTeam(ctx, teamID, args)
}

func (s *contestService) DeleteContestTeam(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) error {
	return s.repo.DeleteContestTeam(ctx, contestID, teamID)
}

func (s *contestService) GetContestTeamMembers(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) ([]*domain.User, error) {
	members, err := s.repo.GetContestTeamMembers(ctx, contestID, teamID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (s *contestService) AddContestTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error {
	return s.repo.AddContestTeamMembers(ctx, teamID, memberIDs)
}

func (s *contestService) EditContestTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error {
	return s.repo.EditContestTeamMembers(ctx, teamID, memberIDs)
}

// Interface guards
var (
	_ ContestService = (*contestService)(nil)
)
