package contest_service

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/domain"

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

func (s ContestService) CreateContest(ctx context.Context, args *repository.CreateContestArgs) (*domain.Contest, error) {
	uid := uuid.Must(uuid.NewV4())
	contest := &domain.Contest{
		ID:          uid,
		Name:        args.Name,
		Description: args.Description,
		Link:        args.Link,
		Since:       args.Since,
		Until:       args.Until,
	}
	contest, err := s.repo.Create(contest)
	if err != nil {
		return nil, err
	}
	return contest, nil
}
func (s ContestService) UpdateContest(ctx context.Context, args *repository.UpdateContestArgs) error {
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
		err := s.repo.Update(changes)
		if err != nil {
			return err
		}
	}
	return nil
}
