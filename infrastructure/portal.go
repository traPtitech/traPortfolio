package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external_e2e"
)

const (
	cacheKey = "portalUsers"
)

type PortalConfig struct {
	cookie        string
	endpoint      string
	isDevelopment bool
}

func NewPortalConfig(cookie, endpoint string, isDevelopment bool) PortalConfig {
	return PortalConfig{
		cookie,
		endpoint,
		isDevelopment,
	}
}

type PortalAPI struct {
	Client *http.Client
	Cache  *cache.Cache
	conf   *PortalConfig
}

func NewPortalAPI(conf *PortalConfig) (external.PortalAPI, error) {
	if conf.isDevelopment {
		return mock_external_e2e.NewMockPortalAPI(), nil
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	cookies := []*http.Cookie{
		{
			Name:  "access_token",
			Value: conf.cookie,
			Path:  "/",
		},
	}
	u, err := url.Parse(conf.endpoint)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(u, cookies)
	c := cache.New(1*time.Hour, 2*time.Hour)
	return &PortalAPI{
		Client: &http.Client{Jar: jar},
		Cache:  c,
		conf:   conf,
	}, nil
}

func (portal *PortalAPI) GetAll() ([]*external.PortalUserResponse, error) {
	portalUsers, found := portal.Cache.Get(cacheKey)
	if found {
		return portalUsers.([]*external.PortalUserResponse), nil
	}

	res, err := apiGet(portal.Client, portal.conf.endpoint, "/user")
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

func (portal *PortalAPI) GetByID(traQID string) (*external.PortalUserResponse, error) {
	if traQID == "" {
		return nil, fmt.Errorf("invalid traQID")
	}

	res, err := apiGet(portal.Client, portal.conf.endpoint, fmt.Sprintf("/user/%v", traQID))
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
