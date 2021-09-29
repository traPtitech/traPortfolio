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
	"gorm.io/gorm"
)

var (
	ids = []uuid.UUID{
		uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111"),
		uuid.FromStringOrNil("22222222-2222-2222-2222-222222222222"),
		uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333"),
	}
	testUsers = []*domain.User{
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
			want:   testUsers,
			setup: func(f fields, want []*domain.User) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.EXPECT().Find(&[]*model.User{}).
					DoAndReturn(func(users *[]*model.User) database.SQLHandler {
						for _, user := range want {
							*users = append(*users, &model.User{
								ID:   user.ID,
								Name: user.Name,
							})
						}
						return f.sqlhandler
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
				sqlhandler.EXPECT().Find(&[]*model.User{}).Return(f.sqlhandler)
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
