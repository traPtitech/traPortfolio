package mockdata

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/infrastructure/external"
)

type TraQUser struct {
	User *external.TraQUserResponse
	Name string
}

var (
	MockKnoqEvents  = CloneMockKnoqEvents()
	MockPortalUsers = CloneMockPortalUsers()
	MockTraQUsers   = CloneMockTraQUsers()
)

func CloneMockKnoqEvents() []*external.EventResponse {
	return []*external.EventResponse{
		{
			ID:          KnoqEventID1(),
			Name:        "第n回進捗回",
			Description: "第n回の進捗会です。",
			Place:       "S516",
			GroupID:     KnoqEventGroupID1(),
			RoomID:      KnoqEventRoomID1(),
			TimeStart:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			TimeEnd:     time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			SharedRoom:  true,
			Admins: []uuid.UUID{
				UserID1(),
			},
		},
		{
			ID:          KnoqEventID2(),
			Name:        "sample event",
			Description: "This is a sample event.",
			Place:       "S516",
			GroupID:     KnoqEventGroupID2(),
			RoomID:      KnoqEventRoomID2(),
			TimeStart:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			TimeEnd:     time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			SharedRoom:  false,
			Admins: []uuid.UUID{
				UserID1(),
				UserID2(),
				UserID3(),
			},
		},
		{
			ID:          KnoqEventID3(),
			Name:        "sample event",
			Description: "This is a sample event.",
			Place:       "S516",
			GroupID:     KnoqEventGroupID3(),
			RoomID:      KnoqEventRoomID3(),
			TimeStart:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			TimeEnd:     time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			SharedRoom:  false,
			Admins: []uuid.UUID{
				UserID2(),
			},
		},
	}
}

func CloneMockPortalUsers() []*external.PortalUserResponse {
	return []*external.PortalUserResponse{
		{
			TraQID:         userName1,
			RealName:       userRealname1,
			AlphabeticName: "user1 user1",
		},
		{
			TraQID:         userName2,
			RealName:       userRealname2,
			AlphabeticName: "user2 user2",
		},
		{
			TraQID:         userName3,
			RealName:       userRealname3,
			AlphabeticName: "Noriko Azuma",
		},
	}
}

func CloneMockTraQUsers() []*TraQUser {
	return []*TraQUser{
		{
			User: &external.TraQUserResponse{
				ID:    UserID1(),
				State: domain.TraqStateActive,
			},
			Name: userName1,
		},
		{
			User: &external.TraQUserResponse{
				ID:    UserID2(),
				State: domain.TraqStateDeactivated,
			},
			Name: userName2,
		},
		{
			User: &external.TraQUserResponse{
				ID:    UserID3(),
				State: domain.TraqStateActive,
			},
			Name: userName3,
		},
	}
}
