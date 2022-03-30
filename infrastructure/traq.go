package infrastructure

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external_e2e"
)

type TraQConfig struct {
	cookie        string
	endpoint      string
	isDevelopment bool
}

func NewTraQConfig(cookie, endpoint string, isDevelopment bool) TraQConfig {
	return TraQConfig{
		cookie,
		endpoint,
		isDevelopment,
	}
}

type TraQAPI struct {
	Client *http.Client
	conf   *TraQConfig
}

func NewTraQAPI(conf *TraQConfig) (external.TraQAPI, error) {
	if conf.isDevelopment {
		return mock_external_e2e.NewMockTraQAPI(), nil
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	cookies := []*http.Cookie{
		{
			Name:  "r_session",
			Value: conf.cookie,
			Path:  "/",
		},
	}
	u, err := url.Parse(conf.endpoint)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(u, cookies)
	return &TraQAPI{Client: &http.Client{Jar: jar}, conf: conf}, nil
}

func (traQ *TraQAPI) GetAll(args *external.TraQGetAllArgs) ([]*external.TraQUserResponse, error) {
	res, err := apiGet(traQ.Client, traQ.conf.endpoint, fmt.Sprintf("/users?include-suspended=%t&name=%s", args.IncludeSuspended, args.Name))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /users failed: %v", res.Status)
	}

	var usersResponse []*external.TraQUserResponse
	if err := json.NewDecoder(res.Body).Decode(&usersResponse); err != nil {
		return nil, fmt.Errorf("decode failed: %v", err)
	}
	return usersResponse, nil
}

func (traQ *TraQAPI) GetByID(userID uuid.UUID) (*external.TraQUserResponse, error) {
	res, err := apiGet(traQ.Client, traQ.conf.endpoint, fmt.Sprintf("/users/%v", userID))
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
