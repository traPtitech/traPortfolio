package mockdata

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
)

var (
	HMockContest        = CloneHandlerMockContest()
	HMockContestTeam    = CloneHandlerMockContestTeam()
	HMockEvents         = CloneHandlerMockEvents()
	HMockEventDetails   = CloneHandlerMockEventDetails()
	HMockGroup          = CloneHandlerMockGroup()
	HMockGroupMembers   = CloneHandlerMockGroupMembers()
	HMockProjects       = CloneHandlerMockProjects()
	HMockProject        = CloneHandlerMockProject()
	HMockProjectMembers = CloneHandlerMockProjectMembers()
	HMockUsers          = CloneHandlerMockUsers()
	HMockUserDetails    = CloneHandlerMockUserDetails()
	HMockUserAccounts   = CloneHandlerMockUserAccounts()
	HMockUserEvents     = CloneHandlerMockUserEvents()
	HMockUserContests   = CloneHandlerMockUserContests()
	HMockUserGroups     = CloneHandlerMockUserGroups()
	HMockUserProjects   = CloneHandlerMockUserProjects()
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
		Id:          mContestTeam.ID,
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

func CloneHandlerMockEvents() []handler.Event {
	var (
		knoqEvents = CloneMockKnoqEvents()
		hEvents    = make([]handler.Event, len(knoqEvents))
	)

	for i, e := range knoqEvents {
		var (
			hostname = make([]handler.User, len(e.Admins))
		)

		for j, uid := range e.Admins {
			hostname[j] = getUser(uid)
		}

		hEvents[i] = handler.Event{
			Duration: handler.Duration{
				Since: e.TimeStart,
				Until: &e.TimeEnd,
			},
			Id:   e.ID,
			Name: e.Name,
		}
	}

	return hEvents
}

func CloneHandlerMockEventDetails() []handler.EventDetail {
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
		mGroupAdmins  = CloneMockGroupUserAdmin()
		mAdmin        = make([]handler.User, len(mGroupAdmins))
	)

	for i, adm := range mGroupAdmins {
		mAdmin[i] = getUser(adm.UserID)
	}

	return handler.GroupDetail{
		Description: mGroup.Description,
		Id:          mGroup.GroupID,
		Admin:       mAdmin,
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

func CloneHandlerMockProjects() []handler.Project {
	var (
		mProjects = CloneMockProjects()
		hProjects = make([]handler.Project, len(mProjects))
	)

	for i, p := range mProjects {
		hProjects[i] = handler.Project{
			Id:   p.ID,
			Name: p.Name,
			Duration: handler.YearWithSemesterDuration{
				Since: handler.YearWithSemester{
					Year:     p.SinceYear,
					Semester: handler.Semester(p.SinceSemester),
				},
				Until: &handler.YearWithSemester{
					Year:     p.UntilYear,
					Semester: handler.Semester(p.UntilSemester),
				},
			},
		}
	}

	return hProjects
}

func CloneHandlerMockProject() handler.ProjectDetail {
	var (
		mProject        = CloneMockProjects()[0]
		hProjectMembers = CloneHandlerMockProjectMembers()[0:2]
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
		mProjectMembers = CloneMockProjectMembers()
		hProjectMembers = make([]handler.ProjectMember, len(mProjectMembers))
	)

	for i, pm := range mProjectMembers {
		hUser := getUser(pm.UserID)
		hProjectMembers[i] = handler.ProjectMember{
			Duration: handler.YearWithSemesterDuration{
				Since: handler.YearWithSemester{
					Year:     pm.SinceYear,
					Semester: handler.Semester(pm.SinceSemester),
				},
				Until: &handler.YearWithSemester{
					Year:     pm.UntilYear,
					Semester: handler.Semester(pm.UntilSemester),
				},
			},
			Id:       hUser.Id,
			Name:     hUser.Name,
			RealName: hUser.RealName,
		}
	}

	return hProjectMembers
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
		hAccounts   = CloneHandlerMockUserAccounts()
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

		if mu.ID == userID1.uuid() {
			hUsers[i].Accounts = hAccounts
		}
	}

	return hUsers
}

func CloneHandlerMockUserAccounts() []handler.Account {
	var mAccount = CloneMockAccount()

	return []handler.Account{
		{
			DisplayName: mAccount.Name,
			Id:          mAccount.ID,
			PrPermitted: handler.PrPermitted(mAccount.Check),
			Type:        handler.AccountType(mAccount.Type),
			Url:         mAccount.URL,
		},
	}
}

func CloneHandlerMockUserEvents() []handler.Event {
	var (
		hEventDetails = CloneHandlerMockEventDetails()
		mUserEvents   = make([]handler.Event, len(hEventDetails))
	)

	for i, e := range hEventDetails {
		mUserEvents[i] = handler.Event{
			Duration: e.Duration,
			Id:       e.Id,
			Name:     e.Name,
		}
	}

	return mUserEvents
}

func CloneHandlerMockUserContests() []handler.ContestTeamWithContestName {
	var (
		hContest     = CloneHandlerMockContest()
		hContestTeam = CloneHandlerMockContestTeam()
	)

	return []handler.ContestTeamWithContestName{
		{
			ContestName: hContest.Name,
			Id:          hContestTeam.Id,
			Name:        hContestTeam.Name,
			Result:      hContestTeam.Result,
		},
	}
}

func CloneHandlerMockUserGroups() []handler.UserGroup {
	var (
		hGroup        = CloneHandlerMockGroup()
		hGroupMembers = CloneHandlerMockGroupMembers()
		hUserGroups   = make([]handler.UserGroup, len(hGroupMembers))
	)

	for i, gm := range hGroupMembers {
		hUserGroups[i] = handler.UserGroup{
			Duration: gm.Duration,
			Id:       hGroup.Id,
			Name:     hGroup.Name,
		}
	}

	return hUserGroups
}

func CloneHandlerMockUserProjects() []handler.UserProject {
	var (
		hProject        = CloneHandlerMockProject()
		hProjectMembers = CloneHandlerMockProjectMembers()
		hUserProjects   = make([]handler.UserProject, len(hProjectMembers))
	)

	for i, pm := range hProjectMembers {
		hUserProjects[i] = handler.UserProject{
			Duration:     hProject.Duration,
			Id:           hProject.Id,
			Name:         hProject.Name,
			UserDuration: pm.Duration,
		}
	}

	return hUserProjects
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
