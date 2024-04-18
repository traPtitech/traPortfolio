//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package external

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
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
}

func NewTraQAPI(conf config.TraqConfig) (TraQAPI, error) {
	apiClient := traq.NewAPIClient(traq.NewConfiguration())

	return &traQAPI{
		apiClient:   apiClient,
		accessToken: conf.AccessToken,
	}, nil
}

func (a *traQAPI) GetUsers(args *TraQGetAllArgs) ([]*TraQUserResponse, error) {
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

	return usersResponse, nil
}

func (a *traQAPI) GetUser(userID uuid.UUID) (*TraQUserResponse, error) {
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

	return userResponse, nil
}

// Interface guards
var (
	_ TraQAPI = (*traQAPI)(nil)
)
