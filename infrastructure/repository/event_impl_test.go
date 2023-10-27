package repository

import (
	"context"
	"math/rand"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/infrastructure/external"
	"github.com/traPtitech/traPortfolio/infrastructure/external/mock_external"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
	"gorm.io/gorm"
)

type mockEventRepositoryFields struct {
	h    *MockSQLHandler
	knoq *mock_external.MockKnoqAPI
}

func newMockEventRepositoryFields(t *testing.T, ctrl *gomock.Controller) mockEventRepositoryFields {
	t.Helper()
	return mockEventRepositoryFields{
		h:    NewMockSQLHandler(),
		knoq: mock_external.NewMockKnoqAPI(ctrl),
	}
}

func TestEventRepository_GetEvents(t *testing.T) {
	var (
		since1, until1 = random.SinceAndUntil()
		since2, until2 = random.SinceAndUntil()
	)

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
					TimeStart: since1,
					TimeEnd:   until1,
				},
				{
					ID:        random.UUID(),
					Name:      random.AlphaNumeric(),
					TimeStart: since2,
					TimeEnd:   until2,
				},
			},
			setup: func(f mockEventRepositoryFields, want []*domain.Event) {
				f.knoq.EXPECT().GetEvents().Return(makeKnoqEvents(t, want), nil)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT `id` FROM `event_level_relations` WHERE level = ? AND id IN (?,?)")).
					WithArgs(domain.EventLevelPrivate, want[0].ID.String(), want[1].ID.String()).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "KnoqError",
			want: nil,
			setup: func(f mockEventRepositoryFields, want []*domain.Event) {
				f.knoq.EXPECT().GetEvents().Return(nil, errUnexpected)
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
			f := newMockEventRepositoryFields(t, ctrl)
			tt.setup(f, tt.want)
			repo := NewEventRepository(f.h.Conn, f.knoq)
			// Assertion
			got, err := repo.GetEvents(context.Background())
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEventRepository_GetEvent(t *testing.T) {
	since, until := random.SinceAndUntil()

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
					TimeStart: since,
					TimeEnd:   until,
				},
				Place:       random.AlphaNumeric(),
				Level:       domain.EventLevelPublic,
				HostName:    []*domain.User{{ID: random.UUID()}},
				Description: random.AlphaNumeric(),
				GroupID:     random.UUID(),
				RoomID:      random.UUID(),
			},
			setup: func(f mockEventRepositoryFields, args args, want *domain.EventDetail) {
				f.knoq.EXPECT().GetEvent(args.id).Return(makeKnoqEvent(t, want), nil)
				f.h.Mock.ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "level"}).
							AddRow(args.id, domain.EventLevelPublic),
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
				f.knoq.EXPECT().GetEvent(args.id).Return(nil, repository.ErrNotFound)
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
					TimeStart: since,
					TimeEnd:   until,
				},
				Place:       random.AlphaNumeric(),
				Level:       domain.EventLevelAnonymous,
				HostName:    nil,
				Description: random.AlphaNumeric(),
				GroupID:     random.UUID(),
				RoomID:      random.UUID(),
			},
			setup: func(f mockEventRepositoryFields, args args, want *domain.EventDetail) {
				f.knoq.EXPECT().GetEvent(args.id).Return(makeKnoqEvent(t, want), nil)
				f.h.Mock.ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(repository.ErrNotFound)
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
						TimeStart: since,
						TimeEnd:   until,
					},
					Place:       random.AlphaNumeric(),
					Level:       domain.EventLevelPrivate,
					HostName:    []*domain.User{{ID: random.UUID()}},
					Description: random.AlphaNumeric(),
					GroupID:     random.UUID(),
					RoomID:      random.UUID(),
				}
				f.knoq.EXPECT().GetEvent(args.id).Return(makeKnoqEvent(t, &ed), nil)
				f.h.Mock.ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
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
			f := newMockEventRepositoryFields(t, ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewEventRepository(f.h.Conn, f.knoq)
			// Assertion
			got, err := repo.GetEvent(context.Background(), tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEventRepository_CreateEventLevel(t *testing.T) {
	t.Parallel()

	type args struct {
		args *repository.CreateEventLevelArgs
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
				args: &repository.CreateEventLevelArgs{
					EventID: random.UUID(),
					Level:   domain.EventLevel(rand.Intn(int(domain.EventLevelLimit))),
				},
			},
			setup: func(f mockEventRepositoryFields, args args) {
				since, until := random.SinceAndUntil()
				event := external.EventResponse{
					ID:          args.args.EventID,
					Name:        random.AlphaNumeric(),
					Description: random.AlphaNumeric(),
					Place:       random.AlphaNumeric(),
					GroupID:     random.UUID(),
					RoomID:      random.UUID(),
					TimeStart:   since,
					TimeEnd:     until,
					SharedRoom:  random.Bool(),
				}
				f.knoq.EXPECT().GetEvent(args.args.EventID).Return(&event, nil)
				f.h.Mock.ExpectBegin()
				f.h.Mock.ExpectExec(makeSQLQueryRegexp("INSERT INTO `event_level_relations` (`id`,`level`,`created_at`,`updated_at`) VALUES (?,?,?,?)")).
					WithArgs(args.args.EventID, args.args.Level, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "KnoqNotFound",
			args: args{
				args: &repository.CreateEventLevelArgs{
					EventID: random.UUID(),
					Level:   domain.EventLevel(rand.Intn(int(domain.EventLevelLimit))),
				},
			},
			setup: func(f mockEventRepositoryFields, args args) {
				f.knoq.EXPECT().GetEvent(args.args.EventID).Return(nil, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
		{
			name: "LevelAlreadyExist",
			args: args{
				args: &repository.CreateEventLevelArgs{
					EventID: random.UUID(),
					Level:   domain.EventLevel(rand.Intn(int(domain.EventLevelLimit))),
				},
			},
			setup: func(f mockEventRepositoryFields, args args) {
				since, until := random.SinceAndUntil()
				event := external.EventResponse{
					ID:          args.args.EventID,
					Name:        random.AlphaNumeric(),
					Description: random.AlphaNumeric(),
					Place:       random.AlphaNumeric(),
					GroupID:     random.UUID(),
					RoomID:      random.UUID(),
					TimeStart:   since,
					TimeEnd:     until,
					SharedRoom:  random.Bool(),
				}
				f.knoq.EXPECT().GetEvent(args.args.EventID).Return(&event, nil)
				f.h.Mock.ExpectBegin()
				f.h.Mock.ExpectExec(makeSQLQueryRegexp("INSERT INTO `event_level_relations` (`id`,`level`,`created_at`,`updated_at`) VALUES (?,?,?,?)")).
					WithArgs(args.args.EventID, args.args.Level, anyTime{}, anyTime{}).
					WillReturnError(gorm.ErrRegistered)
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
			f := newMockEventRepositoryFields(t, ctrl)
			tt.setup(f, tt.args)
			repo := NewEventRepository(f.h.Conn, f.knoq)
			// Assertion
			err := repo.CreateEventLevel(context.Background(), tt.args.args)
			tt.assertion(t, err)
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
					Level: optional.From(domain.EventLevelPublic),
				},
			},
			setup: func(f mockEventRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "level"}).
							AddRow(args.id, domain.EventLevelAnonymous),
					)
				f.h.Mock.ExpectExec(makeSQLQueryRegexp("UPDATE `event_level_relations` SET `level`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.arg.Level.ValueOrZero(), anyTime{}, args.id).
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
					Level: optional.From(domain.EventLevelPublic),
				},
			},
			setup: func(f mockEventRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(repository.ErrNotFound)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "DoNotUpdate",
			args: args{
				id: random.UUID(),
				arg: &repository.UpdateEventLevelArgs{
					Level: optional.From(domain.EventLevelPublic),
				},
			},
			setup: func(f mockEventRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
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
					Level: optional.From(domain.EventLevelPublic),
				},
			},
			setup: func(f mockEventRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `event_level_relations` WHERE `event_level_relations`.`id` = ? ORDER BY `event_level_relations`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "level"}).
							AddRow(args.id, domain.EventLevelAnonymous),
					)
				f.h.Mock.ExpectExec(makeSQLQueryRegexp("UPDATE `event_level_relations` SET `level`=?,`updated_at`=? WHERE `id` = ?")).
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
			f := newMockEventRepositoryFields(t, ctrl)
			tt.setup(f, tt.args)
			repo := NewEventRepository(f.h.Conn, f.knoq)
			// Assertion
			tt.assertion(t, repo.UpdateEventLevel(context.Background(), tt.args.id, tt.args.arg))
		})
	}
}

func TestEventRepository_GetUserEvents(t *testing.T) {
	since1, until1 := random.SinceAndUntil()
	since2, until2 := random.SinceAndUntil()
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
					TimeStart: since1,
					TimeEnd:   until1,
				},
				{
					ID:        random.UUID(),
					Name:      random.AlphaNumeric(),
					TimeStart: since2,
					TimeEnd:   until2,
				},
			},
			setup: func(f mockEventRepositoryFields, args args, want []*domain.Event) {
				f.knoq.EXPECT().GetEventsByUserID(args.userID).Return(makeKnoqEvents(t, want), nil)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT `id` FROM `event_level_relations` WHERE level = ? AND id IN (?,?)")).
					WithArgs(domain.EventLevelPrivate, want[0].ID.String(), want[1].ID.String()).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}),
					)
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
				f.knoq.EXPECT().GetEventsByUserID(args.userID).Return(nil, errUnexpected)
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
			f := newMockEventRepositoryFields(t, ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewEventRepository(f.h.Conn, f.knoq)
			// Assertion
			got, err := repo.GetUserEvents(context.Background(), tt.args.userID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
