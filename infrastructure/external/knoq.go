//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/config"
)

type EventResponse struct {
	ID          uuid.UUID   `json:"eventId"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Place       string      `json:"place"`
	GroupID     uuid.UUID   `json:"groupId"`
	RoomID      uuid.UUID   `json:"roomId"`
	TimeStart   time.Time   `json:"timeStart"`
	TimeEnd     time.Time   `json:"timeEnd"`
	SharedRoom  bool        `json:"sharedRoom"`
	Admins      []uuid.UUID `json:"admins"`
}

type KnoqAPI interface {
	GetEvents() ([]*EventResponse, error)
	GetEvent(eventID uuid.UUID) (*EventResponse, error)
	GetEventsByUserID(userID uuid.UUID) ([]*EventResponse, error)
}

type knoqAPI struct {
	apiClient
}

func NewKnoqAPI(conf *config.KnoqConfig) (KnoqAPI, error) {
	jar, err := newCookieJar(conf.API(), "session")
	if err != nil {
		return nil, err
	}

	return &knoqAPI{newAPIClient(jar, conf.API())}, nil
}

func (a *knoqAPI) GetEvents() ([]*EventResponse, error) {
	res, err := a.apiGet("/events")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("GET /events failed")
	}

	var er []*EventResponse
	if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
		return nil, err
	}
	return er, nil
}

func (a *knoqAPI) GetEvent(eventID uuid.UUID) (*EventResponse, error) {
	res, err := a.apiGet(fmt.Sprintf("/events/%s", eventID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, repository.ErrNotFound
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /events/%s failed: %d", eventID, res.StatusCode)
	}

	var er EventResponse
	if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
		return nil, err
	}
	return &er, nil
}

func (a *knoqAPI) GetEventsByUserID(userID uuid.UUID) ([]*EventResponse, error) {
	res, err := a.apiGet(fmt.Sprintf("/users/%s/events", userID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /users/%s/events failed", userID)
	}

	var er []*EventResponse
	if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
		return nil, err
	}
	return er, nil
}

// Interface guards
var (
	_ KnoqAPI = (*knoqAPI)(nil)
)
