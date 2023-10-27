//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package external

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/config"
)

type TraQUserResponse struct {
	ID    uuid.UUID        `json:"id"`
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
}

func NewTraQAPI(conf *config.TraqConfig) (TraQAPI, error) {
	jar, err := newCookieJar(conf.API(), "r_session")
	if err != nil {
		return nil, err
	}

	return &traQAPI{newAPIClient(jar, conf.API())}, nil
}

func (a *traQAPI) GetUsers(args *TraQGetAllArgs) ([]*TraQUserResponse, error) {
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
	return usersResponse, nil
}

func (a *traQAPI) GetUser(userID uuid.UUID) (*TraQUserResponse, error) {
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
	return &userResponse, nil
}

// Interface guards
var (
	_ TraQAPI = (*traQAPI)(nil)
)
