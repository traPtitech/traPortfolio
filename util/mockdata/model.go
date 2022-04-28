package mockdata

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

var (
	MockUsers      = CloneMockUsers()
	CloneMockUsers = func() []*model.User {
		return []*model.User{
			{
				ID:          uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111"),
				Description: "I am user1",
				Check:       true,
				Name:        "user1",
			},
			{
				ID:          uuid.FromStringOrNil("22222222-2222-2222-2222-222222222222"),
				Description: "I am user2",
				Check:       true,
				Name:        "user2",
			},
			{
				ID:          uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333"),
				Description: "I am lolico",
				Check:       false,
				Name:        "lolico",
			},
		}
	}

	MockAccount      = CloneMockAccount()
	CloneMockAccount = func() model.Account {
		return model.Account{
			ID:     uuid.FromStringOrNil("d834e180-2af9-4cfe-838a-8a3930666490"),
			Type:   0,
			Name:   "sample_account_display_name",
			URL:    "https://sample.accounts.com",
			UserID: MockUsers[0].ID,
			Check:  true,
		}
	}

	MockContest      = CloneMockContest()
	CloneMockContest = func() model.Contest {
		return model.Contest{
			ID:          uuid.FromStringOrNil("08eec963-0f29-48d1-929f-004cb67d8ce6"),
			Name:        "sample_contest_name",
			Description: "sample_contest_description",
			Link:        "https://sample.contests.com",
			Since:       time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			Until:       time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		}
	}

	MockContestTeam      = CloneMockContestTeam()
	CloneMockContestTeam = func() model.ContestTeam {
		return model.ContestTeam{
			ID:          uuid.FromStringOrNil("a9d07124-ffee-412f-adfc-02d3db0b750d"),
			ContestID:   MockContest.ID,
			Name:        "sample_contest_team_name",
			Description: "sample_contest_team_description",
			Result:      "sample_contest_team_result",
			Link:        "https://sample.contest_teams.com",
		}
	}

	MockContestTeamUserBelonging      = CloneMockContestTeam()
	CloneMockContestTeamUserBelonging = func() model.ContestTeamUserBelonging {
		return model.ContestTeamUserBelonging{
			TeamID: MockContestTeam.ID,
			UserID: MockUsers[0].ID,
		}
	}

	MockEventLevelRelation      = CloneMockEventLevelRelation()
	CloneMockEventLevelRelation = func() model.EventLevelRelation {
		return model.EventLevelRelation{
			ID:    uuid.FromStringOrNil("e32a0431-aa0e-4825-98e6-479912275bbd"),
			Level: domain.EventLevelPublic,
		}
	}

	MockGroup      = CloneMockGroup()
	CloneMockGroup = func() model.Group {
		return model.Group{
			GroupID:     uuid.FromStringOrNil("455938b1-635f-4b43-ae74-66550b04c5d4"),
			Name:        "sample_group_name",
			Link:        "https://sample.groups.com",
			Description: "sample_group_description",
		}
	}

	MockGroupUserBelonging      = CloneMockGroupUserBelonging()
	CloneMockGroupUserBelonging = func() model.GroupUserBelonging {
		return model.GroupUserBelonging{
			UserID:        MockUsers[0].ID,
			GroupID:       MockGroup.GroupID,
			SinceYear:     2022,
			SinceSemester: 1,
			UntilYear:     2022,
			UntilSemester: 2,
		}
	}

	MockGroupUserAdmin      = CloneMockGroupUserAdmin()
	CloneMockGroupUserAdmin = func() model.GroupUserAdmin {
		return model.GroupUserAdmin{
			UserID:  MockUsers[0].ID,
			GroupID: MockGroup.GroupID,
		}
	}

	MockProject      = CloneMockProject()
	CloneMockProject = func() model.Project {
		return model.Project{
			ID:            uuid.FromStringOrNil("bf9c1aec-7e3a-4587-8adc-651895aa6ec0"),
			Name:          "sample_project_name",
			Description:   "sample_project_description",
			Link:          "https://sample.projects.com",
			SinceYear:     2022,
			SinceSemester: 1,
			UntilYear:     2022,
			UntilSemester: 2,
		}
	}

	MockProjectMember      = CloneMockProjectMember()
	CloneMockProjectMember = func() model.ProjectMember {
		return model.ProjectMember{
			ID:            uuid.FromStringOrNil("a211a49c-9b30-48b9-8dbb-c449c99f12c7"),
			ProjectID:     MockProject.ID,
			UserID:        MockUsers[0].ID,
			SinceYear:     2022,
			SinceSemester: 1,
			UntilYear:     2022,
			UntilSemester: 2,
		}
	}
)

func InsertSampleDataToDB(h database.SQLHandler) error {
	mockUsers := CloneMockUsers()
	if err := h.Create(&mockUsers).Error(); err != nil {
		return err
	}

	mockAccount := CloneMockAccount()
	if err := h.Create(&mockAccount).Error(); err != nil {
		return err
	}

	mockContest := CloneMockContest()
	if err := h.Create(&mockContest).Error(); err != nil {
		return err
	}

	mockContestTeam := CloneMockContestTeam()
	if err := h.Create(&mockContestTeam).Error(); err != nil {
		return err
	}

	mockContestTeamUserBelonging := CloneMockContestTeamUserBelonging()
	if err := h.Create(&mockContestTeamUserBelonging).Error(); err != nil {
		return err
	}

	mockEventLevelRelation := CloneMockEventLevelRelation()
	if err := h.Create(&mockEventLevelRelation).Error(); err != nil {
		return err
	}

	mockGroup := CloneMockGroup()
	if err := h.Create(&mockGroup).Error(); err != nil {
		return err
	}

	mockGroupUserBelonging := CloneMockGroupUserBelonging()
	if err := h.Create(&mockGroupUserBelonging).Error(); err != nil {
		return err
	}

	mockProject := CloneMockProject()
	if err := h.Create(&mockProject).Error(); err != nil {
		return err
	}

	mockProjectMember := CloneMockProjectMember()
	if err := h.Create(&mockProjectMember).Error(); err != nil {
		return err
	}

	return nil
}
