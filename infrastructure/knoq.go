package infrastructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external_e2e"
)

type KnoQConfig struct {
	cookie        string
	endpoint      string
	isDevelopment bool
}

func NewKnoqConfig(cookie, endpoint string, isDevelopment bool) KnoQConfig {
	return KnoQConfig{
		cookie,
		endpoint,
		isDevelopment,
	}
}

type KnoqAPI struct {
	Client *http.Client
	conf   *KnoQConfig
}

func NewKnoqAPI(conf *KnoQConfig) (external.KnoqAPI, error) {
	if conf.isDevelopment {
		return &mock_external_e2e.MockKnoqAPI{}, nil
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	cookies := []*http.Cookie{
		{
			Name:  "session",
			Value: conf.cookie,
			Path:  "/",
		},
	}
	u, err := url.Parse(conf.endpoint)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(u, cookies)
	return &KnoqAPI{Client: &http.Client{Jar: jar}, conf: conf}, nil
}

func (knoq *KnoqAPI) GetAll() ([]*external.EventResponse, error) {
	res, err := apiGet(knoq.Client, knoq.conf.endpoint, "/events")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("GET /events failed")
	}

	var er []*external.EventResponse
	if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
		return nil, err
	}
	return er, nil
}

func (knoq *KnoqAPI) GetByEventID(eventID uuid.UUID) (*external.EventResponse, error) {
	res, err := apiGet(knoq.Client, knoq.conf.endpoint, fmt.Sprintf("/events/%v", eventID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /events/%v failed", eventID)
	}

	var er external.EventResponse
	if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
		return nil, err
	}
	return &er, nil
}

func (knoq *KnoqAPI) GetByUserID(userID uuid.UUID) ([]*external.EventResponse, error) {
	res, err := apiGet(knoq.Client, knoq.conf.endpoint, fmt.Sprintf("/users/%v/events", userID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /users/%v/events failed", userID)
	}

	var er []*external.EventResponse
	if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
		return nil, err
	}
	return er, nil
}

// Interface guards
var (
	_ external.KnoqAPI = (*KnoqAPI)(nil)
)
