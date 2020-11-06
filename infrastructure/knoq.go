package infrastructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/external"
)

var (
	knoQCookie      = os.Getenv("KNOQ_COOKIE")
	knoQAPIEndpoint = os.Getenv("KNOQ_API_ENDPOINT")
)

type KnoqAPI struct {
	Client *http.Client
}

func NewKnoqAPI() (*KnoqAPI, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	cookies := []*http.Cookie{
		{
			Name:  "session",
			Value: knoQCookie,
			Path:  "/",
		},
	}
	u, err := url.Parse(knoQAPIEndpoint)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(u, cookies)
	return &KnoqAPI{Client: &http.Client{Jar: jar}}, nil
}

func (knoq *KnoqAPI) GetAll() ([]*external.EventResponse, error) {
	res, err := apiGet(knoq.Client, "/events")
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

	res, err := apiGet(knoq.Client, fmt.Sprintf("/events/%v", id))
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

func apiGet(client *http.Client, path string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", knoQAPIEndpoint, path), nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
