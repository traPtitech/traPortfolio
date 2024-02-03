//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/traPtitech/traPortfolio/util/config"
)

type PortalUserResponse struct {
	TraQID         string `json:"id"`
	RealName       string `json:"name"`
	AlphabeticName string `json:"alphabeticName"`
}

type PortalAPI interface {
	GetUsers() ([]*PortalUserResponse, error)
	GetUserByTraqID(traQID string) (*PortalUserResponse, error)
}

const (
	cacheKey = "portalUsers"
)

type portalAPI struct {
	apiClient
	cache *cache.Cache
}

func NewPortalAPI(conf config.APIConfig) (PortalAPI, error) {
	jar, err := newCookieJar(conf, "access_token")
	if err != nil {
		return nil, err
	}

	return &portalAPI{
		apiClient: newAPIClient(jar, conf),
		cache:     cache.New(1*time.Hour, 2*time.Hour),
	}, nil
}

func (a *portalAPI) GetUsers() ([]*PortalUserResponse, error) {
	portalUsers, found := a.cache.Get(cacheKey)
	if found {
		return portalUsers.([]*PortalUserResponse), nil
	}

	res, err := a.apiGet("/user")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /user failed: %v", res.Status)
	}
	var userResponses []*PortalUserResponse
	if err := json.NewDecoder(res.Body).Decode(&userResponses); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}
	a.cache.Set(cacheKey, userResponses, cache.DefaultExpiration)
	return userResponses, nil
}

func (a *portalAPI) GetUserByTraqID(traQID string) (*PortalUserResponse, error) {
	if traQID == "" {
		return nil, fmt.Errorf("invalid traQID")
	}

	res, err := a.apiGet(fmt.Sprintf("/user/%v", traQID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /user/%v failed: %v", traQID, res.Status)
	}

	var userResponse PortalUserResponse
	if err := json.NewDecoder(res.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}
	return &userResponse, nil
}

// Interface guards
var (
	_ PortalAPI = (*portalAPI)(nil)
)
