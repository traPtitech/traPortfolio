package mockdata

import "github.com/traPtitech/traPortfolio/interfaces/handler"

var (
	HMockUser1 = handler.User{
		Id:       MockUsers[0].ID,
		Name:     MockUsers[0].Name,
		RealName: MockPortalUsers[0].RealName,
	}
	HMockUser2 = handler.User{
		Id:       MockUsers[1].ID,
		Name:     MockUsers[1].Name,
		RealName: MockPortalUsers[1].RealName,
	}
	HMockUser3 = handler.User{
		Id:       MockUsers[2].ID,
		Name:     MockUsers[2].Name,
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
