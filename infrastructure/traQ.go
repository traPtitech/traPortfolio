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
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external"
)

type TraQConfig struct {
	cookie   string
	endpoint string
}

func NewTraQConfig(cookie, endpoint string) TraQConfig {
	return TraQConfig{
		cookie,
		endpoint,
	}
}

type TraQAPI struct {
	Client *http.Client
	conf   *TraQConfig
}

func NewTraQAPI(conf *TraQConfig) (external.TraQAPI, error) {
	if isDevelopment {
		return &mock_external.MockTraQAPI{}, nil
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

func (traQ *TraQAPI) GetByID(id uuid.UUID) (*external.TraQUserResponse, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid uuid")
	}

	res, err := apiGet(traQ.Client, traQ.conf.endpoint, fmt.Sprintf("/users/%v", id))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /users/%v failed: %v", id, res.Status)
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
