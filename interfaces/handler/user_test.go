package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
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

func TestUserHandler_GetAll(t *testing.T) {
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
					ruser := domain.User{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(),
						RealName: random.AlphaNumeric(),
					}
					huser := User{
						Id:       ruser.ID,
						Name:     ruser.Name,
						RealName: ruser.RealName,
					}

					repoUsers = append(repoUsers, &ruser)
					hresUsers = append(hresUsers, &huser)

				}

				args := repository.GetUsersArgs{}

				s.EXPECT().GetUsers(gomock.Any(), &args).Return(repoUsers, nil)
				return hresUsers, "/api/v1/users"
			},
			statusCode: http.StatusOK,
		},
		// TODO: オプションありのテストを追加する
		// TODO: Validationのテストを追加する
		{
			name: "internal error",
			setup: func(s *mock_service.MockUserService) (hres []*User, path string) {
				args := repository.GetUsersArgs{}

				s.EXPECT().GetUsers(gomock.Any(), &args).Return(nil, errors.New("Internal Server Error"))
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

func TestUserHandler_GetByID(t *testing.T) {
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
						Type:        uint(rand.Intn(int(domain.AccountLimit))),
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

					User: domain.User{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(),
						RealName: random.AlphaNumeric(),
					},
					State:    domain.TraQState(uint8(rand.Intn(int(domain.TraqStateLimit)))),
					Bio:      random.AlphaNumericn(rand.Intn(256) + 1),
					Accounts: rAccounts,
				}

				hresUser := UserDetail{
					User: User{
						Id:       repoUser.User.ID,
						Name:     repoUser.User.Name,
						RealName: repoUser.User.RealName,
					},
					Accounts: hAccounts,
					Bio:      repoUser.Bio,
					State:    UserAccountState(repoUser.State),
				}

				s.EXPECT().GetUser(gomock.Any(), repoUser.User.ID).Return(&repoUser, nil)
				path := fmt.Sprintf("/api/v1/users/%s", hresUser.User.Id)
				return &hresUser, path
			},
			statusCode: http.StatusOK,
		},

		{
			name: "internal error",
			setup: func(s *mock_service.MockUserService) (hres *UserDetail, userpath string) {
				id := random.UUID()
				s.EXPECT().GetUser(gomock.Any(), id).Return(nil, errors.New("Internal Server Error"))
				path := fmt.Sprintf("/api/v1/users/%s", id)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Bad Request: validate error: UUID",
			setup: func(s *mock_service.MockUserService) (hres *UserDetail, userpath string) {
				id := random.UUID()
				s.EXPECT().GetUser(gomock.Any(), id).Return(nil, repository.ErrValidate)
				path := fmt.Sprintf("/api/v1/users/%s", id)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(s *mock_service.MockUserService) (hres *UserDetail, userpath string) {
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

func TestUserHandler_Update(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockUserService) (reqBody *EditUser, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockUserService) (*EditUser, string) {

				userID := random.UUID()
				userBio := random.AlphaNumeric()
				userCheck := false
				if rand.Intn(2) == 1 {
					userCheck = true
				}

				reqBody := &EditUser{
					Bio:   &userBio,
					Check: &userCheck,
				}

				args := repository.UpdateUserArgs{
					Description: optional.StringFrom(&userBio),
					Check:       optional.BoolFrom(&userCheck),
				}

				path := fmt.Sprintf("/api/v1/users/%s", userID)
				s.EXPECT().Update(gomock.Any(), userID, &args).Return(nil)
				return reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "Conflict",
			setup: func(s *mock_service.MockUserService) (*EditUser, string) {

				userID := random.UUID()
				userBio := random.AlphaNumeric()
				userCheck := false
				if rand.Intn(2) == 1 {
					userCheck = true
				}

				reqBody := &EditUser{
					Bio:   &userBio,
					Check: &userCheck,
				}

				args := repository.UpdateUserArgs{
					Description: optional.StringFrom(&userBio),
					Check:       optional.BoolFrom(&userCheck),
				}

				path := fmt.Sprintf("/api/v1/users/%s", userID)
				s.EXPECT().Update(gomock.Any(), userID, &args).Return(repository.ErrAlreadyExists)
				return reqBody, path
			},
			statusCode: http.StatusConflict,
		},
		{
			name: "Not Found",
			setup: func(s *mock_service.MockUserService) (*EditUser, string) {

				userID := random.UUID()
				userBio := random.AlphaNumeric()
				userCheck := false
				if rand.Intn(2) == 1 {
					userCheck = true
				}

				reqBody := &EditUser{
					Bio:   &userBio,
					Check: &userCheck,
				}

				args := repository.UpdateUserArgs{
					Description: optional.StringFrom(&userBio),
					Check:       optional.BoolFrom(&userCheck),
				}

				path := fmt.Sprintf("/api/v1/users/%s", userID)
				s.EXPECT().Update(gomock.Any(), userID, &args).Return(repository.ErrNotFound)
				return reqBody, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error",
			setup: func(s *mock_service.MockUserService) (*EditUser, string) {
				path := fmt.Sprintf("/api/v1/users/%s", "invalid")
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
						Type:        uint(i),
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

				s.EXPECT().GetAccounts(gomock.Any()).Return(rAccounts, nil)
				path = fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return hAccounts, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(s *mock_service.MockUserService) (hres []*Account, path string) {

				userID := random.UUID()
				s.EXPECT().GetAccounts(gomock.Any()).Return(nil, errors.New("Internal Server Error"))
				path = fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Bad Request: validate error: UUID",
			setup: func(s *mock_service.MockUserService) (hres []*Account, path string) {

				userID := random.UUID()
				s.EXPECT().GetAccounts(gomock.Any()).Return(nil, repository.ErrValidate)
				path = fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(s *mock_service.MockUserService) (hres []*Account, path string) {

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
					Type:        uint(rand.Intn(int(domain.AccountLimit))),
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

				s.EXPECT().GetAccount(gomock.Any(), rAccount.ID).Return(&rAccount, nil)
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

				s.EXPECT().GetAccount(gomock.Any(), accountID).Return(nil, errors.New("Internal Server Error"))
				path = fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path

			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Bad Request: validate error: UUID",
			setup: func(s *mock_service.MockUserService) (hres *Account, path string) {

				userID := random.UUID()
				accountID := random.UUID()

				s.EXPECT().GetAccount(gomock.Any(), accountID).Return(nil, repository.ErrValidate)
				path = fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path

			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(s *mock_service.MockUserService) (hres *Account, path string) {

				userID := random.AlphaNumericn(36)
				accountID := random.UUID()

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
		setup      func(s *mock_service.MockUserService) (reqBody *AddUserAccountJSONBody, expectedResBody Account, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockUserService) (*AddUserAccountJSONBody, Account, string) {

				userID := random.UUID()

				reqBody := AddUserAccountJSONBody{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: PrPermitted(random.Bool()),
					Type:        AccountType((rand.Intn(int(domain.AccountLimit)))),
					Url:         random.RandURLString(),
				}

				args := repository.CreateAccountArgs{
					DisplayName: reqBody.DisplayName,
					Type:        uint(reqBody.Type),
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
				s.EXPECT().CreateAccount(gomock.Any(), userID, &args).Return(&want, nil)
				return &reqBody, expectedResBody, path
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "Bad Request: DisplayName is empty",
			setup: func(s *mock_service.MockUserService) (*AddUserAccountJSONBody, Account, string) {

				userID := random.UUID()

				reqBody := AddUserAccountJSONBody{
					DisplayName: "",
					PrPermitted: PrPermitted(random.Bool()),
					Type:        AccountType((rand.Intn(int(domain.AccountLimit)))),
					Url:         random.RandURLString(),
				}

				args := repository.CreateAccountArgs{
					DisplayName: reqBody.DisplayName,
					Type:        uint(reqBody.Type),
					URL:         reqBody.Url,
					PrPermitted: bool(reqBody.PrPermitted),
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				s.EXPECT().CreateAccount(gomock.Any(), userID, &args).Return(nil, repository.ErrInvalidArg)
				return &reqBody, Account{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: Account Type is invalid",
			setup: func(s *mock_service.MockUserService) (*AddUserAccountJSONBody, Account, string) {

				userID := random.UUID()

				reqBody := AddUserAccountJSONBody{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: PrPermitted(random.Bool()),
					Type:        AccountType(domain.AccountLimit),
					Url:         random.RandURLString(),
				}

				args := repository.CreateAccountArgs{
					DisplayName: reqBody.DisplayName,
					Type:        uint(reqBody.Type),
					URL:         reqBody.Url,
					PrPermitted: bool(reqBody.PrPermitted),
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				s.EXPECT().CreateAccount(gomock.Any(), userID, &args).Return(nil, repository.ErrInvalidArg)
				return &reqBody, Account{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: UUID",
			setup: func(s *mock_service.MockUserService) (*AddUserAccountJSONBody, Account, string) {

				userID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, Account{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(s *mock_service.MockUserService) (*AddUserAccountJSONBody, Account, string) {

				userID := random.AlphaNumericn(36)

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, Account{}, path
			},
			statusCode: http.StatusBadRequest,
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
		setup      func(s *mock_service.MockUserService) (reqBody *EditUserAccountJSONBody, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockUserService) (*EditUserAccountJSONBody, string) {

				userID := random.UUID()
				accountID := random.UUID()
				accountType := int64(rand.Intn(int(domain.AccountLimit)))
				accountPermit := random.Bool()

				argsName := random.AlphaNumeric()
				argsPermit := PrPermitted(accountPermit)
				argsType := AccountType(accountType)
				argsURL := random.RandURLString()

				reqBody := EditUserAccountJSONBody{
					DisplayName: &argsName,
					PrPermitted: &argsPermit,
					Type:        &argsType,
					Url:         &argsURL,
				}

				args := repository.UpdateAccountArgs{
					DisplayName: optional.StringFrom(&argsName),
					Type:        optional.Int64From(&accountType),
					URL:         optional.StringFrom(&argsURL),
					PrPermitted: optional.BoolFrom(&accountPermit),
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				s.EXPECT().EditAccount(gomock.Any(), userID, accountID, &args).Return(nil)
				return &reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
		/*{
			name: "Forbidden",
			setup: func(s *mock_service.MockUserService) (*EditUserAccountJSONBody, string) {

				userID := random.UUID()
				accountID := random.UUID()
				accountType := int64(domain.AccountLimit)
				accountPermit := random.Bool()

				argsName := random.AlphaNumeric()
				argsPermit := PrPermitted(accountPermit)
				argsType := AccountType(accountType)
				argsURL := random.RandURLString()

				reqBody := EditUserAccountJSONBody{
					DisplayName: &argsName,
					PrPermitted: &argsPermit,
					Type:        &argsType,
					Url:         &argsURL,
				}

				args := repository.UpdateAccountArgs{
					DisplayName: optional.StringFrom(&argsName),
					Type:        optional.Int64From(&accountType),
					URL:         optional.StringFrom(&argsURL),
					PrPermitted: optional.BoolFrom(&accountPermit),
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				s.EXPECT().EditAccount(gomock.Any(), userID, accountID, &args).Return(repository.ErrForbidden)
				return &reqBody, path
			},
			statusCode: http.StatusForbidden,
		},*/
		{
			name: "Not Found",
			setup: func(s *mock_service.MockUserService) (*EditUserAccountJSONBody, string) {

				userID := random.UUID()
				accountID := random.UUID()
				accountType := int64(domain.AccountLimit)
				accountPermit := random.Bool()

				argsName := random.AlphaNumeric()
				argsPermit := PrPermitted(accountPermit)
				argsType := AccountType(accountType)
				argsURL := random.RandURLString()

				reqBody := EditUserAccountJSONBody{
					DisplayName: &argsName,
					PrPermitted: &argsPermit,
					Type:        &argsType,
					Url:         &argsURL,
				}

				args := repository.UpdateAccountArgs{
					DisplayName: optional.StringFrom(&argsName),
					Type:        optional.Int64From(&accountType),
					URL:         optional.StringFrom(&argsURL),
					PrPermitted: optional.BoolFrom(&accountPermit),
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				s.EXPECT().EditAccount(gomock.Any(), userID, accountID, &args).Return(repository.ErrNotFound)
				return &reqBody, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error: UUID",
			setup: func(s *mock_service.MockUserService) (*EditUserAccountJSONBody, string) {

				userID := random.UUID()
				accountID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: nonUUID1",
			setup: func(s *mock_service.MockUserService) (*EditUserAccountJSONBody, string) {

				userID := random.AlphaNumericn(36)
				accountID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: nonUUID2",
			setup: func(s *mock_service.MockUserService) (*EditUserAccountJSONBody, string) {

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
				s.EXPECT().DeleteAccount(gomock.Any(), userID, accountID).Return(nil)
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
				s.EXPECT().DeleteAccount(gomock.Any(), userID, accountID).Return(repository.ErrForbidden)
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
				s.EXPECT().DeleteAccount(gomock.Any(), userID, accountID).Return(repository.ErrNotFound)
				return path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error: nonUUID1",
			setup: func(s *mock_service.MockUserService) string {

				userID := random.AlphaNumericn(36)
				accountID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: nonUUID2",
			setup: func(s *mock_service.MockUserService) string {

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

/*
func TestUserHandler_DeleteAccount(t *testing.T) {
	type fields struct {
		srv service.UserService
	}
	type args struct {
		_c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &UserHandler{
				srv: tt.fields.srv,
			}
			if err := DeleteAccount(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserDeleteAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHandler_GetProjects(t *testing.T) {
	type fields struct {
		srv service.UserService
	}
	type args struct {
		_c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &UserHandler{
				srv: tt.fields.srv,
			}
			if err := GetProjects(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserGetProjects() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHandler_GetUsers(t *testing.T) {
	type fields struct {
		srv service.UserService
	}
	type args struct {
		_c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &UserHandler{
				srv: tt.fields.srv,
			}
			if err := GetUsers(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserGetUsers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHandler_GetGroupsByUserID(t *testing.T) {
	type fields struct {
		srv service.UserService
	}
	type args struct {
		_c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &UserHandler{
				srv: tt.fields.srv,
			}
			if err := GetGroupsByUserID(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserGetGroupsByUserID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHandler_GetEvents(t *testing.T) {
	type fields struct {
		srv service.UserService
	}
	type args struct {
		_c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &UserHandler{
				srv: tt.fields.srv,
			}
			if err := GetEvents(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserGetEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newUser(t *testing.T) {
	type args struct {
		id       uuid.UUID
		name     string
		realName string
	}
	tests := []struct {
		name string
		args args
		want User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUser(tt.args.id, tt.args.name, tt.args.realName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newUserDetail(t *testing.T) {
	type args struct {
		user     User
		accounts []Account
		bio      string
		state    domain.TraQState
	}
	tests := []struct {
		name string
		args args
		want UserDetail
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUserDetail(tt.args.user, tt.args.accounts, tt.args.bio, tt.args.state); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUserDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newAccount(t *testing.T) {
	type args struct {
		id          uuid.UUID
		name        string
		atype       uint
		url         string
		prPermitted bool
	}
	tests := []struct {
		name string
		args args
		want Account
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newAccount(tt.args.id, tt.args.name, tt.args.atype, tt.args.url, tt.args.prPermitted); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newUserProject(t *testing.T) {
	type args struct {
		id           uuid.UUID
		name         string
		duration     YearWithSemesterDuration
		userDuration YearWithSemesterDuration
	}
	tests := []struct {
		name string
		args args
		want UserProject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUserProject(tt.args.id, tt.args.name, tt.args.duration, tt.args.userDuration); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUserProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newUserTeamWithUserName(t *testing.T) {
	type args struct {
		UserTeam UserTeam
		UserName string
	}
	tests := []struct {
		name string
		args args
		want UserTeamWithUserName
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUserTeamWithUserName(tt.args.UserTeam, tt.args.UserName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUserTeamWithUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newGroup(t *testing.T) {
	type args struct {
		id   uuid.UUID
		name string
	}
	tests := []struct {
		name string
		args args
		want Group
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newGroup(tt.args.id, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newUserGroup(t *testing.T) {
	type args struct {
		group    Group
		Duration YearWithSemesterDuration
	}
	tests := []struct {
		name string
		args args
		want UserGroup
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUserGroup(tt.args.group, tt.args.Duration); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newUserGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
