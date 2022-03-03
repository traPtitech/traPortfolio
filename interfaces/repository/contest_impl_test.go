package repository_test

import (
	"database/sql/driver"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/database/mock_database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external"
	impl "github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
	"gorm.io/gorm"
)

type mockContestRepositoryFields struct {
	h      database.SQLHandler
	portal external.PortalAPI
}

func newMockContestRepositoryFields(ctrl *gomock.Controller) mockContestRepositoryFields {
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
					Name:      random.AlphaNumeric(rand.Intn(30) + 1),
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
			ctrl := gomock.NewController(t)
			f := newMockContestRepositoryFields(ctrl)
			tt.setup(f, tt.want)
			repo := impl.NewContestRepository(f.h, f.portal)
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
					Name:      random.AlphaNumeric(rand.Intn(30) + 1),
					TimeStart: sampleTime,
					TimeEnd:   sampleTime,
				},
				Link:        random.RandURLString(),
				Description: random.AlphaNumeric(rand.Intn(30) + 1),
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
			ctrl := gomock.NewController(t)
			f := newMockContestRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.GetContest(tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestRepository_CreateContest(t *testing.T) {
	cname := random.AlphaNumeric(rand.Intn(30) + 1)       // Successで使用するContest.Name
	link := random.AlphaNumeric(rand.Intn(30) + 1)        // Successで使用するContestDetail.Link
	description := random.AlphaNumeric(rand.Intn(30) + 1) // Successで使用するContestDetail.Description

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
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Link:        optional.NewString(random.RandURLString(), true),
					Since:       sampleTime,
					Until:       optional.NewTime(sampleTime, true),
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
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `contests` (`id`,`name`,`description`,`link`,`since`,`until`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Name, args.args.Description, args.args.Link, args.args.Since, args.args.Until, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				h.Mock.ExpectCommit()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(anyUUID{}).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "name", "since", "until", "link", "description"}).
							// do not compare id
							AddRow(random.UUID(), want.Contest.Name, want.Contest.TimeStart, want.Contest.TimeEnd, want.Link, want.Description),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				args: &repository.CreateContestArgs{
					Name:        random.AlphaNumeric(rand.Intn(30) + 1),
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Link:        optional.NewString(random.RandURLString(), true),
					Since:       sampleTime,
					Until:       optional.NewTime(sampleTime, true),
				},
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestDetail) {
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
			ctrl := gomock.NewController(t)
			f := newMockContestRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.CreateContest(tt.args.args)
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
					Name:        optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Description: optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Link:        optional.NewString(random.RandURLString(), true),
					Since:       optional.NewTime(sampleTime, true),
					Until:       optional.NewTime(sampleTime, true),
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
					WithArgs(args.args.Description.String, args.args.Link.String, args.args.Name.String, args.args.Since.Time, args.args.Until.Time, anyTime{}, args.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until", "created_at", "updated_at"}).
							AddRow(args.id, args.args.Name.String, args.args.Description.String, args.args.Link.String, args.args.Since.Time, args.args.Until.Time, time.Time{}, time.Time{}),
					)
				h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				id: random.UUID(),
				args: &repository.UpdateContestArgs{
					Name:        optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Description: optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Link:        optional.NewString(random.RandURLString(), true),
					Since:       optional.NewTime(sampleTime, true),
					Until:       optional.NewTime(sampleTime, true),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(gorm.ErrRecordNotFound)
				h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError",
			args: args{
				id: random.UUID(),
				args: &repository.UpdateContestArgs{
					Name:        optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Description: optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Link:        optional.NewString(random.RandURLString(), true),
					Since:       optional.NewTime(sampleTime, true),
					Until:       optional.NewTime(sampleTime, true),
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
					WithArgs(args.args.Description.String, args.args.Link.String, args.args.Name.String, args.args.Since.Time, args.args.Until.Time, anyTime{}, args.id).
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
			ctrl := gomock.NewController(t)
			f := newMockContestRepositoryFields(ctrl)
			tt.setup(f, tt.args)
			repo := impl.NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.UpdateContest(tt.args.id, tt.args.args))
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
					WillReturnError(gorm.ErrRecordNotFound)
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
			ctrl := gomock.NewController(t)
			f := newMockContestRepositoryFields(ctrl)
			tt.setup(f, tt.args)
			repo := impl.NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.DeleteContest(tt.args.id))
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
					ID:        random.UUID(),
					ContestID: cid,
					Name:      random.AlphaNumeric(rand.Intn(30) + 1),
					Result:    random.AlphaNumeric(rand.Intn(30) + 1),
				},
			},
			setup: func(f mockContestRepositoryFields, args args, want []*domain.ContestTeam) {
				rows := sqlmock.NewRows([]string{"id", "contest_id", "name", "result"})
				for _, v := range want {
					rows.AddRow(v.ID, v.ContestID, v.Name, v.Result)
				}
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.contestID))
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`contest_id` = ?")).
					WithArgs(args.contestID).
					WillReturnRows(rows)
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
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnError(gorm.ErrRecordNotFound)
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
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ? ORDER BY `contests`.`id` LIMIT 1")).
					WithArgs(args.contestID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.contestID))
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE contest_id = ?")).
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
			f := newMockContestRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.GetContestTeams(tt.args.contestID)
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
					ID:        tid,
					ContestID: cid,
					Name:      random.AlphaNumeric(rand.Intn(30) + 1),
					Result:    random.AlphaNumeric(rand.Intn(30) + 1),
				},
				Link:        random.RandURLString(),
				Description: random.AlphaNumeric(rand.Intn(30) + 1),
				// Members
			},
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestTeamDetail) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? AND `contest_teams`.`contest_id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID, args.contestID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id", "name", "result", "link", "description"}).
							AddRow(want.ContestTeam.ID, want.ContestTeam.ContestID, want.ContestTeam.Name, want.ContestTeam.Result, want.Link, want.Description),
					)
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
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? AND `contest_teams`.`contest_id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID, args.contestID).
					WillReturnError(gorm.ErrRecordNotFound)
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
			f := newMockContestRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.GetContestTeam(tt.args.contestID, tt.args.teamID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestRepository_CreateContestTeam(t *testing.T) {
	cid := random.UUID() // Successで使うcontestID
	successArgs := repository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(rand.Intn(30) + 1),
		Result:      optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
		Link:        optional.NewString(random.RandURLString(), true),
		Description: random.AlphaNumeric(rand.Intn(30) + 1),
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
					// ID: Assertion時にgot.IDと合わせる
					ContestID: cid,
					Name:      successArgs.Name,
					Result:    successArgs.Result.String,
				},
				Link:        successArgs.Link.String,
				Description: successArgs.Description,
				Members:     nil,
			},
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestTeamDetail) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `contest_teams` (`id`,`contest_id`,`name`,`description`,`result`,`link`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.contestID, args._contestTeam.Name, args._contestTeam.Description, args._contestTeam.Result, args._contestTeam.Link, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				contestID: cid,
				_contestTeam: &repository.CreateContestTeamArgs{
					Name:        random.AlphaNumeric(rand.Intn(30) + 1),
					Result:      optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Link:        optional.NewString(random.RandURLString(), true),
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
				},
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want *domain.ContestTeamDetail) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `contest_teams` (`id`,`contest_id`,`name`,`description`,`result`,`link`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.contestID, args._contestTeam.Name, args._contestTeam.Description, args._contestTeam.Result, args._contestTeam.Link, anyTime{}, anyTime{}).
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
			ctrl := gomock.NewController(t)
			f := newMockContestRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewContestRepository(f.h, f.portal)
			// Assertion
			got, err := repo.CreateContestTeam(tt.args.contestID, tt.args._contestTeam)
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
					Name:        optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Description: optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Link:        optional.NewString(random.RandURLString(), true),
					Result:      optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				cid := random.UUID()
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id", "name", "description", "result", "link", "created_at", "updated_at"}).
							AddRow(args.teamID, cid, "", "", "", "", time.Time{}, time.Time{}),
					)
				h.Mock.
					ExpectExec(regexp.QuoteMeta("UPDATE `contest_teams` SET `description`=?,`link`=?,`name`=?,`result`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.Description.String, args.args.Link.String, args.args.Name.String, args.args.Result.String, anyTime{}, args.teamID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id", "name", "description", "result", "link", "created_at", "updated_at"}).
							AddRow(args.teamID, cid, args.args.Name, args.args.Description, args.args.Result, args.args.Link, time.Time{}, time.Time{}),
					)
				h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				teamID: random.UUID(),
				args: &repository.UpdateContestTeamArgs{
					Name:        optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Description: optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Link:        optional.NewString(random.RandURLString(), true),
					Result:      optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnError(gorm.ErrRecordNotFound)
				h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError",
			args: args{
				teamID: random.UUID(),
				args: &repository.UpdateContestTeamArgs{
					Name:        optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Description: optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Link:        optional.NewString(random.RandURLString(), true),
					Result:      optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
				},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id", "name", "description", "result", "link", "created_at", "updated_at"}).
							AddRow(args.teamID, random.UUID(), "", "", "", "", time.Time{}, time.Time{}),
					)
				h.Mock.
					ExpectExec(regexp.QuoteMeta("UPDATE `contest_teams` SET `description`=?,`link`=?,`name`=?,`result`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.Description.String, args.args.Link.String, args.args.Name.String, args.args.Result.String, anyTime{}, args.teamID).
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
			ctrl := gomock.NewController(t)
			f := newMockContestRepositoryFields(ctrl)
			tt.setup(f, tt.args)
			repo := impl.NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.UpdateContestTeam(tt.args.teamID, tt.args.args))
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
				{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(rand.Intn(30) + 1),
					RealName: random.AlphaNumeric(rand.Intn(30) + 1),
				},
			},
			setup: func(f mockContestRepositoryFields, args args, want []*domain.User) {
				u := want[0]
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"team_id", "user_id"}).
							AddRow(args.teamID, u.ID),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ?")).
					WithArgs(u.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name"}).
							AddRow(u.ID, u.Name),
					)
				p := f.portal.(*mock_external.MockPortalAPI)
				p.EXPECT().GetAll().Return(makePortalUsers(want), nil)
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
				{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(rand.Intn(30) + 1),
					RealName: random.AlphaNumeric(rand.Intn(30) + 1),
				},
				{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(rand.Intn(30) + 1),
					RealName: random.AlphaNumeric(rand.Intn(30) + 1),
				},
			},
			setup: func(f mockContestRepositoryFields, args args, want []*domain.User) {
				h := f.h.(*mock_database.MockSQLHandler)
				belongingRows := sqlmock.NewRows([]string{"team_id", "user_id"})
				for _, u := range want {
					belongingRows.AddRow(args.teamID, u.ID)
				}
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(belongingRows)
				userIDs := make([]driver.Value, len(want))
				userRows := sqlmock.NewRows([]string{"id", "name"})
				for i, u := range want {
					userIDs[i] = u.ID
					userRows.AddRow(u.ID, u.Name)
				}
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` IN (?,?)")).
					WithArgs(userIDs...).
					WillReturnRows(userRows)
				p := f.portal.(*mock_external.MockPortalAPI)
				p.EXPECT().GetAll().Return(makePortalUsers(want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			want: nil,
			setup: func(f mockContestRepositoryFields, args args, want []*domain.User) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
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
				u := &domain.User{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(rand.Intn(30) + 1),
					RealName: random.AlphaNumeric(rand.Intn(30) + 1),
				}
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"team_id", "user_id"}).
							AddRow(args.teamID, u.ID),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ?")).
					WithArgs(u.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name"}).
							AddRow(u.ID, u.Name),
					)
				p := f.portal.(*mock_external.MockPortalAPI)
				p.EXPECT().GetAll().Return(nil, errUnexpected)
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
			f := newMockContestRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewContestRepository(f.h, f.portal)
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
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(rows)
				h.Mock.ExpectBegin()
				for _, u := range newUsers {
					h.Mock.
						ExpectExec(regexp.QuoteMeta("INSERT INTO `contest_team_user_belongings` (`team_id`,`user_id`,`created_at`,`updated_at`) VALUES (?,?,?,?)")).
						WithArgs(args.teamID, u, anyTime{}, anyTime{}).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}
				h.Mock.ExpectCommit()
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
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnError(gorm.ErrRecordNotFound)
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
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
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
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "team_id", "user_id"}).
							AddRow(random.UUID(), args.teamID, random.UUID()),
					)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `contest_team_user_belongings` (`team_id`,`user_id`,`created_at`,`updated_at`) VALUES (?,?,?,?)")).
					WithArgs(args.teamID, args.members[0], anyTime{}, anyTime{}).
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
			ctrl := gomock.NewController(t)
			f := newMockContestRepositoryFields(ctrl)
			tt.setup(f, tt.args)
			repo := impl.NewContestRepository(f.h, f.portal)
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
		{
			name: "Success",
			args: args{
				teamID:  random.UUID(),
				members: []uuid.UUID{random.UUID()},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				rows := sqlmock.NewRows([]string{"team_id", "user_id"})
				for _, member := range args.members {
					rows.AddRow(args.teamID, member)
				}
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(rows)
				h.Mock.ExpectBegin()
				for _, v := range args.members {
					h.Mock.
						ExpectExec(regexp.QuoteMeta("DELETE FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ? AND `contest_team_user_belongings`.`user_id` = ?")).
						WithArgs(args.teamID, v).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}
				h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "ContestTeamNotFound",
			args: args{
				teamID:  random.UUID(),
				members: []uuid.UUID{random.UUID()},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnError(gorm.ErrRecordNotFound)
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
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_DeleteBelongings",
			args: args{
				teamID:  random.UUID(),
				members: []uuid.UUID{random.UUID()},
			},
			setup: func(f mockContestRepositoryFields, args args) {
				rows := sqlmock.NewRows([]string{"team_id", "user_id"})
				for _, member := range args.members {
					rows.AddRow(args.teamID, member)
				}
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ? ORDER BY `contest_teams`.`id` LIMIT 1")).
					WithArgs(args.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "contest_id"}).
							AddRow(args.teamID, random.UUID()),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ?")).
					WithArgs(args.teamID).
					WillReturnRows(rows)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectExec(regexp.QuoteMeta("DELETE FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`team_id` = ? AND `contest_team_user_belongings`.`user_id` = ?")).
					WithArgs(args.teamID, args.members[0]).
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
			ctrl := gomock.NewController(t)
			f := newMockContestRepositoryFields(ctrl)
			tt.setup(f, tt.args)
			repo := impl.NewContestRepository(f.h, f.portal)
			// Assertion
			tt.assertion(t, repo.DeleteContestTeamMembers(tt.args.teamID, tt.args.members))
		})
	}
}
