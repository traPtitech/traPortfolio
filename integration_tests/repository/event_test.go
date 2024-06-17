package repository

import (
	"context"
	"math/rand/v2"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/infrastructure/external/mock_external_e2e"
	irepository "github.com/traPtitech/traPortfolio/infrastructure/repository"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	urepository "github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
)

func TestEventRepository_GetEvents(t *testing.T) {
	t.Parallel()

	db := testutils.SetupGormDB(t)
	repo := irepository.NewEventRepository(db, mock_external_e2e.NewMockKnoqAPI())

	expected := make([]*domain.Event, 0)
	for _, e := range mockdata.MockKnoqEvents {
		expected = append(expected, &domain.Event{
			ID:        e.ID,
			Name:      e.Name,
			TimeStart: e.TimeStart,
			TimeEnd:   e.TimeEnd,
		})
	}

	got, err := repo.GetEvents(context.Background())
	assert.NoError(t, err)

	assert.ElementsMatch(t, expected, got)
}

func TestEventRepository_GetEvent(t *testing.T) {
	t.Parallel()

	db := testutils.SetupGormDB(t)
	repo := irepository.NewEventRepository(db, mock_external_e2e.NewMockKnoqAPI())

	levels := createRandomEventLevels(t, repo)
	selected := mockdata.MockKnoqEvents[rand.IntN(len(mockdata.MockKnoqEvents)-1)]
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
			Level:     level,
			TimeStart: selected.TimeStart,
			TimeEnd:   selected.TimeEnd,
		},
		Description: "",
		Place:       selected.Place,
		HostName:    hostName,
		GroupID:     selected.GroupID,
		RoomID:      selected.RoomID,
	}
	switch level {
	case domain.EventLevelPrivate:
		expected = nil
	case domain.EventLevelAnonymous:
		expected.HostName = nil
	case domain.EventLevelPublic:
	// do nothing
	default:
		t.Fatal("invalid level")
	}

	got, err := repo.GetEvent(context.Background(), selected.ID)
	if level == domain.EventLevelPrivate {
		assert.Error(t, err)
		assert.Nil(t, got)
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

// いいやり方が思いつかないので無視
// func TestCreateEventLevel(t *testing.T) {
// }

func TestEventRepository_UpdateEventLevel(t *testing.T) {
	t.Parallel()

	db := testutils.SetupGormDB(t)
	repo := irepository.NewEventRepository(db, mock_external_e2e.NewMockKnoqAPI())

	levels := createRandomEventLevels(t, repo)

	// choose randomly
	var selected *urepository.CreateEventLevelArgs
	for _, v := range levels {
		selected = v
		break
	}

	got, err := repo.GetEvent(context.Background(), selected.EventID)
	if selected.Level == domain.EventLevelPrivate {
		assert.Error(t, err)
		assert.Nil(t, got)
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, selected.Level, got.Level)

	updatedLevel := rand.N(domain.EventLevelLimit)
	err = repo.UpdateEventLevel(context.Background(), selected.EventID, &urepository.UpdateEventLevelArgs{
		Level: optional.From(updatedLevel),
	})
	assert.NoError(t, err)

	got, err = repo.GetEvent(context.Background(), selected.EventID)
	if updatedLevel == domain.EventLevelPrivate {
		assert.Error(t, err)
		assert.Nil(t, got)
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, updatedLevel, got.Level)
}

func TestEventRepository_GetUserEvents(t *testing.T) {
	t.Parallel()

	db := testutils.SetupGormDB(t)
	repo := irepository.NewEventRepository(db, mock_external_e2e.NewMockKnoqAPI())

	expected := make([]*domain.Event, 0)
	for _, e := range mockdata.MockKnoqEvents {
		expected = append(expected, &domain.Event{
			ID:        e.ID,
			Name:      e.Name,
			TimeStart: e.TimeStart,
			TimeEnd:   e.TimeEnd,
		})
	}

	got, err := repo.GetEvents(context.Background())
	assert.NoError(t, err)

	assert.ElementsMatch(t, expected, got)
}

// Create at least 1 event level.
func createRandomEventLevels(t *testing.T, repo urepository.EventRepository) map[uuid.UUID]*urepository.CreateEventLevelArgs {
	t.Helper()

	created := make(map[uuid.UUID]*urepository.CreateEventLevelArgs)
	for idx, e := range mockdata.MockKnoqEvents {
		if idx == 0 || random.Bool() {
			args := &urepository.CreateEventLevelArgs{
				EventID: e.ID,
				Level:   rand.N(domain.EventLevelLimit),
			}

			mustMakeEventLevel(t, repo, args)
			created[args.EventID] = args
		}
	}
	return created
}
