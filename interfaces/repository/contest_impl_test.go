package repository

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database/mock_database"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
)

type mockContestRepositoryFields struct {
	h      *mock_database.MockSQLHandler
	portal *mock_external.MockPortalAPI
}

func newMockContestRepositoryFields(t *testing.T, ctrl *gomock.Controller) mockContestRepositoryFields {
	t.Helper()
	return mockContestRepositoryFields{
		h:      mock_database.NewMockSQLHandler(),
		portal: mock_external.NewMockPortalAPI(ctrl),
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
					Name:      random.AlphaNumeric(),
					TimeStart: sampleTime,
					TimeEnd:   sampleTime,
				},
			},
			setup: func(f mockContestRepositoryFields, want []*domain.Contest) {
				rows := sqlmock.NewRows([]string{"id", "name", "since", "until"})
				for _, v := range want {
					rows.AddRow(v.ID, v.Name, v.TimeStart, v.TimeEnd)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests`")).
					WillReturnRows(rows)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			want: nil,
			setup: func(f mockContestRepositoryFields, want []*domain.Contest) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests`")).
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
			f := newMockContestRepositoryFields(t, ctrl)
			tt.setup(f, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.GetContests(context.Background())
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
					Name:      random.AlphaNumeric(),
					TimeStart: sampleTime,
					TimeEnd:   sampleTime,
				},
				Link:        random.RandURLString(),
				Description: random.AlphaNumeric(),
				// Teams:
			},
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestDetail) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
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
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
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
			f := newMockContestRepositoryFields(t, ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.GetContest(context.Background(), tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestRepository_CreateContest(t *testing.T) {
	cname := random.AlphaNumeric()       // Successで使用するContest.Name
	link := random.AlphaNumeric()        // Successで使用するContestDetail.Link
	description := random.AlphaNumeric() // Successで使用するContestDetail.Description

	t.Parallel()
	type args struct {
		args *repository.CreateContestArgs
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
				args: &repository.CreateContestArgs{
					Name:        cname,
					Description: description,
					Link:        optional.From(link),
					Since:       sampleTime,
					Until:       optional.From(sampleTime),
				},
			},
			want: &domain.ContestDetail{
				Contest: domain.Contest{
					Name:      cname,
					TimeStart: sampleTime,
					TimeEnd:   sampleTime,
				},
				Link:        link,
				Description: description,
			},
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestDetail) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `contests` (`id`,`name`,`description`,`link`,`since`,`until`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Name, args.args.Description, args.args.Link.ValueOrZero(), args.args.Since, args.args.Until.ValueOrZero(), anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(anyUUID{}).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until"}).
							AddRow(uuid.Nil, args.args.Name, args.args.Description, args.args.Link.ValueOrZero(), args.args.Since, args.args.Until.ValueOrZero()),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError Create",
			args: args{
				args: &repository.CreateContestArgs{
					Name:        random.AlphaNumeric(),
					Description: random.AlphaNumeric(),
					Link:        random.Optional(random.RandURLString()),
					Since:       sampleTime,
					Until:       optional.From(sampleTime),
				},
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestDetail) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `contests` (`id`,`name`,`description`,`link`,`since`,`until`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Name, args.args.Description, args.args.Link.ValueOrZero(), args.args.Since, args.args.Until.ValueOrZero(), anyTime{}, anyTime{}).
					WillReturnError(errUnexpected)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError Get",
			args: args{
				args: &repository.CreateContestArgs{
					Name:        random.AlphaNumeric(),
					Description: random.AlphaNumeric(),
					Link:        random.Optional(random.RandURLString()),
					Since:       sampleTime,
					Until:       optional.From(sampleTime),
				},
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestDetail) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `contests` (`id`,`name`,`description`,`link`,`since`,`until`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Name, args.args.Description, args.args.Link.ValueOrZero(), args.args.Since, args.args.Until.ValueOrZero(), anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(anyUUID{}).
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
			f := newMockContestRepositoryFields(t, ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.CreateContest(context.Background(), tt.args.args)
			if tt.want != nil && got != nil {
				tt.want.ID = got.ID // 関数内でIDを生成するためここで合わせる
			}
			tt.assertion(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestRepository_UpdateContest(t *testing.T) {
	t.Parallel()
	type args struct {
		id   uuid.UUID
		args *repository.UpdateContestArgs
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
				args: &repository.UpdateContestArgs{
					Name:        optional.From(random.AlphaNumeric()),
					Description: optional.From(random.AlphaNumeric()),
					Link:        optional.From(random.RandURLString()),
					Since:       optional.From(sampleTime),
					Until:       optional.From(sampleTime),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until", "created_at", "updated_at"}).
							AddRow(args.id, "", "", "", time.Time{}, time.Time{}, time.Time{}, time.Time{}),
					)
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("UPDATE `contests` SET `description`=?,`link`=?,`name`=?,`since`=?,`until`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.Description.ValueOrZero(), args.args.Link.ValueOrZero(), args.args.Name.ValueOrZero(), args.args.Since.ValueOrZero(), args.args.Until.ValueOrZero(), anyTime{}, args.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until", "created_at", "updated_at"}).
							AddRow(args.id, args.args.Name.ValueOrZero(), args.args.Description.ValueOrZero(), args.args.Link.ValueOrZero(), args.args.Since.ValueOrZero(), args.args.Until.ValueOrZero(), time.Time{}, time.Time{}),
					)
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				id: random.UUID(),
				args: &repository.UpdateContestArgs{
					Name:        optional.From(random.AlphaNumeric()),
					Description: optional.From(random.AlphaNumeric()),
					Link:        optional.From(random.RandURLString()),
					Since:       optional.From(sampleTime),
					Until:       optional.From(sampleTime),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(repository.ErrNotFound)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_Update",
			args: args{
				id: random.UUID(),
				args: &repository.UpdateContestArgs{
					Name:        optional.From(random.AlphaNumeric()),
					Description: optional.From(random.AlphaNumeric()),
					Link:        optional.From(random.RandURLString()),
					Since:       optional.From(sampleTime),
					Until:       optional.From(sampleTime),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until", "created_at", "updated_at"}).
							AddRow(args.id, "", "", "", time.Time{}, time.Time{}, time.Time{}, time.Time{}),
					)
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("UPDATE `contests` SET `description`=?,`link`=?,`name`=?,`since`=?,`until`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.Description.ValueOrZero(), args.args.Link.ValueOrZero(), args.args.Name.ValueOrZero(), args.args.Since.ValueOrZero(), args.args.Until.ValueOrZero(), anyTime{}, args.id).
					WillReturnError(errUnexpected)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_CheckExist",
			args: args{
				id: random.UUID(),
				args: &repository.UpdateContestArgs{
					Name:        optional.From(random.AlphaNumeric()),
					Description: optional.From(random.AlphaNumeric()),
					Link:        optional.From(random.RandURLString()),
					Since:       optional.From(sampleTime),
					Until:       optional.From(sampleTime),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until", "created_at", "updated_at"}).
							AddRow(args.id, "", "", "", time.Time{}, time.Time{}, time.Time{}, time.Time{}),
					)
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("UPDATE `contests` SET `description`=?,`link`=?,`name`=?,`since`=?,`until`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.Description.ValueOrZero(), args.args.Link.ValueOrZero(), args.args.Name.ValueOrZero(), args.args.Since.ValueOrZero(), args.args.Until.ValueOrZero(), anyTime{}, args.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
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
			f := newMockContestRepositoryFields(t, ctrl)
			tt.setup(f, tt.args)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.UpdateContest(context.Background(), tt.args.id, tt.args.args))
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
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until", "created_at", "updated_at"}).
							AddRow(args.id, "", "", "", time.Time{}, time.Time{}, time.Time{}, time.Time{}),
					)
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("DELETE FROM `contests` WHERE `contests`.`id` = ?")).
					WithArgs(args.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				id: random.UUID(),
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(repository.ErrNotFound)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError",
			args: args{
				id: random.UUID(),
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until", "created_at", "updated_at"}).
							AddRow(args.id, "", "", "", time.Time{}, time.Time{}, time.Time{}, time.Time{}),
					)
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("DELETE FROM `contests` WHERE `contests`.`id` = ?")).
					WithArgs(args.id).
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
			f := newMockContestRepositoryFields(t, ctrl)
			tt.setup(f, tt.args)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.DeleteContest(context.Background(), tt.args.id))
		})
	}
}

func TestContestRepository_GetContestTeams(t *testing.T) {
	cid := random.UUID() // Successで使うcontestID

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
		{
			name: "Success",
			args: args{
				contestID: cid,
			},
			want: []*domain.ContestTeam{
				{
					ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
						ID:        random.UUID(),
						ContestID: cid,
						Name:      random.AlphaNumeric(),
						Result:    random.AlphaNumeric(),
					},
					Members: []*domain.User{
						domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
						domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
					},
				},
			},
			setup: func(f mockContestRepositoryFields, args args, want []*domain.ContestTeam) {
				rows := sqlmock.NewRows([]string{"id", "contest_id", "name", "result"})
				for _, v := range want {
					rows.AddRow(v.ID, v.ContestID, v.Name, v.Result)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.contestID))
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`contest_id` = ?")).
					WithArgs(args.contestID).
					WillReturnRows(rows)

				teamsIDList := make([]driver.Value, len(want))
				teamRows := sqlmock.NewRows([]string{"team_id", "user_id"})
				for i, u := range want {
					teamsIDList[i] = u.ID
					for _, w := range u.Members {
						teamRows.AddRow(u.ID, w.ID)
					}
				}

				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE team_id IN (?)")).
					WithArgs(teamsIDList...).
					WillReturnRows(teamRows)

				userIDs := make([]driver.Value, len(want[0].Members))
				userRows := sqlmock.NewRows([]string{"id", "name", "check"})
				for i, u := range want[0].Members {
					userIDs[i] = u.ID
					userRows.AddRow(u.ID, u.Name, u.Check)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?,?)")).
					WithArgs(userIDs...).
					WillReturnRows(userRows)

				f.portal.EXPECT().GetAll().Return(makePortalUsers(t, want[0].Members), nil)

			},
			assertion: assert.NoError,
		},
		{
			name: "ContestNotFound",
			args: args{
				contestID: random.UUID(),
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want []*domain.ContestTeam) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnError(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_FindContestTeam",
			args: args{
				contestID: random.UUID(),
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want []*domain.ContestTeam) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.contestID))
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`contest_id` = ?")).
					WithArgs(args.contestID).
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
			f := newMockContestRepositoryFields(t, ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.GetContestTeams(context.Background(), tt.args.contestID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestRepository_GetContestTeam(t *testing.T) {
	cid := random.UUID() // Successで使うcontestID
	tid := random.UUID() // Successで使うteamID

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
		{
			name: "Success",
			args: args{
				contestID: cid,
				teamID:    tid,
			},
			want: &domain.ContestTeamDetail{
				ContestTeam: domain.ContestTeam{
					ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
						ID:        tid,
						ContestID: cid,
						Name:      random.AlphaNumeric(),
						Result:    random.AlphaNumeric(),
					},
					Members: []*domain.User{
						domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
						domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
					},
				},
				Link:        random.RandURLString(),
				Description: random.AlphaNumeric(),
			},
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestTeamDetail) {
				u := want.Members[0]
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? AND `contest_teams`.`contest_id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID, args.contestID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id", "name", "result", "link", "description"}).
							AddRow(want.ContestTeam.ID, want.ContestTeam.ContestID, want.ContestTeam.Name, want.ContestTeam.Result, want.Link, want.Description),
					)

				belongingRows := sqlmock.NewRows([]string{"team_id", "user_id"})
				for _, u := range want.Members {
					belongingRows.AddRow(args.teamID, u.ID)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(belongingRows)

				userIDs := make([]driver.Value, len(want.Members))
				userRows := sqlmock.NewRows([]string{"id", "name", "check"})
				for i, u := range want.Members {
					userIDs[i] = u.ID
					userRows.AddRow(u.ID, u.Name, u.Check)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?,?)")).
					WithArgs(userIDs...).
					WillReturnRows(userRows)

				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ?")).
					WithArgs(u.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "check"}).
							AddRow(u.ID, u.Name, u.Check),
					)

				f.portal.EXPECT().GetAll().Return(makePortalUsers(t, want.Members), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestTeamDetail) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? AND `contest_teams`.`contest_id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID, args.contestID).
					WillReturnError(repository.ErrNotFound)
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
			f := newMockContestRepositoryFields(t, ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.GetContestTeam(context.Background(), tt.args.contestID, tt.args.teamID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestRepository_CreateContestTeam(t *testing.T) {
	cid := random.UUID() // Successで使うcontestID
	successArgs := repository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      optional.From(random.AlphaNumeric()),
		Link:        optional.From(random.RandURLString()),
		Description: random.AlphaNumeric(),
	}

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
		{
			name: "Success",
			args: args{
				contestID:    cid,
				_contestTeam: &successArgs,
			},
			want: &domain.ContestTeamDetail{
				ContestTeam: domain.ContestTeam{
					ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
						// ID: Assertion時にgot.IDと合わせる
						ContestID: cid,
						Name:      successArgs.Name,
						Result:    successArgs.Result.ValueOrZero(),
					},
					Members: make([]*domain.User, 0),
				},
				Link:        successArgs.Link.ValueOrZero(),
				Description: successArgs.Description,
			},
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestTeamDetail) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.contestID))
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `contest_teams` (`id`,`contest_id`,`name`,`description`,`result`,`link`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.contestID, args._contestTeam.Name, args._contestTeam.Description, args._contestTeam.Result, args._contestTeam.Link, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "ContestNotFound",
			args: args{
				contestID: cid,
				_contestTeam: &repository.CreateContestTeamArgs{
					Name:        random.AlphaNumeric(),
					Result:      random.Optional(random.AlphaNumeric()),
					Link:        random.Optional(random.RandURLString()),
					Description: random.AlphaNumeric(),
				},
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestTeamDetail) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnError(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
		{
			name: "CreateContestFailed",
			args: args{
				contestID: cid,
				_contestTeam: &repository.CreateContestTeamArgs{
					Name:        random.AlphaNumeric(),
					Result:      random.Optional(random.AlphaNumeric()),
					Link:        random.Optional(random.RandURLString()),
					Description: random.AlphaNumeric(),
				},
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestTeamDetail) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.contestID))
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `contest_teams` (`id`,`contest_id`,`name`,`description`,`result`,`link`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.contestID, args._contestTeam.Name, args._contestTeam.Description, args._contestTeam.Result, args._contestTeam.Link, anyTime{}, anyTime{}).
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
			f := newMockContestRepositoryFields(t, ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.CreateContestTeam(context.Background(), tt.args.contestID, tt.args._contestTeam)
			if tt.want != nil && got != nil {
				tt.want.ID = got.ID // 関数内でIDを生成するためここで合わせる
			}
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestRepository_UpdateContestTeam(t *testing.T) {
	t.Parallel()
	type args struct {
		teamID uuid.UUID
		args   *repository.UpdateContestTeamArgs
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
				teamID: random.UUID(),
				args: &repository.UpdateContestTeamArgs{
					Name:        optional.From(random.AlphaNumeric()),
					Description: optional.From(random.AlphaNumeric()),
					Link:        optional.From(random.RandURLString()),
					Result:      optional.From(random.AlphaNumeric()),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				cid := random.UUID()
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id", "name", "description", "result", "link", "created_at", "updated_at"}).
							AddRow(args.teamID, cid, "", "", "", "", time.Time{}, time.Time{}),
					)
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("UPDATE `contest_teams` SET `description`=?,`link`=?,`name`=?,`result`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.Description.ValueOrZero(), args.args.Link.ValueOrZero(), args.args.Name.ValueOrZero(), args.args.Result.ValueOrZero(), anyTime{}, args.teamID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id", "name", "description", "result", "link", "created_at", "updated_at"}).
							AddRow(args.teamID, cid, args.args.Name, args.args.Description, args.args.Result, args.args.Link, time.Time{}, time.Time{}),
					)
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				teamID: random.UUID(),
				args: &repository.UpdateContestTeamArgs{
					Name:        optional.From(random.AlphaNumeric()),
					Description: optional.From(random.AlphaNumeric()),
					Link:        optional.From(random.RandURLString()),
					Result:      optional.From(random.AlphaNumeric()),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnError(repository.ErrNotFound)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_Update",
			args: args{
				teamID: random.UUID(),
				args: &repository.UpdateContestTeamArgs{
					Name:        optional.From(random.AlphaNumeric()),
					Description: optional.From(random.AlphaNumeric()),
					Link:        optional.From(random.RandURLString()),
					Result:      optional.From(random.AlphaNumeric()),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id", "name", "description", "result", "link", "created_at", "updated_at"}).
							AddRow(args.teamID, random.UUID(), "", "", "", "", time.Time{}, time.Time{}),
					)
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("UPDATE `contest_teams` SET `description`=?,`link`=?,`name`=?,`result`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.Description.ValueOrZero(), args.args.Link.ValueOrZero(), args.args.Name.ValueOrZero(), args.args.Result.ValueOrZero(), anyTime{}, args.teamID).
					WillReturnError(errUnexpected)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_CheckExist",
			args: args{
				teamID: random.UUID(),
				args: &repository.UpdateContestTeamArgs{
					Name:        optional.From(random.AlphaNumeric()),
					Description: optional.From(random.AlphaNumeric()),
					Link:        optional.From(random.RandURLString()),
					Result:      optional.From(random.AlphaNumeric()),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id", "name", "description", "result", "link", "created_at", "updated_at"}).
							AddRow(args.teamID, random.UUID(), "", "", "", "", time.Time{}, time.Time{}),
					)
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("UPDATE `contest_teams` SET `description`=?,`link`=?,`name`=?,`result`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.Description.ValueOrZero(), args.args.Link.ValueOrZero(), args.args.Name.ValueOrZero(), args.args.Result.ValueOrZero(), anyTime{}, args.teamID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
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
			f := newMockContestRepositoryFields(t, ctrl)
			tt.setup(f, tt.args)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.UpdateContestTeam(context.Background(), tt.args.teamID, tt.args.args))
		})
	}
}

func TestContestRepository_DeleteContestTeam(t *testing.T) {
	t.Parallel()
	type args struct {
		contestID uuid.UUID
		teamID    uuid.UUID
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
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? AND `contest_teams`.`contest_id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID, args.contestID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, args.contestID),
					)
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("DELETE FROM `contest_teams` WHERE `contest_teams`.`id` = ?")).
					WithArgs(args.teamID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
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
			f := newMockContestRepositoryFields(t, ctrl)
			tt.setup(f, tt.args)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.DeleteContestTeam(context.Background(), tt.args.contestID, tt.args.teamID))
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
		{
			name: "Success_Single",
			args: args{
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			want: []*domain.User{
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
			},
			setup: func(f mockContestRepositoryFields, args args, want []*domain.User) {
				u := want[0]
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"team_id", "user_id"}).
							AddRow(args.teamID, u.ID),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ?")).
					WithArgs(u.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "check"}).
							AddRow(u.ID, u.Name, u.Check),
					)
				f.portal.EXPECT().GetAll().Return(makePortalUsers(t, want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_Multiple",
			args: args{
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			want: []*domain.User{
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
			},
			setup: func(f mockContestRepositoryFields, args args, want []*domain.User) {
				belongingRows := sqlmock.NewRows([]string{"team_id", "user_id"})
				for _, u := range want {
					belongingRows.AddRow(args.teamID, u.ID)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(belongingRows)
				userIDs := make([]driver.Value, len(want))
				userRows := sqlmock.NewRows([]string{"id", "name", "check"})
				for i, u := range want {
					userIDs[i] = u.ID
					userRows.AddRow(u.ID, u.Name, u.Check)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?,?)")).
					WithArgs(userIDs...).
					WillReturnRows(userRows)
				f.portal.EXPECT().GetAll().Return(makePortalUsers(t, want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ContestNotFound",
			args: args{
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want []*domain.User) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnError(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
		{
			name: "ContestTeamNotFound",
			args: args{
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want []*domain.User) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnError(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError",
			args: args{
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want []*domain.User) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "PortalError",
			args: args{
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want []*domain.User) {
				u := domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool())
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"team_id", "user_id"}).
							AddRow(args.teamID, u.ID),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ?")).
					WithArgs(u.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name"}).
							AddRow(u.ID, u.Name),
					)
				f.portal.EXPECT().GetAll().Return(nil, errUnexpected)
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
			f := newMockContestRepositoryFields(t, ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.GetContestTeamMembers(context.Background(), tt.args.contestID, tt.args.teamID)
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
		{
			name: "Success",
			args: args{
				teamID: random.UUID(),
				members: []uuid.UUID{
					random.UUID(), // 新たに追加するメンバー
					random.UUID(), // すでに存在するメンバー
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				rows := sqlmock.NewRows([]string{"team_id", "user_id"})
				newUsers := make([]uuid.UUID, 0, len(args.members))
				for i, u := range args.members {
					if i%2 == 0 {
						rows.AddRow(args.teamID, u)
					} else {
						newUsers = append(newUsers, u)
					}
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(rows)
				f.h.Mock.ExpectBegin()
				for _, u := range newUsers {
					f.h.Mock.
						ExpectExec(makeSQLQueryRegexp("INSERT INTO `contest_team_user_belongings` (`team_id`,`user_id`,`created_at`,`updated_at`) VALUES (?,?,?,?)")).
						WithArgs(args.teamID, u, anyTime{}, anyTime{}).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "ZeroMembers",
			args: args{
				teamID:  random.UUID(),
				members: []uuid.UUID{},
			},
			setup:     func(f mockContestRepositoryFields, args args) {},
			assertion: assert.Error,
		},
		{
			name: "ContestTeamNotFound",
			args: args{
				teamID:  random.UUID(),
				members: []uuid.UUID{random.UUID()},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnError(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_FindBelongings",
			args: args{
				teamID:  random.UUID(),
				members: []uuid.UUID{random.UUID()},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_CreateNewBelongings",
			args: args{
				teamID:  random.UUID(),
				members: []uuid.UUID{random.UUID()},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "team_id", "user_id"}).
							AddRow(random.UUID(), args.teamID, random.UUID()),
					)
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `contest_team_user_belongings` (`team_id`,`user_id`,`created_at`,`updated_at`) VALUES (?,?,?,?)")).
					WithArgs(args.teamID, args.members[0], anyTime{}, anyTime{}).
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
			f := newMockContestRepositoryFields(t, ctrl)
			tt.setup(f, tt.args)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.AddContestTeamMembers(context.Background(), tt.args.teamID, tt.args.members))
		})
	}
}

func TestContestRepository_EditContestTeamMembers(t *testing.T) {
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
		{
			name: "Success",
			args: args{
				teamID: random.UUID(),
				members: []uuid.UUID{
					random.UUID(),
					random.UUID(),
				}},
			setup: func(f mockContestRepositoryFields, args args) {
				memberToBeRemained := args.members[0]
				memberToBeAdded := args.members[1]
				memberToBeRemoved := random.UUID()
				rows := sqlmock.NewRows([]string{"team_id", "user_id"})
				rows.AddRow(args.teamID, memberToBeRemained)
				rows.AddRow(args.teamID, memberToBeRemoved)

				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(rows)
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `contest_team_user_belongings` (`team_id`,`user_id`,`created_at`,`updated_at`) VALUES (?,?,?,?)")).
					WithArgs(args.teamID, memberToBeAdded, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("DELETE FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ? AND `contest_team_user_belongings`.`user_id` IN (?)")).
					WithArgs(args.teamID, memberToBeRemoved).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "ContestTeamNotFound",
			args: args{
				teamID: random.UUID(),
				members: []uuid.UUID{
					random.UUID(),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnError(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_FindBelongings",
			args: args{
				teamID: random.UUID(),
				members: []uuid.UUID{
					random.UUID(),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_InsertBelongings",
			args: args{
				teamID: random.UUID(),
				members: []uuid.UUID{
					random.UUID(),
					random.UUID(),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				rows := sqlmock.NewRows([]string{"team_id", "user_id"})
				memberToBeRemained := args.members[0]
				memberToBeAdded := args.members[1]
				rows.AddRow(args.teamID, memberToBeRemained)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(rows)
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `contest_team_user_belongings` (`team_id`,`user_id`,`created_at`,`updated_at`) VALUES (?,?,?,?)")).
					WithArgs(args.teamID, memberToBeAdded, anyTime{}, anyTime{}).
					WillReturnError(errUnexpected)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_DeleteBelongings",
			args: args{
				teamID: random.UUID(),
				members: []uuid.UUID{
					random.UUID(),
					random.UUID(),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				rows := sqlmock.NewRows([]string{"team_id", "user_id"})
				memberToBeRemained := args.members[0]
				memberToBeAdded := args.members[1]
				memberToBeRemoved := random.UUID()
				rows.AddRow(args.teamID, memberToBeRemained)
				rows.AddRow(args.teamID, memberToBeRemoved)

				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(rows)
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `contest_team_user_belongings` (`team_id`,`user_id`,`created_at`,`updated_at`) VALUES (?,?,?,?)")).
					WithArgs(args.teamID, memberToBeAdded, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("DELETE FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ? AND `contest_team_user_belongings`.`user_id` IN (?)")).
					WithArgs(args.teamID, memberToBeRemoved).
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
			f := newMockContestRepositoryFields(t, ctrl)
			tt.setup(f, tt.args)
			repo := NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.EditContestTeamMembers(context.Background(), tt.args.teamID, tt.args.members))
		})
	}
}
