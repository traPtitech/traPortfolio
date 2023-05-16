package mockdata

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
)

var (
	HMockContestDetails         = CloneHandlerMockContestDetails()
	HMockContests               = CloneHandlerMockContests()
	HMockContestTeamsByID       = CloneHandlerMockContestTeamsByID()
	HMockContestTeamMembersByID = CloneHandlerMockContestTeamMembersByID()
	HMockEvents                 = CloneHandlerMockEvents()
	HMockEventDetails           = CloneHandlerMockEventDetails()
	HMockGroups                 = CloneHandlerMockGroups()
	HMockGroupMembersByID       = CloneHandlerMockGroupMembersByID()
	HMockProjects               = CloneHandlerMockProjects()
	HMockProjectDetails         = CloneHandlerMockProjectDetails()
	HMockProjectMembers         = CloneHandlerMockProjectMembers()
	HMockUsers                  = CloneHandlerMockUsers()
	HMockUserDetails            = CloneHandlerMockUserDetails()
	HMockUserAccountsByID       = CloneHandlerMockUserAccountsByID()
	HMockUserEvents             = CloneHandlerMockUserEvents()
	HMockUserContestsByID       = CloneHandlerMockUserContestsByID()
	HMockUserGroupsByID         = CloneHandlerMockUserGroupsByID()
	HMockUserProjects           = CloneHandlerMockUserProjects()
)

func CloneHandlerMockContestDetails() []handler.ContestDetail {
	var (
		mContests        = CloneMockContests()
		hContestTeams    = CloneMockContestTeams()
		hTeamMembersByID = CloneHandlerMockContestTeamMembersByID()
		mContestTeams    = make([]handler.ContestTeam, len(hContestTeams))
		hContestDetails  = make([]handler.ContestDetail, len(mContests))
	)

	for i, c := range hContestTeams {
		members := hTeamMembersByID[c.ID]
		if members == nil {
			members = make([]handler.User, 0)
		}
		mContestTeams[i] = handler.ContestTeam{
			Id:      c.ID,
			Members: members,
			Name:    c.Name,
			Result:  c.Result,
		}
	}

	for i, c := range mContests {
		hContestDetails[i] = handler.ContestDetail{
			Description: c.Description,
			Duration: handler.Duration{
				Since: c.Since,
				Until: &c.Until,
			},
			Id:    c.ID,
			Link:  c.Link,
			Name:  c.Name,
			Teams: mContestTeams,
		}
	}
	return hContestDetails
}

func CloneHandlerMockContests() []handler.Contest {
	var (
		mContests = CloneMockContests()
		hContests = make([]handler.Contest, len(mContests))
	)

	for i, c := range mContests {
		hContests[i] = handler.Contest{
			Duration: handler.Duration{
				Since: c.Since,
				Until: &c.Until,
			},
			Id:   c.ID,
			Name: c.Name,
		}
	}
	return hContests
}

func CloneHandlerMockContestTeamsByID() map[uuid.UUID][]handler.ContestTeam {
	var (
		mContestTeams     = CloneMockContestTeams()
		hTeamMembersByID  = CloneHandlerMockContestTeamMembersByID()
		hContestTeamsByID = make(map[uuid.UUID][]handler.ContestTeam)
	)

	for _, ct := range mContestTeams {
		members := hTeamMembersByID[ct.ID]
		if members == nil {
			members = make([]handler.User, 0)
		}
		hContestTeamsByID[ct.ContestID] = append(hContestTeamsByID[ct.ContestID], handler.ContestTeam{
			Id:      ct.ID,
			Members: members,
			Name:    ct.Name,
			Result:  ct.Result,
		})
	}
	return hContestTeamsByID
}

func CloneHandlerMockContestTeamMembersByID() map[uuid.UUID][]handler.User {

	var (
		hContestTeams         = CloneMockContestTeams()
		mockMembersBelongings = CloneMockContestTeamUserBelongings()
		mockMembers           = CloneMockUsers()
		hContestMembers       = make(map[uuid.UUID][]handler.User, len(hContestTeams))
	)

	for _, c := range hContestTeams {
		for _, ct := range mockMembersBelongings {
			if c.ID == ct.TeamID {

				for _, cm := range mockMembers {
					if ct.UserID == cm.ID {
						hContestMembers[ct.TeamID] = append(hContestMembers[c.ID], handler.User{
							Id:   cm.ID,
							Name: cm.Name,
						})
					}
				}

			}
		}
	}

	return hContestMembers
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
		d := *domain.NewUser(u.ID, u.Name, portalUsers[i].RealName, u.Check)
		hUsers[i] = handler.User{
			Id:       d.ID,
			Name:     d.Name,
			RealName: d.RealName(),
		}
	}

	return hUsers
}

func CloneHandlerMockUserDetails() []handler.UserDetail {
	var (
		mUsers      = CloneMockUsers()
		portalUsers = CloneMockPortalUsers()
		traqUsers   = CloneMockTraQUsers()
		hAccounts   = CloneHandlerMockUserAccountsByID()
		hUsers      = make([]handler.UserDetail, len(mUsers))
	)

	for i, mu := range mUsers {
		hUsers[i] = handler.UserDetail{
			Accounts: hAccounts[mu.ID],
			Bio:      mu.Description,
			Id:       mu.ID,
			Name:     mu.Name,
			RealName: portalUsers[i].RealName,
			State:    handler.UserAccountState(traqUsers[i].User.State),
		}
	}

	return hUsers
}

func CloneHandlerMockUserAccountsByID() map[uuid.UUID][]handler.Account {
	var (
		mAccounts = CloneMockAccounts()
		hAccounts = make(map[uuid.UUID][]handler.Account)
	)

	for _, a := range mAccounts {
		hAccounts[a.UserID] = append(hAccounts[a.UserID], handler.Account{
			DisplayName: a.Name,
			Id:          a.ID,
			PrPermitted: a.Check,
			Type:        handler.AccountType(a.Type),
			Url:         a.URL,
		})
	}

	return hAccounts
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

// userが所属するteamが参加したcontestの一覧を返す

func CloneHandlerMockUserContestsByID() map[uuid.UUID][]handler.UserContest {
	var (
		mUsers              = CloneMockUsers()
		hContests           = CloneHandlerMockContests()
		mContestTeams       = CloneMockContestTeams()
		mContestTeamMembers = CloneMockContestTeamUserBelongings()
		hTeamMembersByID    = CloneHandlerMockContestTeamMembersByID()
		hUserContests       = make(map[uuid.UUID][]handler.UserContest, len(hContests))
	)

	// userContestTeams[userID][contestID] = []ContestTeam
	userContestTeams := make(map[uuid.UUID]map[uuid.UUID][]handler.ContestTeam, len(mUsers))
	for _, u := range mUsers {
		userContestTeams[u.ID] = make(map[uuid.UUID][]handler.ContestTeam, len(hContests))
		for _, c := range hContests {
			userContestTeams[u.ID][c.Id] = []handler.ContestTeam{}
			for _, ct := range mContestTeams {
				if ct.ContestID != c.Id {
					continue
				}

				for _, ctm := range mContestTeamMembers {
					if ctm.TeamID != ct.ID || ctm.UserID != u.ID {
						continue
					}

					userContestTeams[u.ID][c.Id] = append(userContestTeams[u.ID][c.Id], handler.ContestTeam{
						Id:      ct.ID,
						Members: hTeamMembersByID[ct.ID],
						Name:    ct.Name,
						Result:  ct.Result,
					})
				}
			}
		}
	}

	for _, u := range mUsers {
		for _, c := range hContests {
			if len(userContestTeams[u.ID][c.Id]) == 0 {
				continue
			}

			hUserContests[u.ID] = append(hUserContests[u.ID], handler.UserContest{
				Duration: c.Duration,
				Id:       c.Id,
				Name:     c.Name,
				Teams:    userContestTeams[u.ID][c.Id],
			})
		}
	}

	return hUserContests
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
