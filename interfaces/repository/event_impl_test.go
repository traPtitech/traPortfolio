package repository_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database/mock_database"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external"
	impl "github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/random"
	"gorm.io/gorm"
)

type mockEventRepositoryFields struct {
	h    *mock_database.MockSQLHandler
	knoq *mock_external.MockKnoqAPI
}

func newMockEventRepositoryFields(ctrl *gomock.Controller) mockEventRepositoryFields {
	return mockEventRepositoryFields{
		h:    mock_database.NewMockSQLHandler(),
		knoq: mock_external.NewMockKnoqAPI(ctrl),
	}
}

func TestEventRepository_GetEvents(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		want      []*domain.Event
		setup     func(f mockEventRepositoryFields, want []*domain.Event)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			want: []*domain.Event{
				{
					ID:        random.UUID(),
					Name:      random.AlphaNumeric(),
					TimeStart: random.Time(),
					TimeEnd:   random.Time(),
				},
				{
					ID:        random.UUID(),
					Name:      random.AlphaNumeric(),
					TimeStart: random.Time(),
					TimeEnd:   random.Time(),
				},
			},
			setup: func(f mockEventRepositoryFields, want []*domain.Event) {
				f.knoq.EXPECT().GetAll().Return(makeKnoqEvents(want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "KnoqError",
			want: nil,
			setup: func(f mockEventRepositoryFields, want []*domain.Event) {
				f.knoq.EXPECT().GetAll().Return(nil, errUnexpected)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			f := newMockEventRepositoryFields(ctrl)
			tt.setup(f, tt.want)
			repo := impl.NewEventRepository(f.h, f.knoq)
			// Assertion
			got, err := repo.GetEvents()
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
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
				id: random.UUID(),
			},
			want: &domain.EventDetail{
				Event: domain.Event{
					ID:        random.UUID(),
					Name:      random.AlphaNumeric(),
					TimeStart: random.Time(),
					TimeEnd:   random.Time(),
				},
				Place:       random.AlphaNumeric(),
				Level:       domain.EventLevelPrivate,
				HostName:    []*domain.User{{ID: random.UUID()}},
				Description: random.AlphaNumeric(),
				GroupID:     random.UUID(),
				RoomID:      random.UUID(),
			},
			setup: func(f mockEventRepositoryFields, args args, want *domain.EventDetail) {
				f.knoq.EXPECT().GetByEventID(args.id).Return(makeKnoqEvent(want), nil)
				f.h.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "level"}).
							AddRow(args.id, domain.EventLevelPrivate),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "KnoqNotFound",
			args: args{
				id: random.UUID(),
			},
			want: nil,
			setup: func(f mockEventRepositoryFields, args args, want *domain.EventDetail) {
				f.knoq.EXPECT().GetByEventID(args.id).Return(nil, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
		{
			name: "LevelNotFound",
			args: args{
				id: random.UUID(),
			},
			want: &domain.EventDetail{
				Event: domain.Event{
					ID:        random.UUID(),
					Name:      random.AlphaNumeric(),
					TimeStart: random.Time(),
					TimeEnd:   random.Time(),
				},
				Place:       random.AlphaNumeric(),
				Level:       domain.EventLevelAnonymous,
				HostName:    []*domain.User{{ID: random.UUID()}},
				Description: random.AlphaNumeric(),
				GroupID:     random.UUID(),
				RoomID:      random.UUID(),
			},
			setup: func(f mockEventRepositoryFields, args args, want *domain.EventDetail) {
				f.knoq.EXPECT().GetByEventID(args.id).Return(makeKnoqEvent(want), nil)
				f.h.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError_getEventLevelByID",
			args: args{
				id: random.UUID(),
			},
			want: nil,
			setup: func(f mockEventRepositoryFields, args args, want *domain.EventDetail) {
				ed := domain.EventDetail{
					Event: domain.Event{
						ID:        random.UUID(),
						Name:      random.AlphaNumeric(),
						TimeStart: random.Time(),
						TimeEnd:   random.Time(),
					},
					Place:       random.AlphaNumeric(),
					Level:       domain.EventLevelPrivate,
					HostName:    []*domain.User{{ID: random.UUID()}},
					Description: random.AlphaNumeric(),
					GroupID:     random.UUID(),
					RoomID:      random.UUID(),
				}
				f.knoq.EXPECT().GetByEventID(args.id).Return(makeKnoqEvent(&ed), nil)
				f.h.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
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
			ctrl := gomock.NewController(t)
			f := newMockEventRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewEventRepository(f.h, f.knoq)
			// Assertion
			got, err := repo.GetEvent(tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEventRepository_UpdateEventLevel(t *testing.T) {
	t.Parallel()
	type args struct {
		id  uuid.UUID
		arg *repository.UpdateEventLevelArgs
	}
	tests := []struct {
		name      string
		args      args
		setup     func(f mockEventRepositoryFields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				id: random.UUID(),
				arg: &repository.UpdateEventLevelArgs{
					Level: domain.EventLevelPublic,
				},
			},
			setup: func(f mockEventRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "level"}).
							AddRow(args.id, domain.EventLevelAnonymous),
					)
				f.h.Mock.ExpectExec(regexp.QuoteMeta("UPDATE `event_level_relations` SET `level`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.arg.Level, anyTime{}, args.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "LevelNotFound",
			args: args{
				id: random.UUID(),
				arg: &repository.UpdateEventLevelArgs{
					Level: domain.EventLevelPublic,
				},
			},
			setup: func(f mockEventRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(gorm.ErrRecordNotFound)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "DoNotUpdate",
			args: args{
				id: random.UUID(),
				arg: &repository.UpdateEventLevelArgs{
					Level: domain.EventLevelPublic,
				},
			},
			setup: func(f mockEventRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "level"}).
							AddRow(args.id, domain.EventLevelPublic), // equal to args.arg.Level
					)
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "UpdateError",
			args: args{
				id: random.UUID(),
				arg: &repository.UpdateEventLevelArgs{
					Level: domain.EventLevelPublic,
				},
			},
			setup: func(f mockEventRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "level"}).
							AddRow(args.id, domain.EventLevelAnonymous),
					)
				f.h.Mock.ExpectExec(regexp.QuoteMeta("UPDATE `event_level_relations` SET `level`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.arg.Level, anyTime{}, args.id).
					WillReturnError(errUnexpected)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			f := newMockEventRepositoryFields(ctrl)
			tt.setup(f, tt.args)
			repo := impl.NewEventRepository(f.h, f.knoq)
			// Assertion
			tt.assertion(t, repo.UpdateEventLevel(tt.args.id, tt.args.arg))
		})
	}
}

func TestEventRepository_GetUserEvents(t *testing.T) {
	t.Parallel()
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.Event
		setup     func(f mockEventRepositoryFields, args args, want []*domain.Event)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				userID: random.UUID(),
			},
			want: []*domain.Event{
				{
					ID:        random.UUID(),
					Name:      random.AlphaNumeric(),
					TimeStart: random.Time(),
					TimeEnd:   random.Time(),
				},
				{
					ID:        random.UUID(),
					Name:      random.AlphaNumeric(),
					TimeStart: random.Time(),
					TimeEnd:   random.Time(),
				},
			},
			setup: func(f mockEventRepositoryFields, args args, want []*domain.Event) {
				f.knoq.EXPECT().GetByUserID(args.userID).Return(makeKnoqEvents(want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				userID: random.UUID(),
			},
			want: nil,
			setup: func(f mockEventRepositoryFields, args args, want []*domain.Event) {
				f.knoq.EXPECT().GetByUserID(args.userID).Return(nil, errUnexpected)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			f := newMockEventRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewEventRepository(f.h, f.knoq)
			// Assertion
			got, err := repo.GetUserEvents(tt.args.userID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
