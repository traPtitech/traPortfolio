package mockdata

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

var (
	MockUsers = []*model.User{
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

	MockAccount = model.Account{
		ID:     uuid.FromStringOrNil("d834e180-2af9-4cfe-838a-8a3930666490"),
		Type:   0,
		Name:   "sample_account_display_name",
		URL:    "https://sample.accounts.com",
		UserID: MockUsers[0].ID,
		Check:  true,
	}

	MockContest = model.Contest{
		ID:          uuid.FromStringOrNil("08eec963-0f29-48d1-929f-004cb67d8ce6"),
		Name:        "sample_contest_name",
		Description: "sample_contest_description",
		Link:        "https://sample.contests.com",
		Since:       time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		Until:       time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
	}

	MockContestTeam = model.ContestTeam{
		ID:          uuid.FromStringOrNil("a9d07124-ffee-412f-adfc-02d3db0b750d"),
		ContestID:   MockContest.ID,
		Name:        "sample_contest_team_name",
		Description: "sample_contest_team_description",
		Result:      "sample_contest_team_result",
		Link:        "https://sample.contest_teams.com",
	}

	MockContestTeamUserBelonging = model.ContestTeamUserBelonging{
		TeamID: MockContestTeam.ID,
		UserID: MockUsers[0].ID,
	}

	MockEventLevelRelation = model.EventLevelRelation{
		ID:    uuid.FromStringOrNil("e32a0431-aa0e-4825-98e6-479912275bbd"),
		Level: domain.EventLevelPublic,
	}

	MockGroup = model.Group{
		GroupID:     uuid.FromStringOrNil("455938b1-635f-4b43-ae74-66550b04c5d4"),
		Name:        "sample_group_name",
		Link:        "https://sample.groups.com",
		Description: "sample_group_description",
	}

	MockGroupUserBelonging = model.GroupUserBelonging{
		UserID:        MockUsers[0].ID,
		GroupID:       MockGroup.GroupID,
		SinceYear:     2022,
		SinceSemester: 1,
		UntilYear:     2022,
		UntilSemester: 2,
	}

	MockProject = model.Project{
		ID:            uuid.FromStringOrNil("bf9c1aec-7e3a-4587-8adc-651895aa6ec0"),
		Name:          "sample_project_name",
		Description:   "sample_project_description",
		Link:          "https://sample.projects.com",
		SinceYear:     2022,
		SinceSemester: 1,
		UntilYear:     2022,
		UntilSemester: 2,
	}

	MockProjectMember = model.ProjectMember{
		ID:            uuid.FromStringOrNil("a211a49c-9b30-48b9-8dbb-c449c99f12c7"),
		ProjectID:     MockProject.ID,
		UserID:        MockUsers[0].ID,
		SinceYear:     2022,
		SinceSemester: 1,
		UntilYear:     2022,
		UntilSemester: 2,
	}
)

func InsertSampleDataToDB(h database.SQLHandler) error {
	if err := h.Create(&MockUsers).Error(); err != nil {
		return err
	}

	if err := h.Create(&MockAccount).Error(); err != nil {
		return err
	}

	if err := h.Create(&MockContest).Error(); err != nil {
		return err
	}

	if err := h.Create(&MockContestTeam).Error(); err != nil {
		return err
	}

	if err := h.Create(&MockContestTeamUserBelonging).Error(); err != nil {
		return err
	}

	if err := h.Create(&MockEventLevelRelation).Error(); err != nil {
		return err
	}

	if err := h.Create(&MockGroup).Error(); err != nil {
		return err
	}

	if err := h.Create(&MockGroupUserBelonging).Error(); err != nil {
		return err
	}

	if err := h.Create(&MockProject).Error(); err != nil {
		return err
	}

	if err := h.Create(&MockProjectMember).Error(); err != nil {
		return err
	}

	return nil
}
