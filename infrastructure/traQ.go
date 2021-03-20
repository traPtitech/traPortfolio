package infrastructure

import (
	"encoding/json"
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
	traQCookie      = os.Getenv("TRAQ_COOKIE")
	traQAPIEndpoint = os.Getenv("TRAQ_API_ENDPOINT")
)

func init() {
	if traQCookie == "" {
		log.Fatal("the environment variable TRAQ_COOKIE should not be empty")
	}
	if traQAPIEndpoint == "" {
		log.Fatal("the environment variable TRAQ_API_ENDPOINT should not be empty")
	}
}

type TraQAPI struct {
	Client *http.Client
}

func NewTraQAPI() (*TraQAPI, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	cookies := []*http.Cookie{
		{
			Name:  "r_session",
			Value: traQCookie,
			Path:  "/",
		},
	}
	u, err := url.Parse(traQAPIEndpoint)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(u, cookies)
	return &TraQAPI{Client: &http.Client{Jar: jar}}, nil
}

func (traQ *TraQAPI) GetByID(id uuid.UUID, traQToken string) (*external.TraQUserResponse, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("invalid uuid")
	}

	res, err := apiGet(traQ.Client, traQAPIEndpoint, fmt.Sprintf("/users/%v", id))
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
