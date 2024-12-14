package mockdata

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/handler/schema"
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

func CloneHandlerMockContestDetails() []schema.ContestDetail {
	var (
		mContests        = CloneMockContests()
		hContestTeams    = CloneMockContestTeams()
		hTeamMembersByID = CloneHandlerMockContestTeamMembersByID()
		mContestTeams    = make([]schema.ContestTeam, len(hContestTeams))
		hContestDetails  = make([]schema.ContestDetail, len(mContests))
	)

	for i, c := range hContestTeams {
		members, ok := hTeamMembersByID[c.ID]
		if !ok {
			members = make([]schema.User, 0)
		}
		mContestTeams[i] = schema.ContestTeam{
			Id:      c.ID,
			Members: members,
			Name:    c.Name,
			Result:  c.Result,
		}
	}

	for i, c := range mContests {
		hContestDetails[i] = schema.ContestDetail{
			Description: c.Description,
			Duration: schema.Duration{
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

func CloneHandlerMockContests() []schema.Contest {
	var (
		mContests = CloneMockContests()
		hContests = make([]schema.Contest, len(mContests))
	)

	for i, c := range mContests {
		hContests[i] = schema.Contest{
			Duration: schema.Duration{
				Since: c.Since,
				Until: &c.Until,
			},
			Id:   c.ID,
			Name: c.Name,
		}
	}
	return hContests
}

func CloneHandlerMockContestTeamsByID() map[uuid.UUID][]schema.ContestTeam {
	var (
		mContestTeams     = CloneMockContestTeams()
		hTeamMembersByID  = CloneHandlerMockContestTeamMembersByID()
		hContestTeamsByID = make(map[uuid.UUID][]schema.ContestTeam)
	)

	for _, ct := range mContestTeams {
		members, ok := hTeamMembersByID[ct.ID]
		if !ok {
			members = make([]schema.User, 0)
		}
		hContestTeamsByID[ct.ContestID] = append(hContestTeamsByID[ct.ContestID], schema.ContestTeam{
			Id:      ct.ID,
			Members: members,
			Name:    ct.Name,
			Result:  ct.Result,
		})
	}
	return hContestTeamsByID
}

func CloneHandlerMockContestTeamMembersByID() map[uuid.UUID][]schema.User {
	var (
		hContestTeams         = CloneMockContestTeams()
		mockMembersBelongings = CloneMockContestTeamUserBelongings()
		mockMembers           = CloneMockUsers()
		portalUsers           = CloneMockPortalUsers()
		hContestMembers       = make(map[uuid.UUID][]schema.User, len(hContestTeams))
	)

	for _, c := range hContestTeams {
		for _, ct := range mockMembersBelongings {
			if c.ID == ct.TeamID {
				for i, cm := range mockMembers {
					if ct.UserID == cm.ID {
						hContestMembers[ct.TeamID] = append(hContestMembers[c.ID], schema.User{
							Id:       cm.ID,
							Name:     cm.Name,
							RealName: portalUsers[i].RealName,
						})
					}
				}
			}
		}
	}

	return hContestMembers
}

func CloneHandlerMockEvents() []schema.Event {
	var (
		eventDetails = CloneHandlerMockEventDetails()
		events       = make([]schema.Event, len(eventDetails))
	)

	for i, e := range eventDetails {
		events[i] = schema.Event{
			Duration: e.Duration,
			Id:       e.Id,
			Name:     e.Name,
			Level:    e.Level,
		}
	}
	return events
}

func CloneHandlerMockEventDetails() []schema.EventDetail {
	var (
		mEventLevels = CloneMockEventLevelRelations()
		knoqEvents   = CloneMockKnoqEvents()
		hEvents      = make([]schema.EventDetail, 0, len(knoqEvents))
	)

	for _, e := range knoqEvents {
		var (
			eventLevel schema.EventLevel
			hostname   = make([]schema.User, len(e.Admins))
		)

		for _, l := range mEventLevels {
			if l.ID == e.ID {
				eventLevel = schema.EventLevel(l.Level)
				break
			}
		}

		for j, uid := range e.Admins {
			hostname[j] = getUser(uid)
		}

		event := schema.EventDetail{
			Description: e.Description,
			Duration: schema.Duration{
				Since: e.TimeStart,
				Until: &e.TimeEnd,
			},
			Level:    eventLevel,
			Hostname: hostname,
			Id:       e.ID,
			Name:     e.Name,
			Place:    e.Place,
		}
		switch eventLevel {
		case schema.EventLevel(domain.EventLevelPrivate):
			continue
		case schema.EventLevel(domain.EventLevelPublic):
			hEvents = append(hEvents, event)
		case schema.EventLevel(domain.EventLevelAnonymous):
			event.Hostname = nil
			hEvents = append(hEvents, event)
		default:
			panic("invalid event level")
		}
	}

	return hEvents
}

func CloneHandlerMockGroups() []schema.GroupDetail {
	var (
		mGroups        = CloneMockGroups()
		hGroupsMembers = CloneHandlerMockGroupMembersByID()
		mGroupAdmins   = CloneMockGroupUserAdmins()
		hGroups        = make([]schema.GroupDetail, len(mGroups))
	)

	for i, g := range mGroups {
		mAdmins := make([]schema.User, 0, len(mGroupAdmins))
		for _, adm := range mGroupAdmins {
			if g.GroupID == adm.GroupID {
				mAdmins = append(mAdmins, getUser(adm.UserID))
			}
		}
		hGroups[i] = schema.GroupDetail{
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

func CloneHandlerMockGroupMembersByID() map[uuid.UUID][]schema.GroupMember {
	var (
		mGroups              = CloneMockGroups()
		mGroupUserbelongings = CloneMockGroupUserBelongings()
		hGroupMembers        = make(map[uuid.UUID][]schema.GroupMember, len(mGroups))
	)

	for _, gub := range mGroupUserbelongings {
		for _, g := range mGroups {
			if gub.GroupID == g.GroupID {
				hUser := getUser(gub.UserID)
				hGroupMembers[g.GroupID] = append(hGroupMembers[g.GroupID],
					schema.GroupMember{
						Duration: schema.YearWithSemesterDuration{
							Since: schema.YearWithSemester{
								Year:     gub.SinceYear,
								Semester: schema.Semester(gub.SinceSemester),
							},
							Until: &schema.YearWithSemester{
								Year:     gub.UntilYear,
								Semester: schema.Semester(gub.UntilSemester),
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

func CloneHandlerMockProjects() []schema.Project {
	var (
		mProjects = CloneMockProjects()
		hProjects = make([]schema.Project, len(mProjects))
	)

	for i, p := range mProjects {
		hProjects[i] = schema.Project{
			Id:   p.ID,
			Name: p.Name,
			Duration: schema.YearWithSemesterDuration{
				Since: schema.YearWithSemester{
					Year:     p.SinceYear,
					Semester: schema.Semester(p.SinceSemester),
				},
				Until: &schema.YearWithSemester{
					Year:     p.UntilYear,
					Semester: schema.Semester(p.UntilSemester),
				},
			},
		}
	}

	return hProjects
}

func CloneHandlerMockProjectDetails() []schema.ProjectDetail {
	var (
		mProjects       = CloneMockProjects()
		hProjectMembers = CloneHandlerMockProjectMembers()
		mProjectMembers = CloneMockProjectMembers()
		hProjects       = make([]schema.ProjectDetail, len(mProjects))
	)

	for i, mp := range mProjects {
		hProjects[i] = schema.ProjectDetail{
			Description: mp.Description,
			Duration: schema.YearWithSemesterDuration{
				Since: schema.YearWithSemester{
					Year:     mp.SinceYear,
					Semester: schema.Semester(mp.SinceSemester),
				},
				Until: &schema.YearWithSemester{
					Year:     mp.UntilYear,
					Semester: schema.Semester(mp.UntilSemester),
				},
			},
			Id:      mp.ID,
			Link:    mp.Link,
			Members: []schema.ProjectMember{},
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

func CloneHandlerMockProjectMembers() []schema.ProjectMember {
	var (
		mProjectMembers = CloneMockProjectMembers()
		hProjectMembers = make([]schema.ProjectMember, len(mProjectMembers))
	)

	for i, pm := range mProjectMembers {
		hUser := getUser(pm.UserID)
		hProjectMembers[i] = schema.ProjectMember{
			Duration: schema.YearWithSemesterDuration{
				Since: schema.YearWithSemester{
					Year:     pm.SinceYear,
					Semester: schema.Semester(pm.SinceSemester),
				},
				Until: &schema.YearWithSemester{
					Year:     pm.UntilYear,
					Semester: schema.Semester(pm.UntilSemester),
				},
			},
			Id:       hUser.Id,
			Name:     hUser.Name,
			RealName: hUser.RealName,
		}
	}

	return hProjectMembers
}

func CloneHandlerMockUsers() []schema.User {
	var (
		mUsers      = CloneMockUsers()
		portalUsers = CloneMockPortalUsers()
		hUsers      = make([]schema.User, len(mUsers))
	)

	for i, u := range mUsers {
		d := *domain.NewUser(u.ID, u.Name, portalUsers[i].RealName, u.Check)
		hUsers[i] = schema.User{
			Id:       d.ID,
			Name:     d.Name,
			RealName: d.RealName(),
		}
	}

	return hUsers
}

func CloneHandlerMockUserDetails() []schema.UserDetail {
	var (
		mUsers      = CloneMockUsers()
		portalUsers = CloneMockPortalUsers()
		hAccounts   = CloneHandlerMockUserAccountsByID()
		hUsers      = make([]schema.UserDetail, len(mUsers))
	)

	for i, mu := range mUsers {
		hUsers[i] = schema.UserDetail{
			Accounts: hAccounts[mu.ID],
			Bio:      mu.Description,
			Id:       mu.ID,
			Name:     mu.Name,
			RealName: portalUsers[i].RealName,
			State:    schema.UserAccountState(mu.State),
		}
	}

	return hUsers
}

func CloneHandlerMockUserAccountsByID() map[uuid.UUID][]schema.Account {
	var (
		mAccounts = CloneMockAccounts()
		hAccounts = make(map[uuid.UUID][]schema.Account)
	)

	for _, a := range mAccounts {
		hAccounts[a.UserID] = append(hAccounts[a.UserID], schema.Account{
			DisplayName: a.Name,
			Id:          a.ID,
			Type:        schema.AccountType(a.Type),
			Url:         a.URL,
		})
	}

	return hAccounts
}

func CloneHandlerMockUserEvents() []schema.Event {
	var (
		hEventDetails = CloneHandlerMockEventDetails()
		mUserEvents   = make([]schema.Event, len(hEventDetails))
	)

	for i, e := range hEventDetails {
		mUserEvents[i] = schema.Event{
			Duration: e.Duration,
			Id:       e.Id,
			Name:     e.Name,
			Level:    e.Level,
		}
	}

	return mUserEvents
}

// userが所属するteamが参加したcontestの一覧を返す

func CloneHandlerMockUserContestsByID() map[uuid.UUID][]schema.UserContest {
	var (
		mUsers              = CloneMockUsers()
		hContests           = CloneHandlerMockContests()
		mContestTeams       = CloneMockContestTeams()
		mContestTeamMembers = CloneMockContestTeamUserBelongings()
		hUserContests       = make(map[uuid.UUID][]schema.UserContest, len(hContests))
	)

	// userContestTeams[userID][contestID] = []ContestTeam
	userContestTeams := make(map[uuid.UUID]map[uuid.UUID][]schema.ContestTeamWithoutMembers, len(mUsers))
	for _, u := range mUsers {
		userContestTeams[u.ID] = make(map[uuid.UUID][]schema.ContestTeamWithoutMembers, len(hContests))
		for _, c := range hContests {
			userContestTeams[u.ID][c.Id] = []schema.ContestTeamWithoutMembers{}
			for _, ct := range mContestTeams {
				if ct.ContestID != c.Id {
					continue
				}

				for _, ctm := range mContestTeamMembers {
					if ctm.TeamID != ct.ID || ctm.UserID != u.ID {
						continue
					}

					userContestTeams[u.ID][c.Id] = append(userContestTeams[u.ID][c.Id], schema.ContestTeamWithoutMembers{
						Id:     ct.ID,
						Name:   ct.Name,
						Result: ct.Result,
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

			hUserContests[u.ID] = append(hUserContests[u.ID], schema.UserContest{
				Duration: c.Duration,
				Id:       c.Id,
				Name:     c.Name,
				Teams:    userContestTeams[u.ID][c.Id],
			})
		}
	}

	return hUserContests
}

func CloneHandlerMockUserGroupsByID() map[uuid.UUID][]schema.UserGroup {
	var (
		hUsers        = CloneHandlerMockUsers()
		hGroups       = CloneHandlerMockGroups()
		hGroupMembers = CloneHandlerMockGroupMembersByID()
		hUserGroups   = make(map[uuid.UUID][]schema.UserGroup, len(hUsers))
	)

	for _, u := range hUsers {
		for _, g := range hGroups {
			for _, gm := range hGroupMembers[g.Id] {
				if u.Id == gm.Id {
					hUserGroups[u.Id] = append(hUserGroups[u.Id], schema.UserGroup{
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

func CloneHandlerMockUserProjects() []schema.UserProject {
	var (
		hProject        = CloneHandlerMockProjects()[0]
		hProjectMembers = CloneHandlerMockProjectMembers()
		hUserProjects   = make([]schema.UserProject, len(hProjectMembers))
	)

	for i, pm := range hProjectMembers {
		hUserProjects[i] = schema.UserProject{
			Duration:     hProject.Duration,
			Id:           hProject.Id,
			Name:         hProject.Name,
			UserDuration: pm.Duration,
		}
	}

	return hUserProjects
}

func getUser(userID uuid.UUID) schema.User {
	var hUsers = CloneHandlerMockUsers()

	for _, hUser := range hUsers {
		if hUser.Id == userID {
			return hUser
		}
	}

	return schema.User{}
}

func AccountTypesMockUserHas(userID uuid.UUID) []schema.AccountType {
	var (
		mAccounts = CloneHandlerMockUserAccountsByID()[userID]
	)

	holdAccounts := []schema.AccountType{}
	for _, account := range mAccounts {
		holdAccounts = append(holdAccounts, schema.AccountType(account.Type))
	}

	return holdAccounts
}

func AccountTypesMockUserDoesntHave(userID uuid.UUID) []schema.AccountType {
	holdAccounts := AccountTypesMockUserHas(userID)
	vacantAccounts := []schema.AccountType{}

	holdAccountsMap := make(map[schema.AccountType]struct{})
	for _, account := range holdAccounts {
		holdAccountsMap[account] = struct{}{}
	}

	for i := range schema.AccountType(domain.AccountLimit) {
		if _, ok := holdAccountsMap[i]; !ok {
			vacantAccounts = append(vacantAccounts, i)
		}
	}

	return vacantAccounts
}
