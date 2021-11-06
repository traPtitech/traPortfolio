package service

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type ContestService struct {
	repo repository.ContestRepository
}

func NewContestService(repo repository.ContestRepository) ContestService {
	return ContestService{
		repo,
	}
}

func (s *ContestService) GetContests(ctx context.Context) ([]*domain.Contest, error) {
	contest, err := s.repo.GetContests()
	if err != nil {
		return nil, err
	}

	return contest, nil
}

func (s *ContestService) GetContest(ctx context.Context, id uuid.UUID) (*domain.ContestDetail, error) {
	contest, err := s.repo.GetContest(id)
	if err != nil {
		return nil, err
	}

	teams, err := s.repo.GetContestTeams(id)
	if err != nil {
		return nil, err
	}

	contest.Teams = teams

	return contest, nil
}

func (s *ContestService) CreateContest(ctx context.Context, args *repository.CreateContestArgs) (*domain.Contest, error) {
	contest, err := s.repo.CreateContest(args)
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (s *ContestService) UpdateContest(ctx context.Context, id uuid.UUID, args *repository.UpdateContestArgs) error {
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

func (s *ContestService) DeleteContest(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteContest(id); err != nil {
		return err
	}

	return nil
}

func (s *ContestService) GetContestTeams(ctx context.Context, contestID uuid.UUID) ([]*domain.ContestTeam, error) {
	contestTeams, err := s.repo.GetContestTeams(contestID)
	if err != nil {
		return nil, err
	}
	return contestTeams, nil
}

func (s *ContestService) GetContestTeam(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) (*domain.ContestTeamDetail, error) {
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

func (s *ContestService) CreateContestTeam(ctx context.Context, contestID uuid.UUID, args *repository.CreateContestTeamArgs) (*domain.ContestTeamDetail, error) {
	contestTeam, err := s.repo.CreateContestTeam(contestID, args)
	if err != nil {
		return nil, err
	}
	return contestTeam, nil
}

func (s *ContestService) UpdateContestTeam(ctx context.Context, teamID uuid.UUID, args *repository.UpdateContestTeamArgs) error {
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

func (s *ContestService) GetContestTeamMembers(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) ([]*domain.User, error) {
	return s.repo.GetContestTeamMembers(contestID, teamID)
}

func (s *ContestService) AddContestTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error {
	err := s.repo.AddContestTeamMembers(teamID, memberIDs)
	return err
}

func (s *ContestService) DeleteContestTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error {
	err := s.repo.DeleteContestTeamMembers(teamID, memberIDs)
	return err
}
