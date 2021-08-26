package service

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util"
)

func TestEventService_GetEvents(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.Event
		setup     func(m *MockRepository, args args, want []*domain.Event)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
			},
			want: []*domain.Event{
				{
					ID:        util.UUID(),
					Name:      util.AlphaNumeric(5),
					TimeStart: time.Now(),
					TimeEnd:   time.Now(),
				},
			},
			setup: func(m *MockRepository, args args, want []*domain.Event) {
				m.event.EXPECT().GetEvents().Return(want, nil)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewEventService(repo.event)
			got, err := s.GetEvents(tt.args.ctx)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEventService_GetEventByID(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.EventDetail
		setup     func(m *MockRepository, args args, want *domain.EventDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  util.UUID(),
			},
			want: &domain.EventDetail{
				Event: domain.Event{
					ID:        util.UUID(),
					Name:      util.AlphaNumeric(5),
					TimeStart: time.Now(),
					TimeEnd:   time.Now(),
				},
				Description: util.AlphaNumeric(10),
				Place:       util.AlphaNumeric(5),
				Level:       domain.EventLevelAnonymous,
				HostName: []*domain.User{
					{
						ID:       util.UUID(),
						Name:     util.AlphaNumeric(5),
						RealName: util.AlphaNumeric(5),
					},
				},
				GroupID: util.UUID(),
				RoomID:  util.UUID(),
			},
			setup: func(m *MockRepository, args args, want *domain.EventDetail) {
				m.event.EXPECT().GetEvent(args.id).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NilID",
			args: args{
				ctx: context.Background(),
				id:  uuid.Nil,
			},
			want: nil,
			setup: func(m *MockRepository, args args, want *domain.EventDetail) {
				m.event.EXPECT().GetEvent(args.id).Return(nil, repository.ErrInvalidID)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewEventService(repo.event)
			got, err := s.GetEventByID(tt.args.ctx, tt.args.id)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEventService_UpdateEvent(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		id  uuid.UUID
		arg *repository.UpdateEventArg
	}
	tests := []struct {
		name      string
		args      args
		setup     func(m *MockRepository, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  util.UUID(),
				arg: &repository.UpdateEventArg{
					Level: domain.EventLevelAnonymous,
				},
			},
			setup: func(m *MockRepository, args args) {
				m.event.EXPECT().UpdateEvent(args.id, args.arg).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NilID",
			args: args{
				ctx: context.Background(),
				id:  uuid.Nil,
				arg: &repository.UpdateEventArg{
					Level: domain.EventLevelAnonymous,
				},
			},
			setup: func(m *MockRepository, args args) {
				m.event.EXPECT().UpdateEvent(args.id, args.arg).Return(repository.ErrInvalidID)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args)

			s := NewEventService(repo.event)

			tt.assertion(t, s.UpdateEvent(tt.args.ctx, tt.args.id, tt.args.arg))
		})
	}
}
