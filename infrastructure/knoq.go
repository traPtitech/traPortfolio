package infrastructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external_e2e"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/config"
)

type KnoqAPI struct {
	apiClient
}

func NewKnoqAPI(conf *config.KnoqConfig, isDevelopment bool) (external.KnoqAPI, error) {
	if isDevelopment {
		return &mock_external_e2e.MockKnoqAPI{}, nil
	}

	jar, err := newCookieJar(conf.API(), "session")
	if err != nil {
		return nil, err
	}

	return &KnoqAPI{newAPIClient(jar, conf.API())}, nil
}

func (knoq *KnoqAPI) GetAll() ([]*external.EventResponse, error) {
	res, err := knoq.apiGet("/events")
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
	res, err := knoq.apiGet(fmt.Sprintf("/events/%s", eventID))
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

	var er external.EventResponse
	if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
		return nil, err
	}
	return &er, nil
}

func (knoq *KnoqAPI) GetByUserID(userID uuid.UUID) ([]*external.EventResponse, error) {
	res, err := knoq.apiGet(fmt.Sprintf("/users/%s/events", userID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET /users/%s/events failed", userID)
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
