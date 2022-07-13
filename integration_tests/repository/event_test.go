package repository

import (
	"math/rand"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external_e2e"
	irepository "github.com/traPtitech/traPortfolio/interfaces/repository"
	urepository "github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
)

func TestEventRepository_GetEvents(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("event_repository_get_events")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	repo := irepository.NewEventRepository(h, mock_external_e2e.NewMockKnoqAPI())

	expected := make([]*domain.Event, 0)
	for _, e := range mockdata.MockKnoqEvents {
		expected = append(expected, &domain.Event{
			ID:        e.ID,
			Name:      e.Name,
			TimeStart: e.TimeStart,
			TimeEnd:   e.TimeEnd,
		})
	}

	got, err := repo.GetEvents()
	assert.NoError(t, err)

	assert.ElementsMatch(t, expected, got)
}

func TestEventRepository_GetEvent(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("event_repository_get_event")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	repo := irepository.NewEventRepository(h, mock_external_e2e.NewMockKnoqAPI())

	levels := createRandomEventLevels(t, repo)
	selected := mockdata.MockKnoqEvents[rand.Intn(len(mockdata.MockKnoqEvents)-1)]
	hostName := make([]*domain.User, 0, len(selected.Admins))
	for _, aid := range selected.Admins {
		hostName = append(hostName, &domain.User{ID: aid})
	}

	var level domain.EventLevel
	if arg, ok := levels[selected.ID]; ok {
		level = arg.Level
	} else {
		level = domain.EventLevelAnonymous
	}
	expected := &domain.EventDetail{
		Event: domain.Event{
			ID:        selected.ID,
			Name:      selected.Name,
			TimeStart: selected.TimeStart,
			TimeEnd:   selected.TimeEnd,
		},
		Description: selected.Description,
		Place:       selected.Place,
		Level:       level,
		HostName:    hostName,
		GroupID:     selected.GroupID,
		RoomID:      selected.RoomID,
	}

	got, err := repo.GetEvent(selected.ID)
	assert.NoError(t, err)

	assert.Equal(t, expected, got)
}

// いいやり方が思いつかないので無視
// func TestCreateEventLevel(t *testing.T) {
// }

func TestEventRepository_UpdateEventLevel(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("event_repository_update_event_level")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	repo := irepository.NewEventRepository(h, mock_external_e2e.NewMockKnoqAPI())

	levels := createRandomEventLevels(t, repo)

	// choose randomly
	var selected *urepository.CreateEventLevelArgs
	for _, v := range levels {
		selected = v
		break
	}

	got, err := repo.GetEvent(selected.EventID)
	assert.NoError(t, err)
	assert.Equal(t, selected.Level, got.Level)

	updatedLevel := uint(rand.Intn(domain.EventLevelLimit))
	err = repo.UpdateEventLevel(selected.EventID, &urepository.UpdateEventLevelArgs{
		Level: optional.NewUint(updatedLevel, true),
	})
	assert.NoError(t, err)

	got, err = repo.GetEvent(selected.EventID)
	assert.NoError(t, err)
	assert.Equal(t, updatedLevel, uint(got.Level))
}

func TestEventRepository_GetUserEvents(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("event_repository_get_user_events")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	repo := irepository.NewEventRepository(h, mock_external_e2e.NewMockKnoqAPI())

	expected := make([]*domain.Event, 0)
	for _, e := range mockdata.MockKnoqEvents {
		expected = append(expected, &domain.Event{
			ID:        e.ID,
			Name:      e.Name,
			TimeStart: e.TimeStart,
			TimeEnd:   e.TimeEnd,
		})
	}

	got, err := repo.GetEvents()
	assert.NoError(t, err)

	assert.ElementsMatch(t, expected, got)
}

// Create at least 1 event level.
func createRandomEventLevels(t *testing.T, repo urepository.EventRepository) map[uuid.UUID]*urepository.CreateEventLevelArgs {
	created := make(map[uuid.UUID]*urepository.CreateEventLevelArgs)

	for idx, e := range mockdata.MockKnoqEvents {
		if idx == 0 || random.Bool() {
			args := &urepository.CreateEventLevelArgs{
				EventID: e.ID,
				Level:   domain.EventLevel(rand.Intn(domain.EventLevelLimit)),
			}

			mustMakeEventLevel(t, repo, args)
			created[args.EventID] = args
		}
	}
	return created
}
