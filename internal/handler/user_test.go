package handler

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/handler/schema"
	"github.com/traPtitech/traPortfolio/internal/repository"
	"github.com/traPtitech/traPortfolio/internal/repository/mock_repository"
	"github.com/traPtitech/traPortfolio/internal/util/optional"
	"github.com/traPtitech/traPortfolio/internal/util/random"
)

func setupUserMock(t *testing.T) (MockRepository, API) {
	t.Helper()

	ctrl := gomock.NewController(t)
	user := mock_repository.NewMockUserRepository(ctrl)
	event := mock_repository.NewMockEventRepository(ctrl)
	mr := MockRepository{user: user, event: event}
	api := NewAPI(nil, NewUserHandler(user, event), nil, nil, nil, nil)

	return mr, api
}

func TestUserHandler_GetUsers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (hres []*schema.User, path string)
		statusCode int
	}{
		{
			name: "Success_NoOpts",
			setup: func(mr MockRepository) (hres []*schema.User, path string) {
				casenum := 2
				repoUsers := []*domain.User{}
				hresUsers := []*schema.User{}

				for range casenum {
					ruser := domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool())
					huser := schema.User{
						Id:       ruser.ID,
						Name:     ruser.Name,
						RealName: ruser.RealName(),
					}

					repoUsers = append(repoUsers, ruser)
					hresUsers = append(hresUsers, &huser)
				}

				args := repository.GetUsersArgs{}

				mr.user.EXPECT().GetUsers(anyCtx{}, &args).Return(repoUsers, nil)
				return hresUsers, "/api/v1/users"
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Success_WithOpts_IncludeSuspended",
			setup: func(mr MockRepository) (hres []*schema.User, path string) {
				casenum := 2
				repoUsers := []*domain.User{}
				hresUsers := []*schema.User{}

				for range casenum {
					ruser := domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool())
					huser := schema.User{
						Id:       ruser.ID,
						Name:     ruser.Name,
						RealName: ruser.RealName(),
					}

					repoUsers = append(repoUsers, ruser)
					hresUsers = append(hresUsers, &huser)
				}

				includeSuspened := random.Bool()
				args := repository.GetUsersArgs{
					IncludeSuspended: optional.From(includeSuspened),
				}

				mr.user.EXPECT().GetUsers(anyCtx{}, &args).Return(repoUsers, nil)
				return hresUsers, fmt.Sprintf("/api/v1/users?includeSuspended=%t", includeSuspened)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Success_WithOpts_Name",
			setup: func(mr MockRepository) (hres []*schema.User, path string) {
				repoUsers := []*domain.User{
					domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
				}
				hresUsers := []*schema.User{
					{
						Id:       repoUsers[0].ID,
						Name:     repoUsers[0].Name,
						RealName: repoUsers[0].RealName(),
					},
				}

				args := repository.GetUsersArgs{
					Name: optional.From(repoUsers[0].Name),
				}

				mr.user.EXPECT().GetUsers(anyCtx{}, &args).Return(repoUsers, nil)
				return hresUsers, fmt.Sprintf("/api/v1/users?name=%s", repoUsers[0].Name)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "invalid args: multiple options",
			setup: func(_ MockRepository) (hres []*schema.User, path string) {
				return nil, fmt.Sprintf("/api/v1/users?includeSuspended=%t&name=%s", random.Bool(), random.AlphaNumeric())
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "internal error",
			setup: func(mr MockRepository) (hres []*schema.User, path string) {
				args := repository.GetUsersArgs{}

				mr.user.EXPECT().GetUsers(anyCtx{}, &args).Return(nil, errors.New("Internal Server Error"))
				return nil, "/api/v1/users"
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "invalid limit with 0",
			setup: func(_ MockRepository) (hres []*schema.User, path string) {
				return nil, fmt.Sprintf("/api/v1/users?limit=%d", 0)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "invalid limit less than 1",
			setup: func(_ MockRepository) (hres []*schema.User, path string) {
				return nil, fmt.Sprintf("/api/v1/users?limit=%d", -1)
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupUserMock(t)

			hresUsers, path := tt.setup(mr)

			var resBody []*schema.User
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_SyncUsers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(mr MockRepository) (path string) {
				mr.user.EXPECT().SyncUsers(anyCtx{}).Return(nil)
				return "/api/v1/users/sync"
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "internal error",
			setup: func(mr MockRepository) (path string) {
				mr.user.EXPECT().SyncUsers(anyCtx{}).Return(errors.New("Internal Server Error"))
				return "/api/v1/users/sync"
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			s, api := setupUserMock(t)

			path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodPost, path, nil, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (hres *schema.UserDetail, userpath string)
		statusCode int
	}{
		{
			name: "success random",
			setup: func(mr MockRepository) (hres *schema.UserDetail, userpath string) {
				const accountNum int = 9
				rAccounts := []*domain.Account{}
				hAccounts := []schema.Account{}

				for range accountNum {
					prRandom := random.Bool()

					raccount := domain.Account{
						ID:          random.UUID(),
						DisplayName: random.AlphaNumeric(),
						Type:        rand.N(domain.AccountLimit),
						PrPermitted: prRandom,
						URL:         random.AlphaNumeric(),
					}

					haccount := schema.Account{
						Id:          raccount.ID,
						DisplayName: raccount.DisplayName,
						PrPermitted: schema.PrPermitted(prRandom),
						Type:        schema.AccountType(raccount.Type),
						Url:         raccount.URL,
					}

					rAccounts = append(rAccounts, &raccount)
					hAccounts = append(hAccounts, haccount)
				}

				repoUser := domain.UserDetail{
					User:     *domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
					State:    rand.N(domain.TraqStateLimit),
					Bio:      random.AlphaNumericN(rand.IntN(256) + 1),
					Accounts: rAccounts,
				}

				hresUser := schema.UserDetail{
					Accounts: hAccounts,
					Bio:      repoUser.Bio,
					Id:       repoUser.User.ID,
					Name:     repoUser.User.Name,
					RealName: repoUser.User.RealName(),
					State:    schema.UserAccountState(repoUser.State),
				}

				mr.user.EXPECT().GetUser(anyCtx{}, repoUser.User.ID).Return(&repoUser, nil)
				path := fmt.Sprintf("/api/v1/users/%s", hresUser.Id)
				return &hresUser, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(mr MockRepository) (hres *schema.UserDetail, userpath string) {
				id := random.UUID()
				mr.user.EXPECT().GetUser(anyCtx{}, id).Return(nil, errors.New("Internal Server Error"))
				path := fmt.Sprintf("/api/v1/users/%s", id)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Bad Request: validate error: UUID",
			setup: func(mr MockRepository) (hres *schema.UserDetail, userpath string) {
				id := random.UUID()
				mr.user.EXPECT().GetUser(anyCtx{}, id).Return(nil, repository.ErrValidate)
				path := fmt.Sprintf("/api/v1/users/%s", id)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(_ MockRepository) (hres *schema.UserDetail, userpath string) {
				id := random.AlphaNumericN(36)
				path := fmt.Sprintf("/api/v1/users/%s", id)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupUserMock(t)

			var resBody *schema.UserDetail

			hresUsers, userpath := tt.setup(mr)

			statusCode, _ := doRequest(t, api, http.MethodGet, userpath, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (reqBody *schema.EditUserRequest, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(mr MockRepository) (*schema.EditUserRequest, string) {
				userID := random.UUID()
				userBio := random.AlphaNumeric()
				userCheck := random.Bool()

				reqBody := &schema.EditUserRequest{
					Bio:   &userBio,
					Check: &userCheck,
				}

				args := repository.UpdateUserArgs{
					Description: optional.FromPtr(&userBio),
					Check:       optional.FromPtr(&userCheck),
				}

				path := fmt.Sprintf("/api/v1/users/%s", userID)
				mr.user.EXPECT().UpdateUser(anyCtx{}, userID, &args).Return(nil)
				return reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "Success with description args(len=256)",
			setup: func(mr MockRepository) (*schema.EditUserRequest, string) {
				userID := random.UUID()
				userBio := strings.Repeat("a", 256)
				userCheck := random.Bool()

				reqBody := &schema.EditUserRequest{
					Bio:   &userBio,
					Check: &userCheck,
				}

				args := repository.UpdateUserArgs{
					Description: optional.FromPtr(&userBio),
					Check:       optional.FromPtr(&userCheck),
				}

				path := fmt.Sprintf("/api/v1/users/%s", userID)
				mr.user.EXPECT().UpdateUser(anyCtx{}, userID, &args).Return(nil)
				return reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "Conflict",
			setup: func(mr MockRepository) (*schema.EditUserRequest, string) {
				userID := random.UUID()
				userBio := random.AlphaNumeric()
				userCheck := random.Bool()

				reqBody := &schema.EditUserRequest{
					Bio:   &userBio,
					Check: &userCheck,
				}

				args := repository.UpdateUserArgs{
					Description: optional.FromPtr(&userBio),
					Check:       optional.FromPtr(&userCheck),
				}

				path := fmt.Sprintf("/api/v1/users/%s", userID)
				mr.user.EXPECT().UpdateUser(anyCtx{}, userID, &args).Return(repository.ErrAlreadyExists)
				return reqBody, path
			},
			statusCode: http.StatusConflict,
		},
		{
			name: "Not Found",
			setup: func(mr MockRepository) (*schema.EditUserRequest, string) {
				userID := random.UUID()
				userBio := random.AlphaNumeric()
				userCheck := random.Bool()

				reqBody := &schema.EditUserRequest{
					Bio:   &userBio,
					Check: &userCheck,
				}

				args := repository.UpdateUserArgs{
					Description: optional.FromPtr(&userBio),
					Check:       optional.FromPtr(&userCheck),
				}

				path := fmt.Sprintf("/api/v1/users/%s", userID)
				mr.user.EXPECT().UpdateUser(anyCtx{}, userID, &args).Return(repository.ErrNotFound)
				return reqBody, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: invalid userID",
			setup: func(_ MockRepository) (*schema.EditUserRequest, string) {
				path := fmt.Sprintf("/api/v1/users/%s", "invalid")
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: too long description(len>256)",
			setup: func(_ MockRepository) (*schema.EditUserRequest, string) {
				userID := random.UUID()
				userBio := strings.Repeat("a", 257)

				reqBody := &schema.EditUserRequest{
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
			t.Parallel()
			// Setup mock
			mr, api := setupUserMock(t)

			reqBody, path := tt.setup(mr)

			statusCode, _ := doRequest(t, api, http.MethodPatch, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}

func TestUserHandler_GetUserAccounts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (hres []*schema.Account, path string)
		statusCode int
	}{
		{
			name: "success",
			setup: func(mr MockRepository) (hres []*schema.Account, path string) {
				userID := random.UUID()
				accountKinds := rand.IntN((1<<domain.AccountLimit)-1) + 1
				//AccountLimit種類のうち、テストに使うものだけbitが立っている
				//例えば0(HOMEPAGE)と2(TWITTER)と7(ATCODER)なら10000101=133
				//0(bitがすべて立っていない)は除外

				rAccounts := []*domain.Account{}
				hAccounts := []*schema.Account{}

				for i := range int(domain.AccountLimit) {
					if (accountKinds>>i)%2 == 0 {
						continue
					}

					prRandom := random.Bool()
					raccount := domain.Account{
						ID:          random.UUID(),
						DisplayName: random.AlphaNumeric(),
						Type:        domain.AccountType(uint8(i)),
						PrPermitted: prRandom,
						URL:         random.AlphaNumeric(),
					}

					haccount := schema.Account{
						Id:          raccount.ID,
						DisplayName: raccount.DisplayName,
						PrPermitted: schema.PrPermitted(prRandom),
						Type:        schema.AccountType(raccount.Type),
						Url:         raccount.URL,
					}

					rAccounts = append(rAccounts, &raccount)
					hAccounts = append(hAccounts, &haccount)
				}

				mr.user.EXPECT().GetAccounts(anyCtx{}, userID).Return(rAccounts, nil)
				path = fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return hAccounts, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(mr MockRepository) (hres []*schema.Account, path string) {
				userID := random.UUID()
				mr.user.EXPECT().GetAccounts(anyCtx{}, userID).Return(nil, errors.New("Internal Server Error"))
				path = fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Bad Request: validate error: UUID",
			setup: func(mr MockRepository) (hres []*schema.Account, path string) {
				userID := random.UUID()
				mr.user.EXPECT().GetAccounts(anyCtx{}, userID).Return(nil, repository.ErrValidate)
				path = fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(_ MockRepository) (hres []*schema.Account, path string) {
				userID := random.AlphaNumericN(36)
				path = fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupUserMock(t)

			hresUsers, path := tt.setup(mr)

			var resBody []*schema.Account
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_GetUserAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (hres *schema.Account, path string)
		statusCode int
	}{
		{
			name: "success",
			setup: func(mr MockRepository) (hres *schema.Account, path string) {
				userID := random.UUID()
				prRandom := random.Bool()

				rAccount := domain.Account{
					ID:          random.UUID(),
					DisplayName: random.AlphaNumeric(),
					Type:        rand.N(domain.AccountLimit),
					PrPermitted: prRandom,
					URL:         random.AlphaNumeric(),
				}
				hAccount := schema.Account{
					Id:          rAccount.ID,
					DisplayName: rAccount.DisplayName,
					PrPermitted: schema.PrPermitted(prRandom),
					Type:        schema.AccountType(rAccount.Type),
					Url:         rAccount.URL,
				}

				mr.user.EXPECT().GetAccount(anyCtx{}, userID, rAccount.ID).Return(&rAccount, nil)
				path = fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, rAccount.ID)
				return &hAccount, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(mr MockRepository) (hres *schema.Account, path string) {
				userID := random.UUID()
				accountID := random.UUID()

				mr.user.EXPECT().GetAccount(anyCtx{}, userID, accountID).Return(nil, errors.New("Internal Server Error"))
				path = fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Bad Request: validate error: invalid userID",
			setup: func(_ MockRepository) (hres *schema.Account, path string) {
				userID := random.AlphaNumericN(36)
				accountID := random.UUID()

				path = fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(_ MockRepository) (hres *schema.Account, path string) {
				userID := random.UUID()
				accountID := random.AlphaNumericN(36)

				path = fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupUserMock(t)

			hresUsers, path := tt.setup(mr)

			var resBody *schema.Account
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_AddUserAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (reqBody *schema.AddAccountRequest, expectedResBody schema.Account, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(mr MockRepository) (*schema.AddAccountRequest, schema.Account, string) {
				userID := random.UUID()
				accountType := rand.N(domain.AccountLimit)

				reqBody := schema.AddAccountRequest{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: schema.PrPermitted(random.Bool()),
					Type:        schema.AccountType(accountType),
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

				expectedResBody := schema.Account{
					Id:          userID,
					DisplayName: reqBody.DisplayName,
					PrPermitted: reqBody.PrPermitted,
					Type:        reqBody.Type,
					Url:         reqBody.Url,
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				mr.user.EXPECT().CreateAccount(anyCtx{}, userID, &args).Return(&want, nil)
				return &reqBody, expectedResBody, path
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "Success: Account Type is 0",
			setup: func(mr MockRepository) (*schema.AddAccountRequest, schema.Account, string) {
				userID := random.UUID()

				reqBody := schema.AddAccountRequest{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: schema.PrPermitted(random.Bool()),
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

				expectedResBody := schema.Account{
					Id:          userID,
					DisplayName: reqBody.DisplayName,
					PrPermitted: reqBody.PrPermitted,
					Type:        reqBody.Type,
					Url:         reqBody.Url,
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				mr.user.EXPECT().CreateAccount(anyCtx{}, userID, &args).Return(&want, nil)
				return &reqBody, expectedResBody, path
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "Bad Request: DisplayName is empty",
			setup: func(_ MockRepository) (*schema.AddAccountRequest, schema.Account, string) {
				userID := random.UUID()
				accountType := rand.N(domain.AccountLimit)

				reqBody := schema.AddAccountRequest{
					DisplayName: "",
					PrPermitted: schema.PrPermitted(random.Bool()),
					Type:        schema.AccountType(accountType),
					Url:         random.AccountURLString(accountType),
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return &reqBody, schema.Account{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: Account Type is invalid",
			setup: func(_ MockRepository) (*schema.AddAccountRequest, schema.Account, string) {
				userID := random.UUID()

				reqBody := schema.AddAccountRequest{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: schema.PrPermitted(random.Bool()),
					Type:        schema.AccountType(domain.AccountLimit),
					Url:         random.RandURLString(),
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return &reqBody, schema.Account{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: UUID",
			setup: func(_ MockRepository) (*schema.AddAccountRequest, schema.Account, string) {
				userID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, schema.Account{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(_ MockRepository) (*schema.AddAccountRequest, schema.Account, string) {
				userID := random.AlphaNumericN(36)

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				return nil, schema.Account{}, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "internal error",
			setup: func(mr MockRepository) (*schema.AddAccountRequest, schema.Account, string) {
				userID := random.UUID()
				accountType := rand.N(domain.AccountLimit)

				reqBody := schema.AddAccountRequest{
					DisplayName: random.AlphaNumeric(),
					PrPermitted: schema.PrPermitted(random.Bool()),
					Type:        schema.AccountType(accountType),
					Url:         random.AccountURLString(accountType),
				}

				args := repository.CreateAccountArgs{
					DisplayName: reqBody.DisplayName,
					Type:        domain.AccountType(uint8(reqBody.Type)),
					URL:         reqBody.Url,
					PrPermitted: bool(reqBody.PrPermitted),
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts", userID)
				mr.user.EXPECT().CreateAccount(anyCtx{}, userID, &args).Return(nil, errors.New("internal error"))
				return &reqBody, schema.Account{}, path
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupUserMock(t)

			reqBody, res, path := tt.setup(mr)

			var resBody schema.Account
			statusCode, _ := doRequest(t, api, http.MethodPost, path, reqBody, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, res, resBody)
		})
	}
}

func TestUserHandler_EditUserAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (reqBody *schema.EditUserAccountRequest, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(mr MockRepository) (*schema.EditUserAccountRequest, string) {
				userID := random.UUID()
				accountID := random.UUID()
				accountType := rand.N(domain.AccountLimit)
				accountPermit := random.Bool()

				argsName := random.AlphaNumeric()
				argsPermit := schema.PrPermitted(accountPermit)
				argsType := schema.AccountType(accountType)
				argsURL := random.AccountURLString(domain.AccountType(accountType))

				reqBody := schema.EditUserAccountRequest{
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
				mr.user.EXPECT().UpdateAccount(anyCtx{}, userID, accountID, &args).Return(nil)
				return &reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "Not Found",
			setup: func(mr MockRepository) (*schema.EditUserAccountRequest, string) {
				userID := random.UUID()
				accountID := random.UUID()
				accountType := rand.N(domain.AccountLimit)
				accountPermit := random.Bool()

				argsName := random.AlphaNumeric()
				argsPermit := schema.PrPermitted(accountPermit)
				argsType := schema.AccountType(accountType)
				argsURL := random.AccountURLString(domain.AccountType(accountType))

				reqBody := schema.EditUserAccountRequest{
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
				mr.user.EXPECT().UpdateAccount(anyCtx{}, userID, accountID, &args).Return(repository.ErrNotFound)
				return &reqBody, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error: empty display name(but not nil)",
			setup: func(_ MockRepository) (*schema.EditUserAccountRequest, string) {
				userID := random.UUID()
				accountID := random.UUID()

				argsName := "" // empty but not nil

				reqBody := schema.EditUserAccountRequest{
					DisplayName: &argsName,
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)

				return &reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: too large account type",
			setup: func(_ MockRepository) (*schema.EditUserAccountRequest, string) {
				userID := random.UUID()
				accountID := random.UUID()

				argsType := schema.AccountType(domain.AccountLimit)

				reqBody := schema.EditUserAccountRequest{
					Type: &argsType,
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)

				return &reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: invalid url",
			setup: func(_ MockRepository) (*schema.EditUserAccountRequest, string) {
				userID := random.UUID()
				accountID := random.UUID()

				argsURL := random.AlphaNumeric()

				reqBody := schema.EditUserAccountRequest{
					Url: &argsURL,
				}

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)

				return &reqBody, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: nonUUID1",
			setup: func(_ MockRepository) (*schema.EditUserAccountRequest, string) {
				userID := random.AlphaNumericN(36)
				accountID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: nonUUID2",
			setup: func(_ MockRepository) (*schema.EditUserAccountRequest, string) {
				userID := random.UUID()
				accountID := random.AlphaNumericN(36)

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupUserMock(t)

			reqBody, path := tt.setup(mr)

			statusCode, _ := doRequest(t, api, http.MethodPatch, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}

func TestUserHandler_DeleteUserAccount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(mr MockRepository) string {
				userID := random.UUID()
				accountID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				mr.user.EXPECT().DeleteAccount(anyCtx{}, userID, accountID).Return(nil)
				return path
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "Forbidden",
			setup: func(mr MockRepository) string {
				userID := random.UUID()
				accountID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				mr.user.EXPECT().DeleteAccount(anyCtx{}, userID, accountID).Return(repository.ErrForbidden)
				return path
			},
			statusCode: http.StatusForbidden,
		},
		{
			name: "Not Found",
			setup: func(mr MockRepository) string {
				userID := random.UUID()
				accountID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				mr.user.EXPECT().DeleteAccount(anyCtx{}, userID, accountID).Return(repository.ErrNotFound)
				return path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error: nonUUID1",
			setup: func(_ MockRepository) string {
				userID := random.AlphaNumericN(36)
				accountID := random.UUID()

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error: nonUUID2",
			setup: func(_ MockRepository) string {
				userID := random.UUID()
				accountID := random.AlphaNumericN(36)

				path := fmt.Sprintf("/api/v1/users/%s/accounts/%s", userID, accountID)
				return path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupUserMock(t)

			path := tt.setup(mr)

			statusCode, _ := doRequest(t, api, http.MethodDelete, path, nil, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}

func TestUserHandler_GetUserProjects(t *testing.T) {
	makeProjects := func(t *testing.T, mr MockRepository, projectsLen int) (hres []*schema.UserProject, path string) {
		t.Helper()

		userID := random.UUID()

		repoProjects := []*domain.UserProject{}
		hresProjects := []*schema.UserProject{}

		for range projectsLen {
			//TODO: DurationはUserDurationを包含しているべき
			rproject := domain.UserProject{
				ID:           random.UUID(),
				Name:         random.AlphaNumeric(),
				Duration:     random.Duration(),
				UserDuration: random.Duration(),
			}

			hproject := schema.UserProject{
				Duration:     schema.ConvertDuration(rproject.Duration),
				Id:           rproject.ID,
				Name:         rproject.Name,
				UserDuration: schema.ConvertDuration(rproject.UserDuration),
			}

			repoProjects = append(repoProjects, &rproject)
			hresProjects = append(hresProjects, &hproject)
		}

		mr.user.EXPECT().GetProjects(anyCtx{}, userID).Return(repoProjects, nil)
		path = fmt.Sprintf("/api/v1/users/%s/projects", userID)
		return hresProjects, path
	}

	t.Parallel()

	tests := []struct {
		name       string
		setup      func(t *testing.T, mr MockRepository) (hres []*schema.UserProject, path string)
		statusCode int
	}{
		{
			name: "success 1",
			setup: func(t *testing.T, mr MockRepository) (hres []*schema.UserProject, path string) {
				return makeProjects(t, mr, 1)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 2",
			setup: func(t *testing.T, mr MockRepository) (hres []*schema.UserProject, path string) {
				return makeProjects(t, mr, 2)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 32",
			setup: func(t *testing.T, mr MockRepository) (hres []*schema.UserProject, path string) {
				return makeProjects(t, mr, 32)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Not Found",
			setup: func(t *testing.T, mr MockRepository) (hres []*schema.UserProject, path string) {
				userID := random.UUID()

				mr.user.EXPECT().GetProjects(anyCtx{}, userID).Return(nil, repository.ErrNotFound)
				path = fmt.Sprintf("/api/v1/users/%s/projects", userID)
				return nil, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error",
			setup: func(t *testing.T, _ MockRepository) (hres []*schema.UserProject, path string) {
				userID := random.AlphaNumericN(36)

				path = fmt.Sprintf("/api/v1/users/%s/projects", userID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupUserMock(t)

			hresUsers, path := tt.setup(t, mr)
			var resBody []*schema.UserProject
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_GetUserContests(t *testing.T) {
	makeContests := func(t *testing.T, mr MockRepository, contestsLen int) (hres []*schema.UserContest, path string) {
		t.Helper()

		userID := random.UUID()

		repoContests := []*domain.UserContest{}
		hresContests := []*schema.UserContest{}

		for range contestsLen {
			rcontest := domain.UserContest{
				ID:        random.UUID(),
				Name:      random.AlphaNumeric(),
				TimeStart: random.Time(),
				TimeEnd:   random.Time(),
				Teams: []*domain.ContestTeamWithoutMembers{
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
				[]schema.ContestTeamWithoutMembers{
					newContestTeamWithoutMembers(rcontest.Teams[0].ID, rcontest.Teams[0].Name, rcontest.Teams[0].Result),
				},
			)

			repoContests = append(repoContests, &rcontest)
			hresContests = append(hresContests, &hcontest)
		}

		mr.user.EXPECT().GetContests(anyCtx{}, userID).Return(repoContests, nil)
		path = fmt.Sprintf("/api/v1/users/%s/contests", userID)
		return hresContests, path
	}

	t.Parallel()

	tests := []struct {
		name       string
		setup      func(t *testing.T, mr MockRepository) (hres []*schema.UserContest, path string)
		statusCode int
	}{
		{
			name: "success 1",
			setup: func(t *testing.T, mr MockRepository) (hres []*schema.UserContest, path string) {
				return makeContests(t, mr, 1)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 2",
			setup: func(t *testing.T, mr MockRepository) (hres []*schema.UserContest, path string) {
				return makeContests(t, mr, 2)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 32",
			setup: func(t *testing.T, mr MockRepository) (hres []*schema.UserContest, path string) {
				return makeContests(t, mr, 32)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Not Found",
			setup: func(t *testing.T, mr MockRepository) (hres []*schema.UserContest, path string) {
				userID := random.UUID()

				mr.user.EXPECT().GetContests(anyCtx{}, userID).Return(nil, repository.ErrNotFound)
				path = fmt.Sprintf("/api/v1/users/%s/contests", userID)
				return nil, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error",
			setup: func(t *testing.T, _ MockRepository) (hres []*schema.UserContest, path string) {
				userID := random.AlphaNumericN(36)

				path = fmt.Sprintf("/api/v1/users/%s/contests", userID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupUserMock(t)

			hresUsers, path := tt.setup(t, mr)
			var resBody []*schema.UserContest
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_GetUserGroups(t *testing.T) {
	makeGroups := func(mr MockRepository, groupsLen int) (hres []*schema.UserGroup, path string) {
		userID := random.UUID()

		repoGroups := []*domain.UserGroup{}
		hresGroups := []*schema.UserGroup{}

		for range groupsLen {
			rgroup := domain.UserGroup{
				ID:       random.UUID(),
				Name:     random.AlphaNumeric(),
				Duration: random.Duration(),
			}

			hgroup := schema.UserGroup{
				Duration: schema.ConvertDuration(rgroup.Duration),
				Id:       rgroup.ID,
				Name:     rgroup.Name,
			}

			repoGroups = append(repoGroups, &rgroup)
			hresGroups = append(hresGroups, &hgroup)
		}

		mr.user.EXPECT().GetGroupsByUserID(anyCtx{}, userID).Return(repoGroups, nil)
		path = fmt.Sprintf("/api/v1/users/%s/groups", userID)
		return hresGroups, path
	}

	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (hres []*schema.UserGroup, path string)
		statusCode int
	}{
		{
			name: "success 0",
			setup: func(mr MockRepository) (hres []*schema.UserGroup, path string) {
				return makeGroups(mr, 0)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 1",
			setup: func(mr MockRepository) (hres []*schema.UserGroup, path string) {
				return makeGroups(mr, 1)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 2",
			setup: func(mr MockRepository) (hres []*schema.UserGroup, path string) {
				return makeGroups(mr, 2)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 32",
			setup: func(mr MockRepository) (hres []*schema.UserGroup, path string) {
				return makeGroups(mr, 32)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Not Found",
			setup: func(mr MockRepository) (hres []*schema.UserGroup, path string) {
				userID := random.UUID()

				mr.user.EXPECT().GetGroupsByUserID(anyCtx{}, userID).Return(nil, repository.ErrNotFound)
				path = fmt.Sprintf("/api/v1/users/%s/groups", userID)
				return nil, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error",
			setup: func(_ MockRepository) (hres []*schema.UserGroup, path string) {
				userID := random.AlphaNumericN(36)

				path = fmt.Sprintf("/api/v1/users/%s/groups", userID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupUserMock(t)

			hresUsers, path := tt.setup(mr)
			var resBody []*schema.UserGroup
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_GetUserEvents(t *testing.T) {
	makeEvents := func(mr MockRepository, eventsLen int) (hres []*schema.Event, path string) {
		userID := random.UUID()

		repoEvents := []*domain.Event{}
		hresEvents := []*schema.Event{}

		for range eventsLen {
			timeStart, timeEnd := random.SinceAndUntil()

			revent := domain.Event{
				ID:        random.UUID(),
				Name:      random.AlphaNumeric(),
				TimeStart: timeStart,
				TimeEnd:   timeEnd,
			}

			hevent := schema.Event{
				Duration: schema.Duration{
					Since: timeStart,
					Until: &timeEnd,
				},
				Id:   revent.ID,
				Name: revent.Name,
			}

			repoEvents = append(repoEvents, &revent)
			hresEvents = append(hresEvents, &hevent)
		}

		mr.event.EXPECT().GetUserEvents(anyCtx{}, userID).Return(repoEvents, nil)
		path = fmt.Sprintf("/api/v1/users/%s/events", userID)
		return hresEvents, path
	}

	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (hres []*schema.Event, path string)
		statusCode int
	}{
		{
			name: "success 0",
			setup: func(mr MockRepository) (hres []*schema.Event, path string) {
				return makeEvents(mr, 0)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 1",
			setup: func(mr MockRepository) (hres []*schema.Event, path string) {
				return makeEvents(mr, 1)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 2",
			setup: func(mr MockRepository) (hres []*schema.Event, path string) {
				return makeEvents(mr, 2)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success 32",
			setup: func(mr MockRepository) (hres []*schema.Event, path string) {
				return makeEvents(mr, 32)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Not Found",
			setup: func(mr MockRepository) (hres []*schema.Event, path string) {
				userID := random.UUID()

				mr.event.EXPECT().GetUserEvents(anyCtx{}, userID).Return(nil, repository.ErrNotFound)
				path = fmt.Sprintf("/api/v1/users/%s/events", userID)
				return nil, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error",
			setup: func(_ MockRepository) (hres []*schema.Event, path string) {
				userID := random.AlphaNumericN(36)

				path = fmt.Sprintf("/api/v1/users/%s/events", userID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupUserMock(t)

			hresUsers, path := tt.setup(mr)
			var resBody []*schema.Event
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_GetMe(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (hres *schema.UserDetail, path string, header map[string]string)
		statusCode int
	}{
		{
			name: "success",
			setup: func(mr MockRepository) (hres *schema.UserDetail, path string, header map[string]string) {
				username := random.AlphaNumeric()
				header = map[string]string{
					"X-Forwarded-User": username,
				}

				userID := random.UUID()

				ruser := domain.NewUser(userID, username, random.AlphaNumeric(), random.Bool())
				rusers := []*domain.User{ruser}
				mr.user.EXPECT().GetUsers(anyCtx{}, &repository.GetUsersArgs{
					Name: optional.From(ruser.Name),
				}).Return(rusers, nil)

				accountType := rand.N(domain.AccountLimit)
				ruserDetail := domain.UserDetail{
					User:  *ruser,
					Bio:   random.AlphaNumeric(),
					State: rand.N(domain.TraqStateLimit),
					Accounts: []*domain.Account{
						{
							ID:          random.UUID(),
							DisplayName: random.AlphaNumeric(),
							Type:        accountType,
							PrPermitted: random.Bool(),
							URL:         random.AccountURLString(accountType),
						},
					},
				}
				mr.user.EXPECT().GetUser(anyCtx{}, userID).Return(&ruserDetail, nil)

				haccounts := []schema.Account{}
				for _, account := range ruserDetail.Accounts {
					haccounts = append(haccounts, schema.Account{
						Id:          account.ID,
						DisplayName: account.DisplayName,
						PrPermitted: schema.PrPermitted(account.PrPermitted),
						Type:        schema.AccountType(account.Type),
						Url:         account.URL,
					})
				}

				huser := schema.UserDetail{
					Id:       userID,
					Name:     ruser.Name,
					RealName: ruser.RealName(),
					Accounts: haccounts,
					Bio:      ruserDetail.Bio,
					State:    schema.UserAccountState(ruserDetail.State),
				}

				path = "/api/v1/users/me"

				return &huser, path, header
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Unauthorized",
			setup: func(mr MockRepository) (hres *schema.UserDetail, path string, header map[string]string) {
				header = map[string]string{}
				path = "/api/v1/users/me"
				return nil, path, header
			},
			statusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupUserMock(t)

			hresUsers, path, header := tt.setup(mr)
			var resBody *schema.UserDetail
			statusCode, _ := doRequestWithHeader(t, api, http.MethodGet, path, nil, &resBody, header)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}
