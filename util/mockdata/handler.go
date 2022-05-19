package mockdata

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
)

var (
	HMockContest        = CloneHandlerMockContest()
	HMockContestTeam    = CloneHandlerMockContestTeam()
	HMockEvents         = CloneHandlerMockEvents()
	HMockGroup          = CloneHandlerMockGroup()
	HMockGroupMembers   = CloneHandlerMockGroupMembers()
	HMockProject        = CloneHandlerMockProject()
	HMockProjectMembers = CloneHandlerMockProjectMembers()
	HMockUsers          = CloneHandlerMockUsers()
	HMockUserAccount    = CloneHandlerMockUserAccount()
	HMockUserContest    = CloneHandlerMockUserContest()
	HMockUserGroup      = CloneHandlerMockUserGroup()
	HMockUserProject    = CloneHandlerMockUserProject()
)

func CloneHandlerMockContest() handler.ContestDetail {
	var (
		mContest     = CloneMockContest()
		hContestTeam = CloneHandlerMockContestTeam()
	)

	return handler.ContestDetail{
		Contest: handler.Contest{
			Duration: handler.Duration{
				Since: mContest.Since,
				Until: &mContest.Until,
			},
			Id:   mContest.ID,
			Name: mContest.Name,
		},
		Description: mContest.Description,
		Link:        mContest.Link,
		Teams: []handler.ContestTeam{
			hContestTeam.ContestTeam,
		},
	}
}

func CloneHandlerMockContestTeam() handler.ContestTeamDetail {
	var (
		mContestTeam              = CloneMockContestTeam()
		mContestTeamUserBelonging = CloneMockContestTeamUserBelonging()
	)

	return handler.ContestTeamDetail{
		ContestTeam: handler.ContestTeam{
			Id:     mContestTeam.ContestID,
			Name:   mContestTeam.Name,
			Result: mContestTeam.Result,
		},
		Description: mContestTeam.Description,
		Link:        mContestTeam.Link,
		Members: []handler.User{
			getUser(mContestTeamUserBelonging.UserID).User,
		},
	}
}

func CloneHandlerMockEvents() []handler.EventDetail {
	var (
		mEventLevel = CloneMockEventLevelRelation()
		knoqEvents  = CloneMockKnoqEvents()
		hEvents     = make([]handler.EventDetail, len(knoqEvents))
	)

	for i, e := range knoqEvents {
		hEvents[i] = handler.EventDetail{
			Event: handler.Event{
				Duration: handler.Duration{
					Since: e.TimeStart,
					Until: &e.TimeEnd,
				},
				Id:   e.ID,
				Name: e.Name,
			},
			Description: e.Description,
			Hostname:    make([]handler.User, len(e.Admins)),
			Place:       e.Place,
		}

		for j, uid := range e.Admins {
			hEvents[i].Hostname[j] = getUser(uid).User
		}

		// TODO: 綺麗にする
		if i == 0 {
			hEvents[i].EventLevel = handler.EventLevel(mEventLevel.Level)
		} else if i == 1 {
			hEvents[i].EventLevel = handler.EventLevel(0)
		}
	}

	return hEvents
}

func CloneHandlerMockGroup() handler.GroupDetail {
	var (
		mGroup        = CloneMockGroup()
		hGroupMembers = CloneHandlerMockGroupMembers()
	)

	return handler.GroupDetail{
		Group: handler.Group{
			Id:   mGroup.GroupID,
			Name: mGroup.Name,
		},
		Description: mGroup.Description,
		Leader:      getUser(mGroup.Leader).User,
		Link:        mGroup.Link,
		Members:     hGroupMembers,
	}
}

func CloneHandlerMockGroupMembers() []handler.GroupMember {
	var (
		mGroupUserbelonging = CloneMockGroupUserBelonging()
	)

	return []handler.GroupMember{
		{
			User: getUser(mGroupUserbelonging.UserID).User,
			Duration: handler.YearWithSemesterDuration{
				Since: handler.YearWithSemester{
					Year:     mGroupUserbelonging.SinceYear,
					Semester: handler.Semester(mGroupUserbelonging.SinceSemester),
				},
				Until: &handler.YearWithSemester{
					Year:     mGroupUserbelonging.UntilYear,
					Semester: handler.Semester(mGroupUserbelonging.UntilSemester),
				},
			},
		},
	}
}

func CloneHandlerMockProject() handler.ProjectDetail {
	var (
		mProject        = CloneMockProject()
		hProjectMembers = CloneHandlerMockProjectMembers()
	)

	return handler.ProjectDetail{
		Project: handler.Project{
			Id:   mProject.ID,
			Name: mProject.Name,
		},
		Description: mProject.Description,
		Link:        mProject.Link,
		Members:     hProjectMembers,
	}
}

func CloneHandlerMockProjectMembers() []handler.ProjectMember {
	var mProjectMember = CloneMockProjectMember()

	return []handler.ProjectMember{
		{
			User: getUser(mProjectMember.UserID).User,
			Duration: handler.YearWithSemesterDuration{
				Since: handler.YearWithSemester{
					Year:     mProjectMember.Project.SinceYear,
					Semester: handler.Semester(mProjectMember.SinceSemester),
				},
				Until: &handler.YearWithSemester{
					Year:     mProjectMember.UntilYear,
					Semester: handler.Semester(mProjectMember.UntilSemester),
				},
			},
		},
	}
}

func CloneHandlerMockUsers() []handler.UserDetail {
	var (
		mUsers      = CloneMockUsers()
		portalUsers = CloneMockPortalUsers()
		traqUsers   = CloneMockTraQUsers()
		hAccount    = CloneHandlerMockUserAccount()
		hUsers      = make([]handler.UserDetail, len(mUsers))
	)

	for i, mu := range mUsers {
		hUsers[i] = handler.UserDetail{
			User: handler.User{
				Id:       mu.ID,
				Name:     mu.Name,
				RealName: portalUsers[i].RealName,
			},
			Accounts: []handler.Account{},
			Bio:      mu.Description,
			State:    handler.UserAccountState(traqUsers[i].User.State),
		}

		if i == 0 {
			hUsers[i].Accounts = append(hUsers[i].Accounts, hAccount)
		}
	}

	return hUsers
}

func CloneHandlerMockUserAccount() handler.Account {
	var mAccount = CloneMockAccount()

	return handler.Account{
		DisplayName: mAccount.Name,
		Id:          mAccount.ID,
		PrPermitted: handler.PrPermitted(mAccount.Check),
		Type:        handler.AccountType(mAccount.Type),
		Url:         mAccount.URL,
	}
}

func CloneHandlerMockUserContest() handler.ContestTeamWithContestName {
	var (
		hContest     = CloneHandlerMockContest()
		hContestTeam = CloneHandlerMockContestTeam()
	)

	return handler.ContestTeamWithContestName{
		ContestTeam: hContestTeam.ContestTeam,
		ContestName: hContest.Name,
	}
}

func CloneHandlerMockUserGroup() handler.UserGroup {
	var (
		hGroup        = CloneHandlerMockGroup()
		hGroupMembers = CloneHandlerMockGroupMembers()
	)

	return handler.UserGroup{
		Group:    hGroup.Group,
		Duration: hGroupMembers[0].Duration,
	}
}

func CloneHandlerMockUserProject() handler.UserProject {
	var (
		hProject        = CloneHandlerMockProject()
		hProjectMembers = CloneHandlerMockProjectMembers()
	)

	return handler.UserProject{
		Project:      hProject.Project,
		UserDuration: hProjectMembers[0].Duration,
	}
}

func getUser(userID uuid.UUID) handler.UserDetail {
	var hUsers = CloneHandlerMockUsers()

	for _, hUser := range hUsers {
		if hUser.User.Id == userID {
			return hUser
		}
	}

	return handler.UserDetail{}
}
