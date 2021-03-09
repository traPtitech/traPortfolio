package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/traPtitech/traPortfolio/interfaces/external"
)

var (
	portalCookie      = os.Getenv("PORTAL_COOKIE")
	portalAPIEndpoint = os.Getenv("PORTAL_API_ENDPOINT")
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
}

func NewPortalAPI() (*PortalAPI, error) {
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
	return &PortalAPI{Client: &http.Client{Jar: jar}}, nil
}

func (portal *PortalAPI) GetAll(portalToken string) ([]*external.PortalUserResponse, error) {
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
	return userResponses, nil
}

func (portal *PortalAPI) GetByID(name string, portalToken string) (*external.PortalUserResponse, error) {
	res, err := apiGet(portal.Client, portalAPIEndpoint, fmt.Sprintf("/user/%s", name))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(("GET /user/%s failed: %v"), name, res.Status)
	}

	var userResponse external.PortalUserResponse
	if err := json.NewDecoder(res.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("decode failed: %v", err)
	}
	return &userResponse, err
}
