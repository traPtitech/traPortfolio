package handler_test

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
)

func TestUserHandler_GetAll(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(th *handler.TestHandlers) (hres []*handler.User, path string)
		statusCode int
	}{
		{
			name: "success",
			setup: func(th *handler.TestHandlers) (hres []*handler.User, path string) {

				casenum := 2
				repoUsers := []*domain.User{}
				hresUsers := []*handler.User{}

				for i := 0; i < casenum; i++ {
					ruser := domain.User{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(rand.Intn(30) + 1),
						RealName: random.AlphaNumeric(rand.Intn(30) + 1),
					}
					huser := handler.User{
						Id:       ruser.ID,
						Name:     ruser.Name,
						RealName: ruser.RealName,
					}

					repoUsers = append(repoUsers, &ruser)
					hresUsers = append(hresUsers, &huser)

				}

				th.Service.MockUserService.EXPECT().GetUsers(gomock.Any()).Return(repoUsers, nil)
				return hresUsers, "/api/v1/users"
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(th *handler.TestHandlers) (hres []*handler.User, path string) {
				th.Service.MockUserService.EXPECT().GetUsers(gomock.Any()).Return(nil, errors.New("Internal Server Error"))
				return nil, "/api/v1/users"
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			ctrl := gomock.NewController(t)
			handlers := SetupTestHandlers(t, ctrl)

			hresUsers, path := tt.setup(&handlers)

			var resBody []*handler.User
			statusCode, _ := doRequest(t, handlers.API, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)
		})
	}
}

func TestUserHandler_GetByID(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(th *handler.TestHandlers) (hres *handler.UserDetail, userpath string)
		statusCode int
	}{
		{
			name: "success random",
			setup: func(th *handler.TestHandlers) (hres *handler.UserDetail, userpath string) {

				const accountNum int = 9
				rAccounts := []*domain.Account{}
				hAccounts := []handler.Account{}

				for i := 0; i < accountNum; i++ {
					prRandom := false
					if rand.Intn(2) == 1 {
						prRandom = true
					}

					raccount := domain.Account{
						ID:          random.UUID(),
						Name:        random.AlphaNumeric(rand.Intn(30) + 1),
						Type:        uint(rand.Intn(int(domain.AccountLimit))),
						PrPermitted: prRandom,
						URL:         random.AlphaNumeric(rand.Intn(30) + 1),
					}

					haccount := handler.Account{
						Id:          raccount.ID,
						Name:        raccount.Name,
						PrPermitted: handler.PrPermitted(prRandom),
						Type:        handler.AccountType(raccount.Type),
						Url:         raccount.URL,
					}

					rAccounts = append(rAccounts, &raccount)
					hAccounts = append(hAccounts, haccount)
				}

				ruser := domain.UserDetail{

					User: domain.User{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(rand.Intn(30) + 1),
						RealName: random.AlphaNumeric(rand.Intn(30) + 1),
					},
					State:    domain.TraQState(uint8(rand.Intn(int(domain.TraqStateLimit)))),
					Bio:      random.AlphaNumeric(rand.Intn(256) + 1),
					Accounts: rAccounts,
				}

				huser := handler.UserDetail{
					User: handler.User{
						Id:       ruser.User.ID,
						Name:     ruser.User.Name,
						RealName: ruser.User.RealName,
					},
					Accounts: hAccounts,
					Bio:      ruser.Bio,
					State:    handler.UserAccountState(ruser.State),
				}

				repoUser := &ruser
				hresUser := &huser

				th.Service.MockUserService.EXPECT().GetUser(gomock.Any(), ruser.User.ID).Return(repoUser, nil)
				path := fmt.Sprintf("/api/v1/users/%s", ruser.User.ID)
				return hresUser, path
			},
			statusCode: http.StatusOK,
		},

		{
			name: "internal error",
			setup: func(th *handler.TestHandlers) (hres *handler.UserDetail, userpath string) {
				id := random.UUID()
				th.Service.MockUserService.EXPECT().GetUser(gomock.Any(), id).Return(nil, errors.New("Internal Server Error"))
				path := fmt.Sprintf("/api/v1/users/%s", id)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			ctrl := gomock.NewController(t)
			handlers := SetupTestHandlers(t, ctrl)

			var resBody *handler.UserDetail

			hresUsers, userpath := tt.setup(&handlers)

			statusCode, _ := doRequest(t, handlers.API, http.MethodGet, userpath, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresUsers, resBody)

		})
	}
}

func TestUserHandler_Update(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(th *handler.TestHandlers) (reqBody *handler.EditUser, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(th *handler.TestHandlers) (*handler.EditUser, string) {

				/*eventID := random.UUID()
				eventLevelUint := (uint)(rand.Intn(domain.EventLevelLimit))
				eventLevelHandler := handler.EventLevel(eventLevelUint)
				eventLevelDomain := domain.EventLevel(eventLevelUint)*/

				userID := random.UUID()
				userBio := random.AlphaNumeric(rand.Intn(30) + 1)
				userCheck := false
				if rand.Intn(2) == 1 {
					userCheck = true
				}

				reqBody := &handler.EditUser{
					Bio:   &userBio,
					Check: &userCheck,
				}

				args := repository.UpdateUserArgs{
					Description: optional.StringFrom(&userBio),
					Check:       optional.BoolFrom(&userCheck),
				}

				path := fmt.Sprintf("/api/v1/users/%s", userID)
				th.Service.MockUserService.EXPECT().Update(gomock.Any(), userID, &args).Return(nil)
				return reqBody, path
			},
			statusCode: http.StatusNoContent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			ctrl := gomock.NewController(t)
			handlers := SetupTestHandlers(t, ctrl)

			reqBody, path := tt.setup(&handlers)

			statusCode, _ := doRequest(t, handlers.API, http.MethodPatch, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}

/*

func TestUserHandler_GetAccounts(t *testing.T) {
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
			if err := handler.GetAccounts(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserHandler.GetAccounts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHandler_GetAccount(t *testing.T) {
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
			if err := handler.GetAccount(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserHandler.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHandler_AddAccount(t *testing.T) {
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
			if err := handler.AddAccount(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserHandler.AddAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHandler_PatchAccount(t *testing.T) {
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
			if err := handler.PatchAccount(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserHandler.PatchAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
			if err := handler.DeleteAccount(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserHandler.DeleteAccount() error = %v, wantErr %v", err, tt.wantErr)
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
			if err := handler.GetProjects(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserHandler.GetProjects() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHandler_GetContests(t *testing.T) {
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
			if err := handler.GetContests(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserHandler.GetContests() error = %v, wantErr %v", err, tt.wantErr)
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
			if err := handler.GetGroupsByUserID(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserHandler.GetGroupsByUserID() error = %v, wantErr %v", err, tt.wantErr)
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
			if err := handler.GetEvents(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("UserHandler.GetEvents() error = %v, wantErr %v", err, tt.wantErr)
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

func Test_newContestTeamWithContestName(t *testing.T) {
	type args struct {
		contestTeam ContestTeam
		contestName string
	}
	tests := []struct {
		name string
		args args
		want ContestTeamWithContestName
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newContestTeamWithContestName(tt.args.contestTeam, tt.args.contestName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newContestTeamWithContestName() = %v, want %v", got, tt.want)
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
