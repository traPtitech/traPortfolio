package mockdata

import (
	"time"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

var (
	MockUsers                    = CloneMockUsers()
	MockAccount                  = CloneMockAccount()
	MockContest                  = CloneMockContest()
	MockContestTeam              = CloneMockContestTeam()
	MockContestTeamUserBelonging = CloneMockContestTeam()
	MockEventLevelRelations      = CloneMockEventLevelRelations()
	MockGroup                    = CloneMockGroup()
	MockGroupUserBelonging       = CloneMockGroupUserBelonging()
	MockProject                  = CloneMockProject()
	MockProjectMember            = CloneMockProjectMember()
)

func CloneMockUsers() []*model.User {
	return []*model.User{
		{
			ID:          userID1.uuid(),
			Description: "I am user1",
			Check:       true,
			Name:        userName1,
		},
		{
			ID:          userID2.uuid(),
			Description: "I am user2",
			Check:       true,
			Name:        userName2,
		},
		{
			ID:          userID3.uuid(),
			Description: "I am lolico",
			Check:       false,
			Name:        userName3,
		},
	}
}

func CloneMockAccount() model.Account {
	return model.Account{
		ID:     accountID.uuid(),
		Type:   0,
		Name:   "sample_account_display_name",
		URL:    "https://sample.accounts.com",
		UserID: userID1.uuid(),
		Check:  true,
	}
}

func CloneMockContest() model.Contest {
	return model.Contest{
		ID:          contestID.uuid(),
		Name:        "sample_contest_name",
		Description: "sample_contest_description",
		Link:        "https://sample.contests.com",
		Since:       time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		Until:       time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
	}
}

func CloneMockContestTeam() model.ContestTeam {
	return model.ContestTeam{
		ID:          contestTeamID.uuid(),
		ContestID:   contestID.uuid(),
		Name:        "sample_contest_team_name",
		Description: "sample_contest_team_description",
		Result:      "sample_contest_team_result",
		Link:        "https://sample.contest_teams.com",
	}
}

func CloneMockContestTeamUserBelonging() model.ContestTeamUserBelonging {
	return model.ContestTeamUserBelonging{
		TeamID: contestTeamID.uuid(),
		UserID: userID1.uuid(),
	}
}

func CloneMockEventLevelRelations() []model.EventLevelRelation {
	return []model.EventLevelRelation{
		{
			ID:    knoqEventID1.uuid(),
			Level: domain.EventLevelPublic,
		},
		{
			ID:    knoqEventID2.uuid(),
			Level: domain.EventLevelPrivate,
		},
	}
}

func CloneMockGroupUserAdmin() []model.GroupUserAdmin {
	return []model.GroupUserAdmin{
		{
			UserID:  userID1.uuid(),
			GroupID: groupID.uuid(),
		},
	}
}

func CloneMockGroup() model.Group {
	return model.Group{
		GroupID:     groupID.uuid(),
		Name:        "sample_group_name",
		Link:        "https://sample.groups.com",
		Description: "sample_group_description",
	}
}

func CloneMockGroupUserBelonging() model.GroupUserBelonging {
	return model.GroupUserBelonging{
		UserID:        userID1.uuid(),
		GroupID:       MockGroup.GroupID,
		SinceYear:     2022,
		SinceSemester: 1,
		UntilYear:     2022,
		UntilSemester: 2,
	}
}

func CloneMockProject() model.Project {
	return model.Project{
		ID:            projectID.uuid(),
		Name:          "sample_project_name",
		Description:   "sample_project_description",
		Link:          "https://sample.projects.com",
		SinceYear:     2022,
		SinceSemester: 1,
		UntilYear:     2022,
		UntilSemester: 2,
	}
}

func CloneMockProjectMember() model.ProjectMember {
	return model.ProjectMember{
		ID:            projectMemberID.uuid(),
		ProjectID:     projectID.uuid(),
		UserID:        userID1.uuid(),
		SinceYear:     2022,
		SinceSemester: 1,
		UntilYear:     2022,
		UntilSemester: 2,
	}
}

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

	mockEventLevelRelations := CloneMockEventLevelRelations()
	if err := h.Create(&mockEventLevelRelations).Error(); err != nil {
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

	mockGroupUserAdmin := CloneMockGroupUserAdmin()
	if err := h.Create(&mockGroupUserAdmin).Error(); err != nil {
		return err
	}

	mockProjectMember := CloneMockProjectMember()
	if err := h.Create(&mockProjectMember).Error(); err != nil {
		return err
	}

	return nil
}
