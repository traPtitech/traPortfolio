package repository

import (
	"fmt"
	"net/http"
	"testing"

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

func TestUserRepository_GetUsers(t *testing.T) {
	t.Parallel()
	type fields struct {
		sqlhandler *mock_database.MockSQLHandler
		portal     *mock_external.MockPortalAPI
		traq       *mock_external.MockTraQAPI
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
					ID:       util.UUID(),
					Name:     util.AlphaNumeric(5),
					RealName: util.AlphaNumeric(5),
				},
			},
			setup: func(f fields, want []*domain.User) {
				f.sqlhandler.EXPECT().Find(&[]*model.User{}).
					DoAndReturn(func(users *[]*model.User) database.SQLHandler {
						for _, user := range want {
							*users = append(*users, &model.User{
								ID:   user.ID,
								Name: user.Name,
							})
						}
						return f.sqlhandler
					})
				f.sqlhandler.EXPECT().Error().Return(nil)
				f.portal.EXPECT().GetAll().DoAndReturn(func() ([]*external.PortalUserResponse, error) {
					p := make([]*external.PortalUserResponse, 0, len(want))
					for _, v := range want {
						p = append(p, &external.PortalUserResponse{
							TraQID:   v.Name,
							RealName: v.RealName,
						})
					}
					return p, nil
				})
			},
			assertion: assert.NoError,
		},
		{
			name:   "Fail_DB",
			fields: fields{},
			want:   nil,
			setup: func(f fields, want []*domain.User) {
				f.sqlhandler.EXPECT().Find(&[]*model.User{}).Return(f.sqlhandler)
				f.sqlhandler.EXPECT().Error().Return(gorm.ErrInvalidDB)
			},
			assertion: assert.Error,
		},
		{
			name:   "Fail_Portal",
			fields: fields{},
			want:   nil,
			setup: func(f fields, want []*domain.User) {
				f.sqlhandler.EXPECT().Find(&[]*model.User{}).Return(f.sqlhandler)
				f.sqlhandler.EXPECT().Error().Return(nil)
				f.portal.EXPECT().GetAll().Return(nil, fmt.Errorf("GET /user failed: %v", http.StatusInternalServerError))
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
				portal:     mock_external.NewMockPortalAPI(ctrl),
				traq:       mock_external.NewMockTraQAPI(ctrl),
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
