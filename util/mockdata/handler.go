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
	HMockUserDetails    = CloneHandlerMockUserDetails()
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
		Description: mContest.Description,
		Duration: handler.Duration{
			Since: mContest.Since,
			Until: &mContest.Until,
		},
		Id:   mContest.ID,
		Link: mContest.Link,
		Name: mContest.Name,
		Teams: []handler.ContestTeam{
			{
				Id:     hContestTeam.Id,
				Name:   hContestTeam.Name,
				Result: hContestTeam.Result,
			},
		},
	}
}

func CloneHandlerMockContestTeam() handler.ContestTeamDetail {
	var (
		mContestTeam              = CloneMockContestTeam()
		mContestTeamUserBelonging = CloneMockContestTeamUserBelonging()
		hUser                     = getUser(mContestTeamUserBelonging.UserID)
	)

	return handler.ContestTeamDetail{
		Description: mContestTeam.Description,
		Id:          mContestTeam.ContestID,
		Link:        mContestTeam.Link,
		Members: []handler.User{
			{
				Id:       hUser.Id,
				Name:     hUser.Name,
				RealName: hUser.RealName,
			},
		},
		Name:   mContestTeam.Name,
		Result: mContestTeam.Result,
	}
}

func CloneHandlerMockEvents() []handler.EventDetail {
	var (
		mEventLevels = CloneMockEventLevelRelations()
		knoqEvents   = CloneMockKnoqEvents()
		hEvents      = make([]handler.EventDetail, len(knoqEvents))
	)

	for i, e := range knoqEvents {
		var (
			eventLevel handler.EventLevel
			hostname   = make([]handler.User, len(e.Admins))
		)

		for _, l := range mEventLevels {
			if l.ID == e.ID {
				eventLevel = handler.EventLevel(l.Level)
				break
			}
		}

		for j, uid := range e.Admins {
			hostname[j] = getUser(uid)
		}

		hEvents[i] = handler.EventDetail{
			Description: e.Description,
			Duration: handler.Duration{
				Since: e.TimeStart,
				Until: &e.TimeEnd,
			},
			EventLevel: eventLevel,
			Hostname:   hostname,
			Id:         e.ID,
			Name:       e.Name,
			Place:      e.Place,
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
		Description: mGroup.Description,
		Id:          mGroup.GroupID,
		Leader:      getUser(mGroup.Leader),
		Link:        mGroup.Link,
		Members:     hGroupMembers,
		Name:        mGroup.Name,
	}
}

func CloneHandlerMockGroupMembers() []handler.GroupMember {
	var (
		mGroupUserbelonging = CloneMockGroupUserBelonging()
		hUser               = getUser(mGroupUserbelonging.UserID)
	)

	return []handler.GroupMember{
		{
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
			Id:       hUser.Id,
			Name:     hUser.Name,
			RealName: hUser.RealName,
		},
	}
}

func CloneHandlerMockProject() handler.ProjectDetail {
	var (
		mProject        = CloneMockProject()
		hProjectMembers = CloneHandlerMockProjectMembers()
	)

	return handler.ProjectDetail{
		Description: mProject.Description,
		Duration: handler.YearWithSemesterDuration{
			Since: handler.YearWithSemester{
				Year:     mProject.SinceYear,
				Semester: handler.Semester(mProject.SinceSemester),
			},
			Until: &handler.YearWithSemester{
				Year:     mProject.UntilYear,
				Semester: handler.Semester(mProject.UntilSemester),
			},
		},
		Id:      mProject.ID,
		Link:    mProject.Link,
		Members: hProjectMembers,
		Name:    mProject.Name,
	}
}

func CloneHandlerMockProjectMembers() []handler.ProjectMember {
	var (
		mProjectMember = CloneMockProjectMember()
		hUser          = getUser(mProjectMember.UserID)
	)

	return []handler.ProjectMember{
		{
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
			Id:       hUser.Id,
			Name:     hUser.Name,
			RealName: hUser.RealName,
		},
	}
}

func CloneHandlerMockUsers() []handler.User {
	var (
		mUsers      = CloneMockUsers()
		portalUsers = CloneMockPortalUsers()
		hUsers      = make([]handler.User, len(mUsers))
	)

	for i, u := range mUsers {
		hUsers[i] = handler.User{
			Id:       u.ID,
			Name:     u.Name,
			RealName: portalUsers[i].RealName,
		}
	}

	return hUsers
}

func CloneHandlerMockUserDetails() []handler.UserDetail {
	var (
		mUsers      = CloneMockUsers()
		portalUsers = CloneMockPortalUsers()
		traqUsers   = CloneMockTraQUsers()
		hAccount    = CloneHandlerMockUserAccount()
		hUsers      = make([]handler.UserDetail, len(mUsers))
	)

	for i, mu := range mUsers {
		hUsers[i] = handler.UserDetail{
			Accounts: []handler.Account{},
			Bio:      mu.Description,
			Id:       mu.ID,
			Name:     mu.Name,
			RealName: portalUsers[i].RealName,
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
		ContestName: hContest.Name,
		Id:          hContestTeam.Id,
		Name:        hContestTeam.Name,
		Result:      hContestTeam.Result,
	}
}

func CloneHandlerMockUserGroup() handler.UserGroup {
	var (
		hGroup        = CloneHandlerMockGroup()
		hGroupMembers = CloneHandlerMockGroupMembers()
	)

	return handler.UserGroup{
		Duration: hGroupMembers[0].Duration,
		Id:       hGroup.Id,
		Name:     hGroup.Name,
	}
}

func CloneHandlerMockUserProject() handler.UserProject {
	var (
		hProject        = CloneHandlerMockProject()
		hProjectMembers = CloneHandlerMockProjectMembers()
	)

	return handler.UserProject{
		Duration:     hProjectMembers[0].Duration,
		Id:           hProject.Id,
		Name:         hProject.Name,
		UserDuration: hProjectMembers[0].Duration,
	}
}

func getUser(userID uuid.UUID) handler.User {
	var hUsers = CloneHandlerMockUsers()

	for _, hUser := range hUsers {
		if hUser.Id == userID {
			return hUser
		}
	}

	return handler.User{}
}
