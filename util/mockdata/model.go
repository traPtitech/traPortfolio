package mockdata

import (
	"time"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

var (
	MockUsers      = CloneMockUsers()
	CloneMockUsers = func() []*model.User {
		return []*model.User{
			{
				ID:          userID1,
				Description: "I am user1",
				Check:       true,
				Name:        userName1,
			},
			{
				ID:          userID2,
				Description: "I am user2",
				Check:       true,
				Name:        userName2,
			},
			{
				ID:          userID3,
				Description: "I am lolico",
				Check:       false,
				Name:        userName3,
			},
		}
	}

	MockAccount      = CloneMockAccount()
	CloneMockAccount = func() model.Account {
		return model.Account{
			ID:     accountID,
			Type:   0,
			Name:   "sample_account_display_name",
			URL:    "https://sample.accounts.com",
			UserID: userID1,
			Check:  true,
		}
	}

	MockContest      = CloneMockContest()
	CloneMockContest = func() model.Contest {
		return model.Contest{
			ID:          contestID,
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
			ID:          contestTeamID,
			ContestID:   contestID,
			Name:        "sample_contest_team_name",
			Description: "sample_contest_team_description",
			Result:      "sample_contest_team_result",
			Link:        "https://sample.contest_teams.com",
		}
	}

	MockContestTeamUserBelonging      = CloneMockContestTeam()
	CloneMockContestTeamUserBelonging = func() model.ContestTeamUserBelonging {
		return model.ContestTeamUserBelonging{
			TeamID: contestTeamID,
			UserID: userID1,
		}
	}

	MockEventLevelRelation      = CloneMockEventLevelRelation()
	CloneMockEventLevelRelation = func() model.EventLevelRelation {
		return model.EventLevelRelation{
			ID:    eventID,
			Level: domain.EventLevelPublic,
		}
	}

	MockGroup      = CloneMockGroup()
	CloneMockGroup = func() model.Group {
		return model.Group{
			GroupID:     groupID,
			Name:        "sample_group_name",
			Link:        "https://sample.groups.com",
			Leader:      userID1,
			Description: "sample_group_description",
		}
	}

	MockGroupUserBelonging      = CloneMockGroupUserBelonging()
	CloneMockGroupUserBelonging = func() model.GroupUserBelonging {
		return model.GroupUserBelonging{
			UserID:        userID1,
			GroupID:       MockGroup.GroupID,
			SinceYear:     2022,
			SinceSemester: 1,
			UntilYear:     2022,
			UntilSemester: 2,
		}
	}

	MockProject      = CloneMockProject()
	CloneMockProject = func() model.Project {
		return model.Project{
			ID:            projectID,
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
			ID:            projectMemberID,
			ProjectID:     projectID,
			UserID:        userID1,
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
