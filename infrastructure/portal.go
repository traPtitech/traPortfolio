package infrastructure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external_e2e"
	"github.com/traPtitech/traPortfolio/util/config"
)

const (
	cacheKey = "portalUsers"
)

type PortalAPI struct {
	apiClient
	Cache *cache.Cache
}

func NewPortalAPI(conf *config.PortalConfig, isDevelopment bool) (external.PortalAPI, error) {
	if isDevelopment {
		return mock_external_e2e.NewMockPortalAPI(), nil
	}

	jar, err := newCookieJar(conf.API(), "access_token")
	if err != nil {
		return nil, err
	}

	return &PortalAPI{
		apiClient: newAPIClient(jar, conf.API()),
		Cache:     cache.New(1*time.Hour, 2*time.Hour),
	}, nil
}

func (portal *PortalAPI) GetAll() ([]*external.PortalUserResponse, error) {
	portalUsers, found := portal.Cache.Get(cacheKey)
	if found {
		return portalUsers.([]*external.PortalUserResponse), nil
	}

	res, err := portal.apiGet("/user")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /user failed: %v", res.Status)
	}
	var userResponses []*external.PortalUserResponse
	if err := json.NewDecoder(res.Body).Decode(&userResponses); err != nil {
		return nil, fmt.Errorf("decode failed: %v", err)
	}
	portal.Cache.Set(cacheKey, userResponses, cache.DefaultExpiration)
	return userResponses, nil
}

func (portal *PortalAPI) GetByTraqID(traQID string) (*external.PortalUserResponse, error) {
	if traQID == "" {
		return nil, fmt.Errorf("invalid traQID")
	}

	res, err := portal.apiGet(fmt.Sprintf("/user/%v", traQID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /user/%v failed: %v", traQID, res.Status)
	}

	var userResponse external.PortalUserResponse
	if err := json.NewDecoder(res.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("decode failed: %v", err)
	}
	return &userResponse, nil
}

// Interface guards
var (
	_ external.PortalAPI = (*PortalAPI)(nil)
)
