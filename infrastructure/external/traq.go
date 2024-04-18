//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package external

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/samber/lo"
	"github.com/traPtitech/go-traq"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/config"
)

type TraQUserResponse struct {
	ID    uuid.UUID        `json:"id"`
	Name  string           `json:"name"`
	State domain.TraQState `json:"state"`
}

type TraQGetAllArgs struct {
	IncludeSuspended bool
	Name             string
}

type TraQAPI interface {
	GetUsers(args *TraQGetAllArgs) ([]*TraQUserResponse, error)
	GetUser(userID uuid.UUID) (*TraQUserResponse, error)
}

type traQAPI struct {
	apiClient   *traq.APIClient
	accessToken string
	cache       *expirable.LRU[uuid.UUID, *TraQUserResponse]
}

func NewTraQAPI(conf config.TraqConfig) (TraQAPI, error) {
	apiClient := traq.NewAPIClient(traq.NewConfiguration())

	return &traQAPI{
		apiClient:   apiClient,
		accessToken: conf.AccessToken,
		cache:       expirable.NewLRU[uuid.UUID, *TraQUserResponse](5000, nil, 24*time.Hour),
	}, nil
}

func (a *traQAPI) GetUsers(args *TraQGetAllArgs) ([]*TraQUserResponse, error) {
	if users := a.cache.Values(); len(users) > 0 {
		if args.Name != "" {
			return lo.Filter(users, func(u *TraQUserResponse, i int) bool {
				return u.Name == args.Name
			}), nil
		}

		if !args.IncludeSuspended {
			return lo.Filter(users, func(u *TraQUserResponse, i int) bool {
				return u.State == domain.TraqStateActive
			}), nil
		}

		return users, nil
	}

	ctx := context.WithValue(context.Background(), traq.ContextAccessToken, a.accessToken)
	users, res, err := a.apiClient.UserApi.GetUsers(ctx).
		IncludeSuspended(args.IncludeSuspended).
		Name(args.Name).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("traQ API GetUsers error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("traQ API GetUsers invalid status: %s", res.Status)
	}

	usersResponse := lo.Map(users, func(u traq.User, i int) *TraQUserResponse {
		uid := uuid.FromStringOrNil(u.Id)
		r := &TraQUserResponse{
			ID:    uid,
			Name:  u.Name,
			State: domain.TraQState(u.State),
		}

		return r
	})

	for _, u := range usersResponse {
		a.cache.Add(u.ID, u)
	}

	return usersResponse, nil
}

func (a *traQAPI) GetUser(userID uuid.UUID) (*TraQUserResponse, error) {
	if user, ok := a.cache.Get(userID); ok {
		return user, nil
	}

	ctx := context.WithValue(context.Background(), traq.ContextAccessToken, a.accessToken)
	user, res, err := a.apiClient.UserApi.GetUser(ctx, userID.String()).Execute()
	if err != nil {
		return nil, fmt.Errorf("traQ API GetUser error: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("traQ API GetUser invalid status: %s", res.Status)
	}

	userResponse := &TraQUserResponse{
		ID:    uuid.FromStringOrNil(user.Id),
		Name:  user.Name,
		State: domain.TraQState(user.State),
	}

	a.cache.Add(userResponse.ID, userResponse)

	return userResponse, nil
}

// Interface guards
var (
	_ TraQAPI = (*traQAPI)(nil)
)
