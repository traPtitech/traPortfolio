package mockdata

import (
	"time"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

var (
	MockUsers                     = CloneMockUsers()
	MockAccount                   = CloneMockAccount()
	MockContests                  = CloneMockContests()
	MockContestTeams              = CloneMockContestTeams()
	MockContestTeamUserBelongings = CloneMockContestTeamUserBelongings()
	MockEventLevelRelations       = CloneMockEventLevelRelations()
	MockGroups                    = CloneMockGroups()
	MockGroupUserBelongings       = CloneMockGroupUserBelongings()
	MockGroupUserAdmins           = CloneMockGroupUserAdmins()
	MockProjects                  = CloneMockProjects()
	MockProjectMembers            = CloneMockProjectMembers()
)

func CloneMockUsers() []*model.User {
	return []*model.User{
		{
			ID:          UserID1(),
			Description: "I am user1",
			Check:       true,
			Name:        userName1,
		},
		{
			ID:          UserID2(),
			Description: "I am user2",
			Check:       true,
			Name:        userName2,
		},
		{
			ID:          UserID3(),
			Description: "I am lolico",
			Check:       false,
			Name:        userName3,
		},
	}
}

func CloneMockAccount() model.Account {
	return model.Account{
		ID:     AccountID(),
		Type:   0,
		Name:   "sample_account_display_name",
		URL:    "https://sample.accounts.com",
		UserID: UserID1(),
		Check:  true,
	}
}

func CloneMockContests() []model.Contest {
	return []model.Contest{
		{
			ID:          ContestID(),
			Name:        "sample_contest_name",
			Description: "sample_contest_description",
			Link:        "https://sample.contests.com",
			Since:       time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			Until:       time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}
}

func CloneMockContestTeams() []model.ContestTeam {
	return []model.ContestTeam{
		{
			ID:          ContestTeamID(),
			ContestID:   ContestID(),
			Name:        "sample_contest_team_name",
			Description: "sample_contest_team_description",
			Result:      "sample_contest_team_result",
			Link:        "https://sample.contest_teams.com",
		},
	}
}

func CloneMockContestTeamUserBelongings() []model.ContestTeamUserBelonging {
	return []model.ContestTeamUserBelonging{
		{
			TeamID: ContestTeamID(),
			UserID: UserID1(),
		},
	}
}

func CloneMockEventLevelRelations() []model.EventLevelRelation {
	return []model.EventLevelRelation{
		{
			ID:    KnoqEventID1(),
			Level: domain.EventLevelPublic,
		},
		{
			ID:    KnoqEventID2(),
			Level: domain.EventLevelPrivate,
		},
	}
}

func CloneMockGroups() []model.Group {
	return []model.Group{
		{
			GroupID:     GroupID(),
			Name:        "sample_group_name",
			Link:        "https://sample.groups.com",
			Description: "sample_group_description",
		},
	}
}

func CloneMockGroupUserBelongings() []model.GroupUserBelonging {
	return []model.GroupUserBelonging{
		{
			UserID:        UserID1(),
			GroupID:       GroupID(),
			SinceYear:     2022,
			SinceSemester: 0,
			UntilYear:     2022,
			UntilSemester: 1,
		},
	}
}

func CloneMockGroupUserAdmins() []model.GroupUserAdmin {
	return []model.GroupUserAdmin{
		{
			UserID:  UserID1(),
			GroupID: GroupID(),
		},
	}
}

func CloneMockProjects() []*model.Project {
	return []*model.Project{
		{
			ID:            ProjectID1(),
			Name:          "sample_project_name1",
			Description:   "sample_project_description1",
			Link:          "https://sample.project1.com",
			SinceYear:     2021,
			SinceSemester: 0,
			UntilYear:     2021,
			UntilSemester: 1,
		},
		{
			ID:            ProjectID2(),
			Name:          "sample_project_name2",
			Description:   "sample_project_description2",
			Link:          "https://sample.project2.com",
			SinceYear:     2022,
			SinceSemester: 0,
			UntilYear:     2022,
			UntilSemester: 1,
		},
		{
			ID:            ProjectID3(),
			Name:          "sample_project_name3",
			Description:   "sample_project_description3",
			Link:          "https://sample.project3.com",
			SinceYear:     2021,
			SinceSemester: 0,
			UntilYear:     2022,
			UntilSemester: 1,
		},
	}
}

func CloneMockProjectMembers() []*model.ProjectMember {
	return []*model.ProjectMember{
		{
			ID:            ProjectMemberID1(),
			ProjectID:     ProjectID1(),
			UserID:        UserID1(),
			SinceYear:     2021,
			SinceSemester: 0,
			UntilYear:     2021,
			UntilSemester: 1,
		},
		{
			ID:            ProjectMemberID2(),
			ProjectID:     ProjectID1(),
			UserID:        UserID2(),
			SinceYear:     2022,
			SinceSemester: 0,
			UntilYear:     2022,
			UntilSemester: 1,
		},
		{
			ID:            ProjectMemberID3(),
			ProjectID:     ProjectID2(),
			UserID:        UserID2(),
			SinceYear:     2021,
			SinceSemester: 0,
			UntilYear:     2022,
			UntilSemester: 1,
		},
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

	mockContests := CloneMockContests()
	if err := h.Create(&mockContests).Error(); err != nil {
		return err
	}

	mockContestTeams := CloneMockContestTeams()
	if err := h.Create(&mockContestTeams).Error(); err != nil {
		return err
	}

	mockContestTeamUserBelongings := CloneMockContestTeamUserBelongings()
	if err := h.Create(&mockContestTeamUserBelongings).Error(); err != nil {
		return err
	}

	mockEventLevelRelations := CloneMockEventLevelRelations()
	if err := h.Create(&mockEventLevelRelations).Error(); err != nil {
		return err
	}

	mockGroups := CloneMockGroups()
	if err := h.Create(&mockGroups).Error(); err != nil {
		return err
	}

	mockGroupUserBelongings := CloneMockGroupUserBelongings()
	if err := h.Create(&mockGroupUserBelongings).Error(); err != nil {
		return err
	}

	mockProjects := CloneMockProjects()
	if err := h.Create(&mockProjects).Error(); err != nil {
		return err
	}

	mockGroupUserAdmins := CloneMockGroupUserAdmins()
	if err := h.Create(&mockGroupUserAdmins).Error(); err != nil {
		return err
	}

	mockProjectMembers := CloneMockProjectMembers()
	if err := h.Create(&mockProjectMembers).Error(); err != nil {
		return err
	}

	return nil
}
