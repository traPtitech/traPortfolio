package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database/mock_database"
	"github.com/traPtitech/traPortfolio/util/random"
)

type mockGroupRepositoryFields struct {
	h *mock_database.MockSQLHandler
}

func newMockGroupRepositoryFields() mockGroupRepositoryFields {
	return mockGroupRepositoryFields{
		h: mock_database.NewMockSQLHandler(),
	}
}

func TestGroupRepository_GetAllGroups(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		want      []*domain.Group
		setup     func(f mockGroupRepositoryFields, want []*domain.Group)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			want: []*domain.Group{
				{
					ID:   random.UUID(),
					Name: random.AlphaNumeric(),
				},
			},
			setup: func(f mockGroupRepositoryFields, want []*domain.Group) {
				rows := sqlmock.NewRows([]string{"group_id", "name", "link", "description", "created_at", "updated_at"})
				for _, v := range want {
					vlink := random.RandURLString()
					vdesc := random.AlphaNumeric()
					vtime := random.Time()
					rows.AddRow(v.ID, v.Name, vlink, vdesc, vtime, vtime)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `groups`")).
					WillReturnRows(rows)
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := newMockGroupRepositoryFields()
			tt.setup(f, tt.want)
			repo := NewGroupRepository(f.h)
			got, err := repo.GetAllGroups()
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
