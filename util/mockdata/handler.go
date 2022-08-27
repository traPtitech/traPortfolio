package mockdata

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
)

var (
	HMockContest          = CloneHandlerMockContest()
	HMockContests         = CloneHandlerMockContests()
	HMockContestTeam      = CloneHandlerMockContestTeam()
	HMockEvents           = CloneHandlerMockEvents()
	HMockEventDetails     = CloneHandlerMockEventDetails()
	HMockGroups           = CloneHandlerMockGroups()
	HMockGroupMembersByID = CloneHandlerMockGroupMembersByID()
	HMockProjects         = CloneHandlerMockProjects()
	HMockProjectDetails   = CloneHandlerMockProjectDetails()
	HMockProjectMembers   = CloneHandlerMockProjectMembers()
	HMockUsers            = CloneHandlerMockUsers()
	HMockUserDetails      = CloneHandlerMockUserDetails()
	HMockUserAccounts     = CloneHandlerMockUserAccounts()
	HMockUserEvents       = CloneHandlerMockUserEvents()
	HMockUserContests     = CloneHandlerMockUserContests()
	HMockUserGroupsByID   = CloneHandlerMockUserGroupsByID()
	HMockUserProjects     = CloneHandlerMockUserProjects()
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

func CloneHandlerMockContests() []handler.Contest {
	var (
		mContest = CloneMockContest()
	)

	return []handler.Contest{
		{
			Duration: handler.Duration{
				Since: mContest.Since,
				Until: &mContest.Until,
			},
			Id:   mContest.ID,
			Name: mContest.Name,
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

func CloneHandlerMockGroups() []handler.GroupDetail {
	var (
		mGroups        = CloneMockGroups()
		hGroupsMembers = CloneHandlerMockGroupMembersByID()
		mGroupAdmins   = CloneMockGroupUserAdmins()
		hGroups        = make([]handler.GroupDetail, len(mGroups))
	)

	for i, g := range mGroups {
		mAdmins := make([]handler.User, 0, len(mGroupAdmins))
		for _, adm := range mGroupAdmins {
			if g.GroupID == adm.GroupID {
				mAdmins = append(mAdmins, getUser(adm.UserID))
			}
		}
		hGroups[i] = handler.GroupDetail{
			Description: g.Description,
			Id:          g.GroupID,
			Admin:       mAdmins,
			Link:        g.Link,
			Members:     hGroupsMembers[g.GroupID],
			Name:        g.Name,
		}
	}

	return hGroups
}

func CloneHandlerMockGroupMembersByID() map[uuid.UUID][]handler.GroupMember {
	var (
		mGroups              = CloneMockGroups()
		mGroupUserbelongings = CloneMockGroupUserBelongings()
		hGroupMembers        = make(map[uuid.UUID][]handler.GroupMember, len(mGroups))
	)

	for _, gub := range mGroupUserbelongings {
		for _, g := range mGroups {
			if gub.GroupID == g.GroupID {
				hUser := getUser(gub.UserID)
				hGroupMembers[g.GroupID] = append(hGroupMembers[g.GroupID],
					handler.GroupMember{
						Duration: handler.YearWithSemesterDuration{
							Since: handler.YearWithSemester{
								Year:     gub.SinceYear,
								Semester: handler.Semester(gub.SinceSemester),
							},
							Until: &handler.YearWithSemester{
								Year:     gub.UntilYear,
								Semester: handler.Semester(gub.UntilSemester),
							},
						},
						Id:       hUser.Id,
						Name:     hUser.Name,
						RealName: hUser.RealName,
					})
			}
		}
	}
	return hGroupMembers
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

func CloneHandlerMockProjectDetails() []handler.ProjectDetail {
	var (
		mProjects       = CloneMockProjects()
		hProjectMembers = CloneHandlerMockProjectMembers()
		mProjectMembers = CloneMockProjectMembers()
		hProjects       = make([]handler.ProjectDetail, len(mProjects))
	)

	for i, mp := range mProjects {
		hProjects[i] = handler.ProjectDetail{
			Description: mp.Description,
			Duration: handler.YearWithSemesterDuration{
				Since: handler.YearWithSemester{
					Year:     mp.SinceYear,
					Semester: handler.Semester(mp.SinceSemester),
				},
				Until: &handler.YearWithSemester{
					Year:     mp.UntilYear,
					Semester: handler.Semester(mp.UntilSemester),
				},
			},
			Id:      mp.ID,
			Link:    mp.Link,
			Members: []handler.ProjectMember{},
			Name:    mp.Name,
		}
		for j, mpm := range mProjectMembers {
			if mpm.ProjectID == mp.ID {
				hProjects[i].Members = append(hProjects[i].Members, hProjectMembers[j])
			}
		}
	}
	return hProjects
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

func CloneHandlerMockUserGroupsByID() map[uuid.UUID][]handler.UserGroup {
	var (
		hUsers        = CloneHandlerMockUsers()
		hGroups       = CloneHandlerMockGroups()
		hGroupMembers = CloneHandlerMockGroupMembersByID()
		hUserGroups   = make(map[uuid.UUID][]handler.UserGroup, len(hUsers))
	)

	for _, u := range hUsers {
		for _, g := range hGroups {
			for _, gm := range hGroupMembers[g.Id] {
				if u.Id == gm.Id {
					hUserGroups[u.Id] = append(hUserGroups[u.Id], handler.UserGroup{
						Duration: gm.Duration,
						Id:       g.Id,
						Name:     g.Name,
					})
				}
			}
		}
	}
	return hUserGroups
}

func CloneHandlerMockUserProjects() []handler.UserProject {
	var (
		hProject        = CloneHandlerMockProjects()[0]
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
