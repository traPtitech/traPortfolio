package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/database/mock_database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/random"
)

type mockContestRepositoryFields struct {
	h      database.SQLHandler
	portal external.PortalAPI
}

func newMockContestRepositoryFields() mockContestRepositoryFields {
	return mockContestRepositoryFields{
		h:      mock_database.NewMockSQLHandler(),
		portal: mock_external.NewMockPortalAPI(),
	}
}

func TestContestRepository_GetContests(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		want      []*domain.Contest
		setup     func(f mockContestRepositoryFields, want []*domain.Contest)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			want: []*domain.Contest{
				{
					ID:        random.UUID(),
					Name:      random.AlphaNumeric(5),
					TimeStart: sampleTime,
					TimeEnd:   sampleTime,
				},
			},
			setup: func(f mockContestRepositoryFields, want []*domain.Contest) {
				rows := sqlmock.NewRows([]string{"id", "name", "since", "until"})
				for _, v := range want {
					rows.AddRow(v.ID, v.Name, v.TimeStart, v.TimeEnd)
				}
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests`")).
					WillReturnRows(rows)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			want: nil,
			setup: func(f mockContestRepositoryFields, want []*domain.Contest) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests`")).
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
			f := newMockContestRepositoryFields()
			tt.setup(f, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.GetContests()
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestRepository_GetContest(t *testing.T) {
	cid := random.UUID() // Successで使うcontestID

	t.Parallel()
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.ContestDetail
		setup     func(f mockContestRepositoryFields, args args, want *domain.ContestDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				id: cid,
			},
			want: &domain.ContestDetail{
				Contest: domain.Contest{
					ID:        cid,
					Name:      random.AlphaNumeric(5),
					TimeStart: sampleTime,
					TimeEnd:   sampleTime,
				},
				Link:        random.AlphaNumeric(5),
				Description: random.AlphaNumeric(10),
				// Teams:
			},
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestDetail) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "name", "since", "until", "link", "description"}).
							AddRow(args.id, want.Contest.Name, want.Contest.TimeStart, want.Contest.TimeEnd, want.Link, want.Description),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				id: random.UUID(),
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestDetail) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
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
			f := newMockContestRepositoryFields()
			tt.setup(f, tt.args, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.GetContest(tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestRepository_CreateContest(t *testing.T) {
	cname := random.AlphaNumeric(5) // Successで使用するContest.Name

	t.Parallel()
	type args struct {
		args *repository.CreateContestArgs
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.Contest
		setup     func(f mockContestRepositoryFields, args args, want *domain.Contest)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				args: &repository.CreateContestArgs{
					Name:        cname,
					Description: random.AlphaNumeric(10),
					Link:        random.AlphaNumeric(5),
					Since:       sampleTime,
					Until:       sampleTime,
				},
			},
			want: &domain.Contest{
				// ID: 比較しない
				Name:      cname,
				TimeStart: sampleTime,
				TimeEnd:   sampleTime,
			},
			setup: func(f mockContestRepositoryFields, args args, want *domain.Contest) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `contests` (`id`,`name`,`description`,`link`,`since`,`until`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Name, args.args.Description, args.args.Link, args.args.Since, args.args.Until, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				args: &repository.CreateContestArgs{
					Name:        random.AlphaNumeric(5),
					Description: random.AlphaNumeric(10),
					Link:        random.AlphaNumeric(5),
					Since:       sampleTime,
					Until:       sampleTime,
				},
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want *domain.Contest) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `contests` (`id`,`name`,`description`,`link`,`since`,`until`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Name, args.args.Description, args.args.Link, args.args.Since, args.args.Until, anyTime{}, anyTime{}).
					WillReturnError(errUnexpected)
				h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockContestRepositoryFields()
			tt.setup(f, tt.args, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.CreateContest(tt.args.args)
			tt.assertion(t, err)
			if got == nil {
				assert.Equal(t, tt.want, got)
			} else {
				// IDは比較しない
				assert.Equal(t, tt.want.Name, got.Name)
				assert.Equal(t, tt.want.TimeStart, got.TimeStart)
				assert.Equal(t, tt.want.TimeEnd, got.TimeEnd)
			}
		})
	}
}

func TestContestRepository_UpdateContest(t *testing.T) {
	t.Parallel()
	type args struct {
		id      uuid.UUID
		changes map[string]interface{}
	}
	tests := []struct {
		name      string
		args      args
		setup     func(f mockContestRepositoryFields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				id: random.UUID(),
				changes: map[string]interface{}{
					"name":        random.AlphaNumeric(5),
					"description": random.AlphaNumeric(10),
					"link":        random.AlphaNumeric(5),
					"since":       sampleTime,
					"until":       sampleTime,
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until", "created_at", "updated_at"}).
							AddRow(args.id, "", "", "", time.Time{}, time.Time{}, time.Time{}, time.Time{}),
					)
				h.Mock.
					ExpectExec(regexp.QuoteMeta("UPDATE `contests` SET `description`=?,`link`=?,`name`=?,`since`=?,`until`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.changes["description"], args.changes["link"], args.changes["name"], args.changes["since"], args.changes["until"], anyTime{}, args.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until", "created_at", "updated_at"}).
							AddRow(args.id, args.changes["name"], args.changes["description"], args.changes["link"], args.changes["since"], args.changes["until"], time.Time{}, time.Time{}),
					)
				h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				id: random.UUID(),
				changes: map[string]interface{}{
					"name":        random.AlphaNumeric(5),
					"description": random.AlphaNumeric(10),
					"link":        random.AlphaNumeric(5),
					"since":       sampleTime,
					"until":       sampleTime,
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(repository.ErrNotFound)
				h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError",
			args: args{
				id: random.UUID(),
				changes: map[string]interface{}{
					"name":        random.AlphaNumeric(5),
					"description": random.AlphaNumeric(10),
					"link":        random.AlphaNumeric(5),
					"since":       sampleTime,
					"until":       sampleTime,
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until", "created_at", "updated_at"}).
							AddRow(args.id, "", "", "", time.Time{}, time.Time{}, time.Time{}, time.Time{}),
					)
				h.Mock.
					ExpectExec(regexp.QuoteMeta("UPDATE `contests` SET `description`=?,`link`=?,`name`=?,`since`=?,`until`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.changes["description"], args.changes["link"], args.changes["name"], args.changes["since"], args.changes["until"], anyTime{}, args.id).
					WillReturnError(errUnexpected)
				h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockContestRepositoryFields()
			tt.setup(f, tt.args)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.UpdateContest(tt.args.id, tt.args.changes))
		})
	}
}

func TestContestRepository_DeleteContest(t *testing.T) {
	t.Parallel()
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		setup     func(f mockContestRepositoryFields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				id: random.UUID(),
			},
			setup: func(f mockContestRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until", "created_at", "updated_at"}).
							AddRow(args.id, "", "", "", time.Time{}, time.Time{}, time.Time{}, time.Time{}),
					)
				h.Mock.
					ExpectExec(regexp.QuoteMeta("DELETE FROM `contests` WHERE `contests`.`id` = ?")).
					WithArgs(args.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				id: random.UUID(),
			},
			setup: func(f mockContestRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(repository.ErrNotFound)
				h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError",
			args: args{
				id: random.UUID(),
			},
			setup: func(f mockContestRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until", "created_at", "updated_at"}).
							AddRow(args.id, "", "", "", time.Time{}, time.Time{}, time.Time{}, time.Time{}),
					)
				h.Mock.
					ExpectExec(regexp.QuoteMeta("DELETE FROM `contests` WHERE `contests`.`id` = ?")).
					WithArgs(args.id).
					WillReturnError(errUnexpected)
				h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockContestRepositoryFields()
			tt.setup(f, tt.args)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.DeleteContest(tt.args.id))
		})
	}
}

func TestContestRepository_GetContestTeams(t *testing.T) {
	t.Parallel()
	type args struct {
		contestID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.ContestTeam
		setup     func(f mockContestRepositoryFields, args args, want []*domain.ContestTeam)
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockContestRepositoryFields()
			tt.setup(f, tt.args, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.GetContestTeams(tt.args.contestID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestRepository_GetContestTeam(t *testing.T) {
	t.Parallel()
	type args struct {
		contestID uuid.UUID
		teamID    uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.ContestTeamDetail
		setup     func(f mockContestRepositoryFields, args args, want *domain.ContestTeamDetail)
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockContestRepositoryFields()
			tt.setup(f, tt.args, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.GetContestTeam(tt.args.contestID, tt.args.teamID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestRepository_CreateContestTeam(t *testing.T) {
	t.Parallel()
	type args struct {
		contestID    uuid.UUID
		_contestTeam *repository.CreateContestTeamArgs
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.ContestTeamDetail
		setup     func(f mockContestRepositoryFields, args args, want *domain.ContestTeamDetail)
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockContestRepositoryFields()
			tt.setup(f, tt.args, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.CreateContestTeam(tt.args.contestID, tt.args._contestTeam)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestRepository_UpdateContestTeam(t *testing.T) {
	t.Parallel()
	type args struct {
		teamID  uuid.UUID
		changes map[string]interface{}
	}
	tests := []struct {
		name      string
		args      args
		setup     func(f mockContestRepositoryFields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockContestRepositoryFields()
			tt.setup(f, tt.args)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.UpdateContestTeam(tt.args.teamID, tt.args.changes))
		})
	}
}

func TestContestRepository_GetContestTeamMembers(t *testing.T) {
	t.Parallel()
	type args struct {
		contestID uuid.UUID
		teamID    uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.User
		setup     func(f mockContestRepositoryFields, args args, want []*domain.User)
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockContestRepositoryFields()
			tt.setup(f, tt.args, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.GetContestTeamMembers(tt.args.contestID, tt.args.teamID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestRepository_AddContestTeamMembers(t *testing.T) {
	t.Parallel()
	type args struct {
		teamID  uuid.UUID
		members []uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		setup     func(f mockContestRepositoryFields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockContestRepositoryFields()
			tt.setup(f, tt.args)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.AddContestTeamMembers(tt.args.teamID, tt.args.members))
		})
	}
}

func TestContestRepository_DeleteContestTeamMembers(t *testing.T) {
	t.Parallel()
	type args struct {
		teamID  uuid.UUID
		members []uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		setup     func(f mockContestRepositoryFields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockContestRepositoryFields()
			tt.setup(f, tt.args)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.DeleteContestTeamMembers(tt.args.teamID, tt.args.members))
		})
	}
}
