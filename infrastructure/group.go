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
)

type GroupConfig struct {
	cookie   string
	endpoint string
}

func NewgGoupConfig(cookie, endpoint string) GroupConfig {
	return GroupConfig{
		cookie,
		endpoint,
	}
}

type GroupAPI struct {
	Client *http.Client
	conf   *GroupConfig
}

func NewGroupAPI(conf *GroupConfig) (external.GroupAPI, error) {
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
	return &GroupAPI{Client: &http.Client{Jar: jar}, conf: conf}, nil
}

func (group *GroupAPI) GetGroupsByID(userID uuid.UUID) ([]*external.GroupUserResponse, error) {
	res, err := apiGet(group.Client, group.conf.endpoint, fmt.Sprintf("/users/%v/groups", userID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /users/%v/groups failed", userID)
	}

	var er []*external.GroupUserResponse
	if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
		return nil, err
	}
	return er, nil
}

func (group *GroupAPI) GetAllGroups() ([]*external.GroupsResponse, error) {
	res, err := apiGet(group.Client, group.conf.endpoint, fmt.Sprintf("/groups"))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /groups failed")
	}

	var er []*external.GroupsResponse
	if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
		return nil, err
	}
	return er, nil
}

func (group *GroupAPI) GetGroup(groupID uuid.UUID) (*external.GroupDetailResponse, error) {
	res, err := apiGet(group.Client, group.conf.endpoint, fmt.Sprintf("/groups/%d", groupID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /groups/%d failed", groupID)
	}

	var er *external.GroupDetailResponse
	if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
		return nil, err
	}
	return er, nil
}

// Interface guards
var (
	_ external.GroupAPI = (*GroupAPI)(nil)
)
