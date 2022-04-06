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
	MockKnoqEvents = []*external.EventResponse{
		{
			ID:          uuid.FromStringOrNil("d1274c6e-15cc-4ca0-b720-1c03ea3a60ec"),
			Name:        "第n回進捗回",
			Description: "第n回の進捗会です。",
			Place:       "S516",
			GroupID:     uuid.FromStringOrNil("7ecabb2a-8e2c-4ebe-bb0b-13254a6eae05"),
			RoomID:      uuid.FromStringOrNil("68319c0c-be20-45c1-a05d-7651473bd396"),
			TimeStart:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			TimeEnd:     time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			SharedRoom:  true,
			Admins: []uuid.UUID{
				MockUsers[0].ID,
			},
		},
		{
			ID:          uuid.FromStringOrNil("e28ec610-226d-49c5-be7c-86af54f6839d"),
			Name:        "sample event",
			Description: "This is a sample event.",
			Place:       "S516",
			GroupID:     uuid.FromStringOrNil("9c592124-52a5-4981-a2c8-1e218c64a8e5"),
			RoomID:      uuid.FromStringOrNil("cbd48b1f-6b20-41c8-b122-a9826bd968ed"),
			TimeStart:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			TimeEnd:     time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
			SharedRoom:  false,
			Admins: []uuid.UUID{
				MockUsers[0].ID,
				MockUsers[1].ID,
				MockUsers[2].ID,
			},
		},
	}
	MockPortalUsers = []*external.PortalUserResponse{
		{
			TraQID:         MockUsers[0].Name,
			RealName:       "ユーザー1 ユーザー1",
			AlphabeticName: "user1 user1",
		},
		{
			TraQID:         MockUsers[1].Name,
			RealName:       "ユーザー2 ユーザー2",
			AlphabeticName: "user2 user2",
		},
		{
			TraQID:         MockUsers[2].Name,
			RealName:       "東 工子",
			AlphabeticName: "Noriko Azuma",
		},
	}
	MockTraQUsers = []*TraQUser{
		{
			User: &external.TraQUserResponse{
				ID:    MockUsers[0].ID,
				State: domain.TraqStateActive,
			},
			Name: MockUsers[0].Name,
		},
		{
			User: &external.TraQUserResponse{
				ID:    MockUsers[1].ID,
				State: domain.TraqStateDeactivated,
			},
			Name: MockUsers[1].Name,
		},
		{
			User: &external.TraQUserResponse{
				ID:    MockUsers[2].ID,
				State: domain.TraqStateActive,
			},
			Name: MockUsers[2].Name,
		},
	}
)
