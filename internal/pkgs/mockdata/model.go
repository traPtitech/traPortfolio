package mockdata

import (
	"time"

	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

var (
	MockUsers                     = CloneMockUsers()
	MockAccounts                  = CloneMockAccounts()
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
			State:       domain.TraqStateActive,
		},
		{
			ID:          UserID2(),
			Description: "I am user2",
			Check:       true,
			Name:        userName2,
			State:       domain.TraqStateDeactivated,
		},
		{
			ID:          UserID3(),
			Description: "I am lolico",
			Check:       false,
			Name:        userName3,
			State:       domain.TraqStateActive,
		},
	}
}

func CloneMockAccounts() []model.Account {
	return []model.Account{
		{
			ID:     AccountID1(),
			Type:   2,
			Name:   "sample_account_display_name",
			URL:    "https://twitter.com/sample_account",
			UserID: UserID1(),
			Check:  true,
		},
	}
}

func CloneMockContests() []model.Contest {
	return []model.Contest{
		{
			ID:          ContestID1(),
			Name:        "sample_contest_name",
			Description: "sample_contest_description",
			Since:       time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			Until:       time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          ContestID2(),
			Name:        "sample_contest_name2",
			Description: "sample_contest_description2",
			Since:       time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			Until:       time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          ContestID3(),
			Name:        "sample_contest_name3",
			Description: "sample_contest_description3",
			Since:       time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			Until:       time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}
}

func CloneMockContestLinks() []model.ContestLink {
	return []model.ContestLink{
		{
			ID:    ContestID1(),
			Order: 0,
			Link:  "https://sample.contests1.com",
		},
		{
			ID:    ContestID1(),
			Order: 1,
			Link:  "https://twitter.com/contest",
		},
		{
			ID:    ContestID2(),
			Order: 0,
			Link:  "https://sample.contests2.com",
		},
		{
			ID:    ContestID3(),
			Order: 0,
			Link:  "https://sample.contests3.com",
		},
	}
}

func CloneMockContestTeams() []model.ContestTeam {
	return []model.ContestTeam{
		{
			ID:          ContestTeamID1(),
			ContestID:   ContestID1(),
			Name:        "sample_contest_team_name",
			Description: "sample_contest_team_description",
			Result:      "sample_contest_team_result",
		},
		{
			ID:          ContestTeamID2(),
			ContestID:   ContestID1(),
			Name:        "sample_contest_team_name2",
			Description: "sample_contest_team_description2",
			Result:      "sample_contest_team_result2",
		},
		{
			ID:          ContestTeamID3(),
			ContestID:   ContestID1(),
			Name:        "sample_contest_team_name3",
			Description: "sample_contest_team_description3",
			Result:      "sample_contest_team_result3",
		},
	}
}

func CloneMockContestTeamUserBelongings() []model.ContestTeamUserBelonging {
	return []model.ContestTeamUserBelonging{
		{
			TeamID: ContestTeamID1(),
			UserID: UserID1(),
		},
	}
}

func CloneMockContestTeamLinks() []model.ContestTeamLink {
	return []model.ContestTeamLink{
		{
			ID:    ContestTeamID1(),
			Order: 0,
			Link:  "https://sample.contest_teams1.com",
		},
		{
			ID:    ContestTeamID1(),
			Order: 1,
			Link:  "https://twitter.com/contest_team1",
		},
		{
			ID:    ContestTeamID2(),
			Order: 0,
			Link:  "https://sample.contest_teams2.com",
		},
		{
			ID:    ContestTeamID3(),
			Order: 0,
			Link:  "https://sample.contest_teams3.com",
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
		{
			ID:    KnoqEventID3(),
			Level: domain.EventLevelAnonymous,
		},
	}
}

func CloneMockGroups() []model.Group {
	return []model.Group{
		{
			GroupID:     GroupID1(),
			Name:        "sample_group_name",
			Description: "sample_group_description",
		},
	}
}

func CloneMockGroupUserBelongings() []model.GroupUserBelonging {
	return []model.GroupUserBelonging{
		{
			UserID:        UserID1(),
			GroupID:       GroupID1(),
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
			GroupID: GroupID1(),
		},
	}
}

func CloneMockGroupLinks() []model.GroupLink {
	return []model.GroupLink{
		{
			ID:    GroupID1(),
			Order: 0,
			Link:  "https://sample.group1.com",
		},
		{
			ID:    GroupID1(),
			Order: 1,
			Link:  "https://twitter.com/group1",
		},
	}
}

func CloneMockProjects() []*model.Project {
	return []*model.Project{
		{
			ID:            ProjectID1(),
			Name:          "sample_project_name1",
			Description:   "sample_project_description1",
			SinceYear:     2021,
			SinceSemester: 0,
			UntilYear:     2021,
			UntilSemester: 1,
		},
		{
			ID:            ProjectID2(),
			Name:          "sample_project_name2",
			Description:   "sample_project_description2",
			SinceYear:     2022,
			SinceSemester: 0,
			UntilYear:     2022,
			UntilSemester: 1,
		},
		{
			ID:            ProjectID3(),
			Name:          "sample_project_name3",
			Description:   "sample_project_description3",
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
			ProjectID:     ProjectID1(),
			UserID:        UserID1(),
			SinceYear:     2021,
			SinceSemester: 0,
			UntilYear:     2021,
			UntilSemester: 1,
		},
		{
			ProjectID:     ProjectID1(),
			UserID:        UserID2(),
			SinceYear:     2022,
			SinceSemester: 0,
			UntilYear:     2022,
			UntilSemester: 1,
		},
		{
			ProjectID:     ProjectID2(),
			UserID:        UserID2(),
			SinceYear:     2021,
			SinceSemester: 0,
			UntilYear:     2022,
			UntilSemester: 1,
		},
	}
}

func CloneMockProjectLinks() []model.ProjectLink {
	return []model.ProjectLink{
		{
			ID:    ProjectID1(),
			Order: 0,
			Link:  "https://sample.project1.com",
		},
		{
			ID:    ProjectID1(),
			Order: 1,
			Link:  "https://twitter.com/project1",
		},
		{
			ID:    ProjectID2(),
			Order: 0,
			Link:  "https://sample.project2.com",
		},
		{
			ID:    ProjectID3(),
			Order: 0,
			Link:  "https://sample.project3.com",
		},
	}
}

func InsertSampleDataToDB(h *gorm.DB) error {
	mockUsers := CloneMockUsers()
	if err := h.Create(&mockUsers).Error; err != nil {
		return err
	}

	mockAccounts := CloneMockAccounts()
	if err := h.Create(&mockAccounts).Error; err != nil {
		return err
	}

	mockContests := CloneMockContests()
	if err := h.Create(&mockContests).Error; err != nil {
		return err
	}

	mockContestLinks := CloneMockContestLinks()
	if err := h.Create(&mockContestLinks).Error; err != nil {
		return err
	}

	mockContestTeams := CloneMockContestTeams()
	if err := h.Create(&mockContestTeams).Error; err != nil {
		return err
	}

	mockContestTeamUserBelongings := CloneMockContestTeamUserBelongings()
	if err := h.Create(&mockContestTeamUserBelongings).Error; err != nil {
		return err
	}

	mockContestTeamLinks := CloneMockContestTeamLinks()
	if err := h.Create(&mockContestTeamLinks).Error; err != nil {
		return err
	}

	mockEventLevelRelations := CloneMockEventLevelRelations()
	if err := h.Create(&mockEventLevelRelations).Error; err != nil {
		return err
	}

	mockGroups := CloneMockGroups()
	if err := h.Create(&mockGroups).Error; err != nil {
		return err
	}

	mockGroupLinks := CloneMockGroupLinks()
	if err := h.Create(mockGroupLinks).Error; err != nil {
		return err
	}

	mockGroupUserBelongings := CloneMockGroupUserBelongings()
	if err := h.Create(&mockGroupUserBelongings).Error; err != nil {
		return err
	}

	mockProjects := CloneMockProjects()
	if err := h.Create(&mockProjects).Error; err != nil {
		return err
	}

	mockProjectLinks := CloneMockProjectLinks()
	if err := h.Create(&mockProjectLinks).Error; err != nil {
		return err
	}

	mockGroupUserAdmins := CloneMockGroupUserAdmins()
	if err := h.Create(&mockGroupUserAdmins).Error; err != nil {
		return err
	}

	mockProjectMembers := CloneMockProjectMembers()
	return h.Create(&mockProjectMembers).Error
}
