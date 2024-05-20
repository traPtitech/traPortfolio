//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package external

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/go-traq"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/util/config"
)

type TraQUserResponse struct {
	ID    uuid.UUID
	Name  string
	State domain.TraQState
	Bot   bool
}

type TraQGetAllArgs struct {
	IncludeSuspended bool
	Name             string
}

type TraQAPI interface {
	GetUsers(args *TraQGetAllArgs) ([]*TraQUserResponse, error)
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
			Bot:   u.Bot,
		}

		return r
	})

	return usersResponse, nil
}
