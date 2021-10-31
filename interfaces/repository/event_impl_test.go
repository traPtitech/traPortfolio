package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/traPtitech/traPortfolio/util/random"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/database/mock_database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

var (
	sampleUUID = uuid.FromStringOrNil("3fa85f64-5717-4562-b3fc-2c963f66afa6")
	sampleTime = time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
)

type mockEventRepositoryFields struct {
	h   database.SQLHandler
	api external.KnoqAPI
}

func newMockEventRepositoryFields() mockEventRepositoryFields {
	return mockEventRepositoryFields{
		h:   mock_database.NewMockSQLHandler(),
		api: mock_external.NewMockKnoqAPI(),
	}
}

func TestEventRepository_GetEvent(t *testing.T) {
	t.Parallel()
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.EventDetail
		setup     func(f mockEventRepositoryFields, args args, want *domain.EventDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				id: sampleUUID,
			},
			want: &domain.EventDetail{
				Event: domain.Event{
					ID:        sampleUUID,
					Name:      "第n回進捗回",
					TimeStart: sampleTime,
					TimeEnd:   sampleTime,
				},
				Place:       "S516",
				Level:       domain.EventLevelPrivate,
				HostName:    []*domain.User{{ID: sampleUUID}},
				Description: "第n回の進捗会です。",
				GroupID:     sampleUUID,
				RoomID:      sampleUUID,
			},
			setup: func(f mockEventRepositoryFields, args args, want *domain.EventDetail) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "level"}).
							AddRow(args.id, domain.EventLevelPrivate),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "LevelNotFound",
			args: args{
				id: sampleUUID,
			},
			want: &domain.EventDetail{
				Event: domain.Event{
					ID:        sampleUUID,
					Name:      "第n回進捗回",
					TimeStart: sampleTime,
					TimeEnd:   sampleTime,
				},
				Place:       "S516",
				Level:       domain.EventLevelAnonymous,
				HostName:    []*domain.User{{ID: sampleUUID}},
				Description: "第n回の進捗会です。",
				GroupID:     sampleUUID,
				RoomID:      sampleUUID,
			},
			setup: func(f mockEventRepositoryFields, args args, want *domain.EventDetail) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(repository.ErrNotFound)
			},
			assertion: assert.NoError,
		},
		{
			name: "KnoqNotFound",
			args: args{
				id: random.UUID(),
			},
			want:      nil,
			setup:     func(f mockEventRepositoryFields, args args, want *domain.EventDetail) {},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError",
			args: args{
				id: sampleUUID,
			},
			want: nil,
			setup: func(f mockEventRepositoryFields, args args, want *domain.EventDetail) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockEventRepositoryFields()
			tt.setup(f, tt.args, tt.want)
			repo := NewEventRepository(f.h, f.api)
			// Assertion
			got, err := repo.GetEvent(tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
