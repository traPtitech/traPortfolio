package service

import (
	"context"
	"net/http"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/labstack/echo/v4"

	"github.com/traPtitech/traPortfolio/interfaces/repository/model"

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

func (s *ContestService) CreateContest(ctx context.Context, args *repository.CreateContestArgs) (*model.Contest, error) {
	uid := uuid.Must(uuid.NewV4())
	contest := &model.Contest{
		ID:          uid,
		Name:        args.Name,
		Description: args.Description,
		Link:        args.Link,
		Since:       args.Since,
		Until:       args.Until,
	}
	contest, err := s.repo.CreateContest(contest)
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (s *ContestService) UpdateContest(ctx context.Context, id uuid.UUID, args *repository.UpdateContestArgs) error {
	if id == uuid.Nil {
		return repository.ErrInvalidID
	}
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
		if err != nil && err == repository.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ContestService) CreateContestTeam(ctx context.Context, contestID uuid.UUID, args repository.CreateContestTeamArgs) (*domain.ContestTeamDetail, error) {
	contestTeam, err := s.repo.CreateContestTeam(contestID, args)
	if err != nil {
		return nil, err
	}
	return contestTeam, nil
}
