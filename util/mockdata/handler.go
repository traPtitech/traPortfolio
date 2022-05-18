package mockdata

import "github.com/traPtitech/traPortfolio/interfaces/handler"

var (
	HMockUser1 = handler.User{
		Id:       userID1,
		Name:     userName1,
		RealName: MockPortalUsers[0].RealName,
	}
	HMockUser2 = handler.User{
		Id:       userID2,
		Name:     userName2,
		RealName: MockPortalUsers[1].RealName,
	}
	HMockUser3 = handler.User{
		Id:       userID3,
		Name:     userName2,
		RealName: MockPortalUsers[2].RealName,
	}

	HMockAccount = handler.Account{
		DisplayName: MockAccount.Name,
		Id:          MockAccount.ID,
		PrPermitted: handler.PrPermitted(MockAccount.Check),
		Type:        handler.AccountType(MockAccount.Type),
		Url:         MockAccount.URL,
	}

	HMockUserDetail1 = handler.UserDetail{
		User:     HMockUser1,
		Accounts: []handler.Account{HMockAccount},
		Bio:      MockUsers[0].Description,
		State:    handler.UserAccountState(MockTraQUsers[0].User.State),
	}
)
