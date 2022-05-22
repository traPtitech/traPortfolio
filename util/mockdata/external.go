package mockdata

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
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
			ID:          knoqEventID1.uuid(),
			Name:        "第n回進捗回",
			Description: "第n回の進捗会です。",
			Place:       "S516",
			GroupID:     knoqEventGroupID1.uuid(),
			RoomID:      knoqEventRoomID1.uuid(),
			TimeStart:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			TimeEnd:     time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			SharedRoom:  true,
			Admins: []uuid.UUID{
				userID1.uuid(),
			},
		},
		{
			ID:          knoqEventID2.uuid(),
			Name:        "sample event",
			Description: "This is a sample event.",
			Place:       "S516",
			GroupID:     knoqEventGroupID2.uuid(),
			RoomID:      knoqEventRoomID2.uuid(),
			TimeStart:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			TimeEnd:     time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			SharedRoom:  false,
			Admins: []uuid.UUID{
				userID1.uuid(),
				userID2.uuid(),
				userID3.uuid(),
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
				ID:    userID1.uuid(),
				State: domain.TraqStateActive,
			},
			Name: userName1,
		},
		{
			User: &external.TraQUserResponse{
				ID:    userID2.uuid(),
				State: domain.TraqStateDeactivated,
			},
			Name: userName2,
		},
		{
			User: &external.TraQUserResponse{
				ID:    userID3.uuid(),
				State: domain.TraqStateActive,
			},
			Name: userName3,
		},
	}
}
