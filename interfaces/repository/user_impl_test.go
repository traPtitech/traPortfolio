package repository

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/database/mock_database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/util"
	"gorm.io/gorm"
)

var (
	ids = []uuid.UUID{
		uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111"),
		uuid.FromStringOrNil("22222222-2222-2222-2222-222222222222"),
		uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333"),
	}
)

func TestUserRepository_GetUsers(t *testing.T) {
	t.Parallel()
	type fields struct {
		sqlhandler database.SQLHandler
		portal     external.PortalAPI
		traq       external.TraQAPI
	}
	tests := []struct {
		name      string
		fields    fields
		want      []*domain.User
		setup     func(f fields, want []*domain.User)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:   "Success",
			fields: fields{},
			want: []*domain.User{
				{
					ID:       ids[0],
					Name:     "user1",
					RealName: "ユーザー1 ユーザー1",
				},
				{
					ID:       ids[1],
					Name:     "user2",
					RealName: "ユーザー2 ユーザー2",
				},
				{
					ID:       ids[2],
					Name:     "lolico",
					RealName: "東 工子",
				},
			},
			setup: func(f fields, want []*domain.User) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.EXPECT().Find(&[]*model.User{}).
					DoAndReturn(func(users *[]*model.User) database.SQLHandler {
						_users := make([]*model.User, 0, len(want))
						for _, u := range want {
							_users = append(_users, &model.User{
								ID:   u.ID,
								Name: u.Name,
							})
						}
						*users = _users
						return sqlhandler
					})
				sqlhandler.EXPECT().Error().Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name:   "Fail_DB",
			fields: fields{},
			want:   nil,
			setup: func(f fields, want []*domain.User) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.EXPECT().Find(&[]*model.User{}).Return(sqlhandler)
				sqlhandler.EXPECT().Error().Return(gorm.ErrInvalidDB)
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
			tt.fields = fields{
				sqlhandler: mock_database.NewMockSQLHandler(ctrl),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			}
			tt.setup(tt.fields, tt.want)
			repo := NewUserRepository(tt.fields.sqlhandler, tt.fields.portal, tt.fields.traq)
			// Assertion
			got, err := repo.GetUsers()
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	t.Parallel()
	type fields struct {
		sqlhandler database.SQLHandler
		portal     external.PortalAPI
		traq       external.TraQAPI
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *domain.UserDetail
		setup     func(f fields, args args, want *domain.UserDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:   "Success",
			fields: fields{},
			args:   args{ids[0]},
			want: &domain.UserDetail{
				User: domain.User{
					ID:       ids[0],
					Name:     "user1",
					RealName: "ユーザー1 ユーザー1",
				},
				State: 1,
				Bio:   util.AlphaNumeric(5),
				Accounts: []*domain.Account{
					{
						ID:          util.UUID(),
						Type:        domain.HOMEPAGE,
						PrPermitted: true,
					},
				},
			},
			setup: func(f fields, args args, want *domain.UserDetail) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.EXPECT().Preload("Accounts").Return(sqlhandler)
				sqlhandler.EXPECT().First(&model.User{ID: args.id}).DoAndReturn(func(user *model.User) database.SQLHandler {
					user.ID = args.id
					user.Name = want.Name
					user.Description = want.Bio

					for _, v := range want.Accounts {
						user.Accounts = append(user.Accounts, &model.Account{
							ID:    v.ID,
							Type:  v.Type,
							Check: v.PrPermitted,
						})
					}

					return sqlhandler
				})
				sqlhandler.EXPECT().Error().Return(nil)
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				sqlhandler: mock_database.NewMockSQLHandler(ctrl),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			}
			tt.setup(tt.fields, tt.args, tt.want)
			repo := NewUserRepository(tt.fields.sqlhandler, tt.fields.portal, tt.fields.traq)
			// Assertion
			got, err := repo.GetUser(tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
