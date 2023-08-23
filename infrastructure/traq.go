package infrastructure

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external_e2e"
	"github.com/traPtitech/traPortfolio/util/config"
)

type TraQAPI struct {
	apiClient
}

func NewTraQAPI(conf *config.TraqConfig, isDevelopment bool) (external.TraQAPI, error) {
	if isDevelopment {
		return mock_external_e2e.NewMockTraQAPI(), nil
	}

	jar, err := newCookieJar(conf.API(), "r_session")
	if err != nil {
		return nil, err
	}

	return &TraQAPI{newAPIClient(jar, conf.API())}, nil
}

func (a *TraQAPI) GetTraqUsers(args *external.TraQGetAllArgs) ([]*external.TraQUserResponse, error) {
	res, err := a.apiGet(fmt.Sprintf("/users?include-suspended=%t&name=%s", args.IncludeSuspended, args.Name))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /users failed: %v", res.Status)
	}

	var usersResponse []*external.TraQUserResponse
	if err := json.NewDecoder(res.Body).Decode(&usersResponse); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}
	return usersResponse, nil
}

func (a *TraQAPI) GetByUserID(userID uuid.UUID) (*external.TraQUserResponse, error) {
	res, err := a.apiGet(fmt.Sprintf("/users/%v", userID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /users/%v failed: %v", userID, res.Status)
	}

	var userResponse external.TraQUserResponse
	if err := json.NewDecoder(res.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("decode failed: %v", err)
	}
	return &userResponse, nil
}

// Interface guards
var (
	_ external.TraQAPI = (*TraQAPI)(nil)
)
