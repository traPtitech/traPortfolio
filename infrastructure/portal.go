package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/traPtitech/traPortfolio/interfaces/external"
)

var (
	portalCookie      = os.Getenv("PORTAL_COOKIE")
	portalAPIEndpoint = os.Getenv("PORTAL_API_ENDPOINT")
	cacheKey          = "portalUsers"
)

func init() {
	if portalCookie == "" {
		log.Fatal("the environment variable PORTAL_COOKIE should not be empty")
	}
	if portalAPIEndpoint == "" {
		log.Fatal("the environment variable PORTAL_API_ENDPOINT should not be empty")
	}
}

type PortalAPI struct {
	Client *http.Client
	Cache  *cache.Cache
}

func NewPortalAPI() (external.PortalAPI, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	cookies := []*http.Cookie{
		{
			Name:  "access_token",
			Value: portalCookie,
			Path:  "/",
		},
	}
	u, err := url.Parse(portalAPIEndpoint)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(u, cookies)
	c := cache.New(1*time.Hour, 2*time.Hour)
	return &PortalAPI{
		Client: &http.Client{Jar: jar},
		Cache:  c,
	}, nil
}

func (portal *PortalAPI) GetAll() ([]*external.PortalUserResponse, error) {
	portalUsers, found := portal.Cache.Get(cacheKey)
	if found {
		return portalUsers.([]*external.PortalUserResponse), nil
	}

	res, err := apiGet(portal.Client, portalAPIEndpoint, "/user")
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

	res, err := apiGet(portal.Client, portalAPIEndpoint, fmt.Sprintf("/user/%v", traQID))
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
