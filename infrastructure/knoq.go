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
)

type KnoQConfig struct {
	cookie   string
	endpoint string
}

func NewKnoqConfig(cookie, endpoint string) KnoQConfig {
	return KnoQConfig{
		cookie,
		endpoint,
	}
}

type KnoqAPI struct {
	Client *http.Client
	conf   *KnoQConfig
}

func NewKnoqAPI(conf *KnoQConfig) (external.KnoqAPI, error) {
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

func (knoq *KnoqAPI) GetByID(id uuid.UUID) (*external.EventResponse, error) {
	if id == uuid.Nil {
		return nil, nil
	}

	res, err := apiGet(knoq.Client, knoq.conf.endpoint, fmt.Sprintf("/events/%v", id))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /events/%v failed", id)
	}

	var er external.EventResponse
	if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
		return nil, err
	}
	return &er, nil
}

func (knoq *KnoqAPI) GetByUserID(id uuid.UUID) ([]*external.EventResponse, error) {
	if id == uuid.Nil {
		return nil, nil
	}

	res, err := apiGet(knoq.Client, knoq.conf.endpoint, fmt.Sprintf("/users/%v/events", id))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /users/%v/events failed", id)
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
