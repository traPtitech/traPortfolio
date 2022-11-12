package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/service/mock_service"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
)

func setupUserMock(t *testing.T) (*mock_service.MockUserService, API) {
	t.Helper()

	ctrl := gomock.NewController(t)
	s := mock_service.NewMockUserService(ctrl)
	api := NewAPI(nil, NewUserHandler(s), nil, nil, nil, nil)

	return s, api
}

func TestUserHandler_GetUsers(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockUserService) (hres []*User, path string)
		statusCode int
	}{
		{
			name: "Success_NoOpts",
			setup: func(s *mock_service.MockUserService) (hres []*User, path string) {

				casenum := 2
				repoUsers := []*domain.User{}
				hresUsers := []*User{}

				for i := 0; i < casenum; i++ {
					ruser := domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool())
					huser := User{
						Id:       ruser.ID,
						Name:     ruser.Name,
						RealName: ruser.RealName(),
					}

					repoUsers = append(repoUsers, ruser)
					hresUsers = append(hresUsers, &huser)

				}

				args := repository.GetUsersArgs{}

				s.EXPECT().GetUsers(anyCtx{}, &args).Return(repoUsers, nil)
				return hresUsers, "/api/v1/users"
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Success_WithOpts_IncludeSuspended",
			setup: func(s *mock_service.MockUserService) (hres []*User, path string) {
				casenum := 2
				repoUsers := []*domain.User{}
				hresUsers := []*User{}

				for i := 0; i < casenum; i++ {
					ruser := domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool())
					huser := User{
						Id:       ruser.ID,
						Name:     ruser.Name,
						RealName: ruser.RealName(),
					}

					repoUsers = append(repoUsers, ruser)
					hresUsers = append(hresUsers, &huser)
				}

				includeSuspened := random.Bool()
				args := repository.GetUsersArgs{
					IncludeSuspended: optional.New(includeSuspened, true),
				}

				s.EXPECT().GetUsers(anyCtx{}, &args).Return(repoUsers, nil)
				return hresUsers, fmt.Sprintf("/api/v1/users?includeSuspended=%t", includeSuspened)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Success_WithOpts_Name",
			setup: func(s *mock_service.MockUserService) (hres []*User, path string) {
				repoUsers := []*domain.User{
					domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
				}
				hresUsers := []*User{
					{
						Id:       repoUsers[0].ID,
						Name:     repoUsers[0].Name,
						RealName: repoUsers[0].RealName(),
					},
				}

				args := repository.GetUsersArgs{
					Name: optional.From(repoUsers[0].Name),
				}

				s.EXPECT().GetUsers(anyCtx{}, &args).Return(repoUsers, nil)
				return hresUsers, fmt.Sprintf("/api/v1/users?name=%s", repoUsers[0].Name)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "invalid args: multiple options",
			setup: func(s *mock_service.MockUserService) (hres []*User, path string) {
				return nil, fmt.Sprintf("/api/v1/users?includeSuspended=%t&name=%s", random.Bool(), random.AlphaNumeric())
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "internal error",
			setup: func(s *mock_service.MockUserService) (hres []*User, path string) {
				args := repository.GetUsersArgs{}

				s.EXPECT().GetUsers(anyCtx{}, &args).Return(nil, errors.New("Internal Server Error"))
				return nil, "/api/v1/users"
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupUserMock(t)

			hresUsers, path := tt.setup(s)

			var resBody []*User
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockUserService) (hres *UserDetail, userpath string)
		statusCode int
	}{
		{
			name: "success random",
			setup: func(s *mock_service.MockUserService) (hres *UserDetail, userpath string) {

				const accountNum int = 9
				rAccounts := []*domain.Account{}
				hAccounts := []Account{}

				for i := 0; i < accountNum; i++ {
					prRandom := false
					if rand.Intn(2) == 1 {
						prRandom = true
					}

					raccount := domain.Account{
						ID:          random.UUID(),
						DisplayName: random.AlphaNumeric(),
						Type:        domain.AccountType(random.Uint8n(uint8(domain.AccountLimit))),
						PrPermitted: prRandom,
						URL:         random.AlphaNumeric(),
					}

					haccount := Account{
						Id:          raccount.ID,
						DisplayName: raccount.DisplayName,
						PrPermitted: PrPermitted(prRandom),
						Type:        AccountType(raccount.Type),
						Url:         raccount.URL,
					}

					rAccounts = append(rAccounts, &raccount)
					hAccounts = append(hAccounts, haccount)
				}

				repoUser := domain.UserDetail{
					User:     *domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
					State:    domain.TraQState(random.Uint8n(uint8(domain.TraqStateLimit))),
					Bio:      random.AlphaNumericn(rand.Intn(256) + 1),
					Accounts: rAccounts,
				}

				hresUser := UserDetail{
					Accounts: hAccounts,
					Bio:      repoUser.Bio,
					Id:       repoUser.User.ID,
					Name:     repoUser.User.Name,
					RealName: repoUser.User.RealName(),
					State:    UserAccountState(repoUser.State),
				}

				s.EXPECT().GetUser(anyCtx{}, repoUser.User.ID).Return(&repoUser, nil)
				path := fmt.Sprintf("/api/v1/users/%s", hresUser.Id)
				return &hresUser, path
			},
			statusCode: http.StatusOK,
		},

		{
			name: "internal error",
			setup: func(s *mock_service.MockUserService) (hres *UserDetail, userpath string) {
				id := random.UUID()
				s.EXPECT().GetUser(anyCtx{}, id).Return(nil, errors.New("Internal Server Error"))
				path := fmt.Sprintf("/api/v1/users/%s", id)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Bad Request: validate error: UUID",
			setup: func(s *mock_service.MockUserService) (hres *UserDetail, userpath string) {
				id := random.UUID()
				s.EXPECT().GetUser(anyCtx{}, id).Return(nil, repository.ErrValidate)
				path := fmt.Sprintf("/api/v1/users/%s", id)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(_ *mock_service.MockUserService) (hres *UserDetail, userpath string) {
				id := random.AlphaNumericn(36)
				path := fmt.Sprintf("/api/v1/users/%s", id)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupUserMock(t)

			var resBody *UserDetail

			hresUsers, userpath := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodGet, userpath, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)

		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockUserService) (reqBody *EditUserJSONRequestBody, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockUserService) (*EditUserJSONRequestBody, string) {

				userID := random.UUID()
				userBio := random.AlphaNumeric()
				userCheck := false
				if rand.Intn(2) == 1 {
					userCheck = true
				}

				reqBody := &EditUserJSONRequestBody{
					Bio:   &userBio,
					Check: &userCheck,
				}

				args := repository.UpdateUserArgs{
					Description: optional.FromPtr(&userBio),
					Check:       optional.FromPtr(&userCheck),
				}

				path := fmt.Sprintf("/api/v1/users/%s", userID)
				s.EXPECT().Update(anyCtx{}, userID, &args).Return(nil)
				return reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "Success with description args(len=256)",
			setup: func(s *mock_service.MockUserService) (*EditUserJSONRequestBody, string) {

				userID := random.UUID()
				userBio := strings.Repeat("a", 256)
				userCheck := false
				if rand.Intn(2) == 1 {
					userCheck = true
				}

				reqBody := &EditUserJSONRequestBody{
					Bio:   &userBio,
					Check: &userCheck,
				}

				args := repository.UpdateUserArgs{
					Description: optional.FromPtr(&userBio),
					Check:       optional.FromPtr(&userCheck),
				}

				path := fmt.Sprintf("/api/v1/users/%s", userID)
				s.EXPECT().Update(anyCtx{}, userID, &args).Return(nil)
				return reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "Conflict",
			setup: func(s *mock_service.MockUserService) (*EditUserJSONRequestBody, string) {

				userID := random.UUID()
				userBio := random.AlphaNumeric()
				userCheck := false
				if rand.Intn(2) == 1 {
					userCheck = true
				}

				reqBody := &EditUserJSONRequestBody{
					Bio:   &userBio,
					Check: &userCheck,
				}

				args := repository.UpdateUserArgs{
					Description: optional.FromPtr(&userBio),
					Check:       optional.FromPtr(&userCheck),
				}

				path := fmt.Sprintf("/api/v1/users/%s", userID)
				s.EXPECT().Update(anyCtx{}, userID, &args).Return(repository.ErrAlreadyExists)
				return reqBody, path
			},
			statusCode: http.StatusConflict,
		},
		{
			name: "Not Found",
			setup: func(s *mock_service.MockUserService) (*EditUserJSONRequestBody, string) {

				userID := random.UUID()
				userBio := random.AlphaNumeric()
				userCheck := false
				if rand.Intn(2) == 1 {
					userCheck = true
				}

				reqBody := &EditUserJSONRequestBody{
					Bio:   &userBio,
					Check: &userCheck,
				}

				args := repository.UpdateUserArgs{
					Description: optional.FromPtr(&userBio),
					Check:       optional.FromPtr(&userCheck),
				}

				path := fmt.Sprintf("/api/v1/users/%s", userID)
				s.EXPECT().Update(anyCtx{}, userID, &args).Return(repository.ErrNotFound)
				return reqBody, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: invalid userID",
			setup: func(_ *mock_service.MockUserService) (*EditUserJSONRequestBody, string) {
				path := fmt.Sprintf("/api/v1/users/%s", "invalid")
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: too long description(len>256)",
			setup: func(_ *mock_service.MockUserService) (*EditUserJSONRequestBody, string) {
				userID := random.UUID()
				userBio := strings.Repeat("a", 257)

				reqBody := &EditUserJSONRequestBody{
					Bio: &userBio,
				}

				path := fmt.Sprintf("/api/v1/users/%s", userID)
				return reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupUserMock(t)

			reqBody, path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodPatch, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}

func TestUserHandler_GetUserAccounts(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockUserService) (hres []*Account, path string)
		statusCode int
	}{
		{
			name: "success",
			setup: func(s *mock_service.MockUserService) (hres []*Account, path string) {

				userID := random.UUID()
				accountKinds := rand.Intn((1<<domain.AccountLimit)-1) + 1
				//AccountLimit種類のうち、テストに使うものだけbitが立っている
				//例えば0(HOMEPAGE)と2(TWITTER)と7(ATCODER)なら10000101=133
				//0(bitがすべて立っていない)は除外

				rAccounts := []*domain.Account{}
				hAccounts := []*Account{}

				for i := 0; i < int(domain.AccountLimit); i++ {
					if (accountKinds>>i)%2 == 0 {
						continue
					}

					prRandom := false
					if rand.Intn(2) == 1 {
						prRandom = true
					}

					raccount := domain.Account{
						ID:          random.UUID(),
						DisplayName: random.AlphaNumeric(),
						Type:        domain.AccountType(uint8(i)),
						PrPermitted: prRandom,
						URL:         random.AlphaNumeric(),
					}

					haccount := Account{
						Id:          raccount.ID,
						DisplayName: raccount.DisplayName,
						PrPermitted: PrPermitted(prRandom),
						Type:        AccountType(raccount.Type),
						Url:         raccount.URL,
					}

					rAccounts = append(rAccounts, &raccount)
					hAccounts = append(hAccounts, &haccount)

				}

				s.EXPECT().GetAccounts(userID).Return(rAccounts, nil)
				path = fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return hAccounts, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(s *mock_service.MockUserService) (hres []*Account, path string) {

				userID := random.UUID()
				s.EXPECT().GetAccounts(userID).Return(nil, errors.New("Internal Server Error"))
				path = fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Bad Request: validate error: UUID",
			setup: func(s *mock_service.MockUserService) (hres []*Account, path string) {

				userID := random.UUID()
				s.EXPECT().GetAccounts(userID).Return(nil, repository.ErrValidate)
				path = fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(_ *mock_service.MockUserService) (hres []*Account, path string) {

				userID := random.AlphaNumericn(36)
				path = fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupUserMock(t)

			hresUsers, path := tt.setup(s)

			var resBody []*Account
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_GetUserAccount(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockUserService) (hres *Account, path string)
		statusCode int
	}{
		{
			name: "success",
			setup: func(s *mock_service.MockUserService) (hres *Account, path string) {

				userID := random.UUID()
				prRandom := false
				if rand.Intn(2) == 1 {
					prRandom = true
				}

				rAccount := domain.Account{
					ID:          random.UUID(),
					DisplayName: random.AlphaNumeric(),
					Type:        domain.AccountType(uint8(random.Uint8n(uint8(domain.AccountLimit)))),
					PrPermitted: prRandom,
					URL:         random.AlphaNumeric(),
				}
				hAccount := Account{
					Id:          rAccount.ID,
					DisplayName: rAccount.DisplayName,
					PrPermitted: PrPermitted(prRandom),
					Type:        AccountType(rAccount.Type),
					Url:         rAccount.URL,
				}

				s.EXPECT().GetAccount(userID, rAccount.ID).Return(&rAccount, nil)
				path = fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, rAccount.ID)
				return &hAccount, path

			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(s *mock_service.MockUserService) (hres *Account, path string) {

				userID := random.UUID()
				accountID := random.UUID()

				s.EXPECT().GetAccount(userID, accountID).Return(nil, errors.New("Internal Server Error"))
				path = fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path

			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Bad Request: validate error: invalid userID",
			setup: func(_ *mock_service.MockUserService) (hres *Account, path string) {

				userID := random.AlphaNumericn(36)
				accountID := random.UUID()

				path = fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path

			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(_ *mock_service.MockUserService) (hres *Account, path string) {

				userID := random.UUID()
				accountID := random.AlphaNumericn(36)

				path = fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path

			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupUserMock(t)

			hresUsers, path := tt.setup(s)

			var resBody *Account
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_AddUserAccount(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockUserService) (reqBody *AddUserAccountJSONRequestBody, expectedResBody Account, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockUserService) (*AddUserAccountJSONRequestBody, Account, string) {

				userID := random.UUID()
				accountType := domain.AccountType(rand.Intn(int(domain.AccountLimit)))

				reqBody := AddUserAccountJSONRequestBody{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: PrPermitted(random.Bool()),
					Type:        AccountType(accountType),
					Url:         random.AccountURLString(accountType),
				}

				args := repository.CreateAccountArgs{
					DisplayName: reqBody.DisplayName,
					Type:        domain.AccountType(uint8(reqBody.Type)),
					URL:         reqBody.Url,
					PrPermitted: bool(reqBody.PrPermitted),
				}

				want := domain.Account{
					ID:          userID,
					DisplayName: args.DisplayName,
					Type:        args.Type,
					PrPermitted: args.PrPermitted,
					URL:         args.URL,
				}

				expectedResBody := Account{
					Id:          userID,
					DisplayName: reqBody.DisplayName,
					PrPermitted: reqBody.PrPermitted,
					Type:        reqBody.Type,
					Url:         reqBody.Url,
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				s.EXPECT().CreateAccount(anyCtx{}, userID, &args).Return(&want, nil)
				return &reqBody, expectedResBody, path
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "Success: Account Type is 0",
			setup: func(s *mock_service.MockUserService) (*AddUserAccountJSONRequestBody, Account, string) {
				userID := random.UUID()

				reqBody := AddUserAccountJSONRequestBody{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: PrPermitted(random.Bool()),
					Type:        0,
					Url:         random.AccountURLString(0),
				}

				args := repository.CreateAccountArgs{
					DisplayName: reqBody.DisplayName,
					Type:        domain.AccountType(uint8(reqBody.Type)),
					URL:         reqBody.Url,
					PrPermitted: bool(reqBody.PrPermitted),
				}

				want := domain.Account{
					ID:          userID,
					DisplayName: args.DisplayName,
					Type:        args.Type,
					PrPermitted: args.PrPermitted,
					URL:         args.URL,
				}

				expectedResBody := Account{
					Id:          userID,
					DisplayName: reqBody.DisplayName,
					PrPermitted: reqBody.PrPermitted,
					Type:        reqBody.Type,
					Url:         reqBody.Url,
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				s.EXPECT().CreateAccount(anyCtx{}, userID, &args).Return(&want, nil)
				return &reqBody, expectedResBody, path
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "Bad Request: DisplayName is empty",
			setup: func(_ *mock_service.MockUserService) (*AddUserAccountJSONRequestBody, Account, string) {

				userID := random.UUID()
				accountType := domain.AccountType(rand.Intn(int(domain.AccountLimit)))

				reqBody := AddUserAccountJSONRequestBody{
					DisplayName: "",
					PrPermitted: PrPermitted(random.Bool()),
					Type:        AccountType(accountType),
					Url:         random.AccountURLString(accountType),
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return &reqBody, Account{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: Account Type is invalid",
			setup: func(_ *mock_service.MockUserService) (*AddUserAccountJSONRequestBody, Account, string) {

				userID := random.UUID()

				reqBody := AddUserAccountJSONRequestBody{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: PrPermitted(random.Bool()),
					Type:        AccountType(domain.AccountLimit),
					Url:         random.RandURLString(),
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return &reqBody, Account{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: UUID",
			setup: func(_ *mock_service.MockUserService) (*AddUserAccountJSONRequestBody, Account, string) {

				userID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, Account{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(_ *mock_service.MockUserService) (*AddUserAccountJSONRequestBody, Account, string) {

				userID := random.AlphaNumericn(36)

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, Account{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "internal error",
			setup: func(s *mock_service.MockUserService) (*AddUserAccountJSONRequestBody, Account, string) {
				userID := random.UUID()
				accountType := domain.AccountType(rand.Intn(int(domain.AccountLimit)))

				reqBody := AddUserAccountJSONRequestBody{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: PrPermitted(random.Bool()),
					Type:        AccountType(accountType),
					Url:         random.AccountURLString(accountType),
				}

				args := repository.CreateAccountArgs{
					DisplayName: reqBody.DisplayName,
					Type:        domain.AccountType(uint8(reqBody.Type)),
					URL:         reqBody.Url,
					PrPermitted: bool(reqBody.PrPermitted),
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				s.EXPECT().CreateAccount(anyCtx{}, userID, &args).Return(nil, errors.New("internal error"))
				return &reqBody, Account{}, path
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupUserMock(t)

			reqBody, res, path := tt.setup(s)

			var resBody Account
			statusCode, _ := doRequest(t, api, http.MethodPost, path, reqBody, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, res, resBody)
		})
	}
}

func TestUserHandler_EditUserAccount(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockUserService) (reqBody *EditUserAccountJSONRequestBody, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockUserService) (*EditUserAccountJSONRequestBody, string) {

				userID := random.UUID()
				accountID := random.UUID()
				accountType := domain.AccountType(rand.Intn(int(domain.AccountLimit))) // TODO: domain.AccountType型にする
				accountPermit := random.Bool()

				argsName := random.AlphaNumeric()
				argsPermit := PrPermitted(accountPermit)
				argsType := AccountType(accountType)
				argsURL := random.AccountURLString(domain.AccountType(accountType))

				reqBody := EditUserAccountJSONRequestBody{
					DisplayName: &argsName,
					PrPermitted: &argsPermit,
					Type:        &argsType,
					Url:         &argsURL,
				}

				args := repository.UpdateAccountArgs{
					DisplayName: optional.FromPtr(&argsName),
					Type:        optional.FromPtr(&accountType),
					URL:         optional.FromPtr(&argsURL),
					PrPermitted: optional.FromPtr(&accountPermit),
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				s.EXPECT().EditAccount(anyCtx{}, userID, accountID, &args).Return(nil)
				return &reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "Not Found",
			setup: func(s *mock_service.MockUserService) (*EditUserAccountJSONRequestBody, string) {

				userID := random.UUID()
				accountID := random.UUID()
				accountType := domain.AccountType(rand.Intn(int(domain.AccountLimit)))
				accountPermit := random.Bool()

				argsName := random.AlphaNumeric()
				argsPermit := PrPermitted(accountPermit)
				argsType := AccountType(accountType)
				argsURL := random.AccountURLString(domain.AccountType(accountType))

				reqBody := EditUserAccountJSONRequestBody{
					DisplayName: &argsName,
					PrPermitted: &argsPermit,
					Type:        &argsType,
					Url:         &argsURL,
				}

				args := repository.UpdateAccountArgs{
					DisplayName: optional.FromPtr(&argsName),
					Type:        optional.FromPtr(&accountType),
					URL:         optional.FromPtr(&argsURL),
					PrPermitted: optional.FromPtr(&accountPermit),
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				s.EXPECT().EditAccount(anyCtx{}, userID, accountID, &args).Return(repository.ErrNotFound)
				return &reqBody, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error: empty display name(but not nil)",
			setup: func(_ *mock_service.MockUserService) (*EditUserAccountJSONRequestBody, string) {
				userID := random.UUID()
				accountID := random.UUID()

				argsName := "" // empty but not nil

				reqBody := EditUserAccountJSONRequestBody{
					DisplayName: &argsName,
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)

				return &reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: too large account type",
			setup: func(_ *mock_service.MockUserService) (*EditUserAccountJSONRequestBody, string) {
				userID := random.UUID()
				accountID := random.UUID()

				argsType := AccountType(domain.AccountLimit)

				reqBody := EditUserAccountJSONRequestBody{
					Type: &argsType,
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)

				return &reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: invalid url",
			setup: func(_ *mock_service.MockUserService) (*EditUserAccountJSONRequestBody, string) {
				userID := random.UUID()
				accountID := random.UUID()

				argsURL := random.AlphaNumeric()

				reqBody := EditUserAccountJSONRequestBody{
					Url: &argsURL,
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)

				return &reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: nonUUID1",
			setup: func(_ *mock_service.MockUserService) (*EditUserAccountJSONRequestBody, string) {

				userID := random.AlphaNumericn(36)
				accountID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: nonUUID2",
			setup: func(_ *mock_service.MockUserService) (*EditUserAccountJSONRequestBody, string) {

				userID := random.UUID()
				accountID := random.AlphaNumericn(36)

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupUserMock(t)

			reqBody, path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodPatch, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}

func TestUserHandler_DeleteUserAccount(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockUserService) (path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockUserService) string {

				userID := random.UUID()
				accountID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				s.EXPECT().DeleteAccount(anyCtx{}, userID, accountID).Return(nil)
				return path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "Forbidden",
			setup: func(s *mock_service.MockUserService) string {

				userID := random.UUID()
				accountID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				s.EXPECT().DeleteAccount(anyCtx{}, userID, accountID).Return(repository.ErrForbidden)
				return path
			},
			statusCode: http.StatusForbidden,
		},
		{
			name: "Not Found",
			setup: func(s *mock_service.MockUserService) string {

				userID := random.UUID()
				accountID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				s.EXPECT().DeleteAccount(anyCtx{}, userID, accountID).Return(repository.ErrNotFound)
				return path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error: nonUUID1",
			setup: func(_ *mock_service.MockUserService) string {

				userID := random.AlphaNumericn(36)
				accountID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: nonUUID2",
			setup: func(_ *mock_service.MockUserService) string {

				userID := random.UUID()
				accountID := random.AlphaNumericn(36)

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupUserMock(t)

			path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodDelete, path, nil, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}

func TestUserHandler_GetUserProjects(t *testing.T) {

	makeProjects := func(s *mock_service.MockUserService, projectsLen int) (hres []*UserProject, path string) {
		userID := random.UUID()

		repoProjects := []*domain.UserProject{}
		hresProjects := []*UserProject{}

		for i := 0; i < projectsLen; i++ {

			//TODO: DurationはUserDurationを包含しているべき
			rproject := domain.UserProject{
				ID:           random.UUID(),
				Name:         random.AlphaNumeric(),
				Duration:     random.Duration(),
				UserDuration: random.Duration(),
			}

			hproject := UserProject{
				Duration:     ConvertDuration(rproject.Duration),
				Id:           rproject.ID,
				Name:         rproject.Name,
				UserDuration: ConvertDuration(rproject.UserDuration),
			}

			repoProjects = append(repoProjects, &rproject)
			hresProjects = append(hresProjects, &hproject)

		}

		s.EXPECT().GetUserProjects(anyCtx{}, userID).Return(repoProjects, nil)
		path = fmt.Sprintf("/api/v1/users/%s/projects", userID)
		return hresProjects, path
	}

	tests := []struct {
		name       string
		setup      func(s *mock_service.MockUserService) (hres []*UserProject, path string)
		statusCode int
	}{
		{
			name: "success 1",
			setup: func(s *mock_service.MockUserService) (hres []*UserProject, path string) {
				return makeProjects(s, 1)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 2",
			setup: func(s *mock_service.MockUserService) (hres []*UserProject, path string) {
				return makeProjects(s, 2)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 32",
			setup: func(s *mock_service.MockUserService) (hres []*UserProject, path string) {
				return makeProjects(s, 32)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Not Found",
			setup: func(s *mock_service.MockUserService) (hres []*UserProject, path string) {

				userID := random.UUID()

				s.EXPECT().GetUserProjects(anyCtx{}, userID).Return(nil, repository.ErrNotFound)
				path = fmt.Sprintf("/api/v1/users/%s/projects", userID)
				return nil, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error",
			setup: func(_ *mock_service.MockUserService) (hres []*UserProject, path string) {

				userID := random.AlphaNumericn(36)

				path = fmt.Sprintf("/api/v1/users/%s/projects", userID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupUserMock(t)

			hresUsers, path := tt.setup(s)
			var resBody []*UserProject
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)

		})
	}
}

func TestUserHandler_GetUserContests(t *testing.T) {

	makeContests := func(s *mock_service.MockUserService, contestsLen int) (hres []*UserContest, path string) {
		userID := random.UUID()

		repoContests := []*domain.UserContest{}
		hresContests := []*UserContest{}

		for i := 0; i < contestsLen; i++ {

			rcontest := domain.UserContest{
				ID:        random.UUID(),
				Name:      random.AlphaNumeric(),
				TimeStart: random.Time(),
				TimeEnd:   random.Time(),
				Teams: []*domain.ContestTeam{
					{
						ID:        random.UUID(),
						ContestID: random.UUID(),
						Name:      random.AlphaNumeric(),
						Result:    random.AlphaNumeric(),
					},
				},
			}

			hcontest := newUserContest(
				newContest(rcontest.ID, rcontest.Name, rcontest.TimeStart, rcontest.TimeEnd),
				[]ContestTeam{
					newContestTeam(rcontest.Teams[0].ID, rcontest.Teams[0].Name, rcontest.Teams[0].Result),
				},
			)

			repoContests = append(repoContests, &rcontest)
			hresContests = append(hresContests, &hcontest)

		}

		s.EXPECT().GetUserContests(anyCtx{}, userID).Return(repoContests, nil)
		path = fmt.Sprintf("/api/v1/users/%s/contests", userID)
		return hresContests, path

	}

	tests := []struct {
		name       string
		setup      func(s *mock_service.MockUserService) (hres []*UserContest, path string)
		statusCode int
	}{
		{
			name: "success 1",
			setup: func(s *mock_service.MockUserService) (hres []*UserContest, path string) {
				return makeContests(s, 1)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 2",
			setup: func(s *mock_service.MockUserService) (hres []*UserContest, path string) {
				return makeContests(s, 2)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 32",
			setup: func(s *mock_service.MockUserService) (hres []*UserContest, path string) {
				return makeContests(s, 32)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Not Found",
			setup: func(s *mock_service.MockUserService) (hres []*UserContest, path string) {

				userID := random.UUID()

				s.EXPECT().GetUserContests(anyCtx{}, userID).Return(nil, repository.ErrNotFound)
				path = fmt.Sprintf("/api/v1/users/%s/contests", userID)
				return nil, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error",
			setup: func(_ *mock_service.MockUserService) (hres []*UserContest, path string) {

				userID := random.AlphaNumericn(36)

				path = fmt.Sprintf("/api/v1/users/%s/contests", userID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupUserMock(t)

			hresUsers, path := tt.setup(s)
			var resBody []*UserContest
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_GetUserGroups(t *testing.T) {

	makeGroups := func(s *mock_service.MockUserService, groupsLen int) (hres []*UserGroup, path string) {
		userID := random.UUID()

		repoGroups := []*domain.UserGroup{}
		hresGroups := []*UserGroup{}

		for i := 0; i < groupsLen; i++ {

			rgroup := domain.UserGroup{
				ID:       random.UUID(),
				Name:     random.AlphaNumeric(),
				Duration: random.Duration(),
			}

			hgroup := UserGroup{
				Duration: ConvertDuration(rgroup.Duration),
				Id:       rgroup.ID,
				Name:     rgroup.Name,
			}

			repoGroups = append(repoGroups, &rgroup)
			hresGroups = append(hresGroups, &hgroup)

		}

		s.EXPECT().GetGroupsByUserID(anyCtx{}, userID).Return(repoGroups, nil)
		path = fmt.Sprintf("/api/v1/users/%s/groups", userID)
		return hresGroups, path

	}

	tests := []struct {
		name       string
		setup      func(s *mock_service.MockUserService) (hres []*UserGroup, path string)
		statusCode int
	}{
		{
			name: "success 0",
			setup: func(s *mock_service.MockUserService) (hres []*UserGroup, path string) {
				return makeGroups(s, 0)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 1",
			setup: func(s *mock_service.MockUserService) (hres []*UserGroup, path string) {
				return makeGroups(s, 1)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 2",
			setup: func(s *mock_service.MockUserService) (hres []*UserGroup, path string) {
				return makeGroups(s, 2)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 32",
			setup: func(s *mock_service.MockUserService) (hres []*UserGroup, path string) {
				return makeGroups(s, 32)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Not Found",
			setup: func(s *mock_service.MockUserService) (hres []*UserGroup, path string) {

				userID := random.UUID()

				s.EXPECT().GetGroupsByUserID(anyCtx{}, userID).Return(nil, repository.ErrNotFound)
				path = fmt.Sprintf("/api/v1/users/%s/groups", userID)
				return nil, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error",
			setup: func(_ *mock_service.MockUserService) (hres []*UserGroup, path string) {

				userID := random.AlphaNumericn(36)

				path = fmt.Sprintf("/api/v1/users/%s/groups", userID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupUserMock(t)

			hresUsers, path := tt.setup(s)
			var resBody []*UserGroup
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_GetUserEvents(t *testing.T) {

	makeEvents := func(s *mock_service.MockUserService, eventsLen int) (hres []*Event, path string) {
		userID := random.UUID()

		repoEvents := []*domain.Event{}
		hresEvents := []*Event{}

		for i := 0; i < eventsLen; i++ {

			timeStart, timeEnd := random.SinceAndUntil()

			revent := domain.Event{
				ID:        random.UUID(),
				Name:      random.AlphaNumeric(),
				TimeStart: timeStart,
				TimeEnd:   timeEnd,
			}

			hevent := Event{
				Duration: Duration{
					Since: timeStart,
					Until: &timeEnd,
				},
				Id:   revent.ID,
				Name: revent.Name,
			}

			repoEvents = append(repoEvents, &revent)
			hresEvents = append(hresEvents, &hevent)

		}

		s.EXPECT().GetUserEvents(anyCtx{}, userID).Return(repoEvents, nil)
		path = fmt.Sprintf("/api/v1/users/%s/events", userID)
		return hresEvents, path

	}

	tests := []struct {
		name       string
		setup      func(s *mock_service.MockUserService) (hres []*Event, path string)
		statusCode int
	}{
		{
			name: "success 0",
			setup: func(s *mock_service.MockUserService) (hres []*Event, path string) {
				return makeEvents(s, 0)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 1",
			setup: func(s *mock_service.MockUserService) (hres []*Event, path string) {
				return makeEvents(s, 1)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 2",
			setup: func(s *mock_service.MockUserService) (hres []*Event, path string) {
				return makeEvents(s, 2)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 32",
			setup: func(s *mock_service.MockUserService) (hres []*Event, path string) {
				return makeEvents(s, 32)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Not Found",
			setup: func(s *mock_service.MockUserService) (hres []*Event, path string) {

				userID := random.UUID()

				s.EXPECT().GetUserEvents(anyCtx{}, userID).Return(nil, repository.ErrNotFound)
				path = fmt.Sprintf("/api/v1/users/%s/events", userID)
				return nil, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error",
			setup: func(_ *mock_service.MockUserService) (hres []*Event, path string) {

				userID := random.AlphaNumericn(36)

				path = fmt.Sprintf("/api/v1/users/%s/events", userID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupUserMock(t)

			hresUsers, path := tt.setup(s)
			var resBody []*Event
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}
