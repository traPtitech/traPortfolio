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
	GetContest(ctx context.Context, id uuid.UUID) (*domain.ContestDetail, error)
	CreateContest(ctx context.Context, args *repository.CreateContestArgs) (*domain.Contest, error)
	UpdateContest(ctx context.Context, id uuid.UUID, args *repository.UpdateContestArgs) error
	DeleteContest(ctx context.Context, id uuid.UUID) error
	GetContestTeams(ctx context.Context, contestID uuid.UUID) ([]*domain.ContestTeam, error)
	GetContestTeam(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) (*domain.ContestTeamDetail, error)
	CreateContestTeam(ctx context.Context, contestID uuid.UUID, args *repository.CreateContestTeamArgs) (*domain.ContestTeamDetail, error)
	UpdateContestTeam(ctx context.Context, teamID uuid.UUID, args *repository.UpdateContestTeamArgs) error
	DeleteContestTeam(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) error
	GetContestTeamMembers(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) ([]*domain.User, error)
	AddContestTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error
	DeleteContestTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error
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
	contest, err := s.repo.GetContests()
	if err != nil {
		return nil, err
	}

	return contest, nil
}

func (s *contestService) GetContest(ctx context.Context, id uuid.UUID) (*domain.ContestDetail, error) {
	contest, err := s.repo.GetContest(id)
	if err != nil {
		return nil, err
	}

	teams, err := s.repo.GetContestTeams(id)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	contest.Teams = teams // TODO: repositoryで行うべきな気がする

	return contest, nil
}

func (s *contestService) CreateContest(ctx context.Context, args *repository.CreateContestArgs) (*domain.Contest, error) {
	contest, err := s.repo.CreateContest(args)
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (s *contestService) UpdateContest(ctx context.Context, id uuid.UUID, args *repository.UpdateContestArgs) error {
	changes := map[string]interface{}{}
	if args.Name.Valid {
		changes["name"] = args.Name.String
	}
	if args.Description.Valid {
		changes["description"] = args.Description.String
	}
	if args.Link.Valid {
		changes["link"] = args.Link.String
	}
	if args.Since.Valid {
		changes["since"] = args.Since.Time
	}
	if args.Until.Valid {
		changes["until"] = args.Until.Time
	}
	if len(changes) > 0 {
		err := s.repo.UpdateContest(id, changes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *contestService) DeleteContest(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteContest(id); err != nil {
		return err
	}

	return nil
}

func (s *contestService) GetContestTeams(ctx context.Context, contestID uuid.UUID) ([]*domain.ContestTeam, error) {
	contestTeams, err := s.repo.GetContestTeams(contestID)
	if err != nil {
		return nil, err
	}
	return contestTeams, nil
}

func (s *contestService) GetContestTeam(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) (*domain.ContestTeamDetail, error) {
	contestTeam, err := s.repo.GetContestTeam(contestID, teamID)
	if err != nil {
		return nil, err
	}

	members, err := s.repo.GetContestTeamMembers(contestID, teamID)
	if err != nil {
		return nil, err
	}

	contestTeam.Members = members

	return contestTeam, nil
}

func (s *contestService) CreateContestTeam(ctx context.Context, contestID uuid.UUID, args *repository.CreateContestTeamArgs) (*domain.ContestTeamDetail, error) {
	contestTeam, err := s.repo.CreateContestTeam(contestID, args)
	if err != nil {
		return nil, err
	}
	return contestTeam, nil
}

func (s *contestService) UpdateContestTeam(ctx context.Context, teamID uuid.UUID, args *repository.UpdateContestTeamArgs) error {
	changes := map[string]interface{}{}
	if args.Name.Valid {
		changes["name"] = args.Name.String
	}
	if args.Description.Valid {
		changes["description"] = args.Description.String
	}
	if args.Link.Valid {
		changes["link"] = args.Link.String
	}
	if args.Result.Valid {
		changes["result"] = args.Result.String
	}
	if len(changes) > 0 {
		err := s.repo.UpdateContestTeam(teamID, changes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *contestService) DeleteContestTeam(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) error {
	if err := s.repo.DeleteContestTeam(contestID, teamID); err != nil {
		return err
	}

	return nil
}

func (s *contestService) GetContestTeamMembers(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) ([]*domain.User, error) {
	return s.repo.GetContestTeamMembers(contestID, teamID)
}

func (s *contestService) AddContestTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error {
	err := s.repo.AddContestTeamMembers(teamID, memberIDs)
	return err
}

func (s *contestService) DeleteContestTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error {
	err := s.repo.DeleteContestTeamMembers(teamID, memberIDs)
	return err
}

// Interface guards
var (
	_ ContestService = (*contestService)(nil)
)
