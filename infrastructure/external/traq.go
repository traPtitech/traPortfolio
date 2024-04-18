//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/samber/lo"
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
	apiClient
	cache *expirable.LRU[uuid.UUID, *TraQUserResponse]
}

func NewTraQAPI(conf config.APIConfig) (TraQAPI, error) {
	jar, err := newCookieJar(conf, "r_session")
	if err != nil {
		return nil, err
	}

	return &traQAPI{
		apiClient: newAPIClient(jar, conf),
		cache:     expirable.NewLRU[uuid.UUID, *TraQUserResponse](5000, nil, 24*time.Hour),
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
				return u.State != domain.TraqStateActive
			}), nil
		}

		return users, nil
	}

	res, err := a.apiGet(fmt.Sprintf("/users?include-suspended=%t&name=%s", args.IncludeSuspended, args.Name))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /users failed: %v", res.Status)
	}

	var usersResponse []*TraQUserResponse
	if err := json.NewDecoder(res.Body).Decode(&usersResponse); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	for _, user := range usersResponse {
		a.cache.Add(user.ID, user)
	}

	return usersResponse, nil
}

func (a *traQAPI) GetUser(userID uuid.UUID) (*TraQUserResponse, error) {
	if user, ok := a.cache.Get(userID); ok {
		return user, nil
	}

	res, err := a.apiGet(fmt.Sprintf("/users/%v", userID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /users/%v failed: %v", userID, res.Status)
	}

	var userResponse TraQUserResponse
	if err := json.NewDecoder(res.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("decode failed: %v", err)
	}

	a.cache.Add(userResponse.ID, &userResponse)

	return &userResponse, nil
}

// Interface guards
var (
	_ TraQAPI = (*traQAPI)(nil)
)
