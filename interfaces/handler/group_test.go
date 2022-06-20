package handler

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/service/mock_service"
	"github.com/traPtitech/traPortfolio/util/random"
)

func setupGroupMock(t *testing.T) (*mock_service.MockGroupService, API) {
	t.Helper()

	ctrl := gomock.NewController(t)
	s := mock_service.NewMockGroupService(ctrl)
	api := NewAPI(nil, nil, nil, nil, nil, NewGroupHandler(s))

	return s, api
}

func TestGroupHandler_GetGroups(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockGroupService) (hres []*Group, path string)
		statusCode int
	}{
		{
			name: "success",
			setup: func(s *mock_service.MockGroupService) (hres []*Group, path string) {

				casenum := 2
				repoGroups := []*domain.Group{}
				hresGroups := []*Group{}

				for i := 0; i < casenum; i++ {

					rgroup := domain.Group{
						ID:   random.UUID(),
						Name: random.AlphaNumeric(),
					}

					hgroup := Group{
						Id:   rgroup.ID,
						Name: rgroup.Name,
					}

					repoGroups = append(repoGroups, &rgroup)
					hresGroups = append(hresGroups, &hgroup)

				}

				s.EXPECT().GetAllGroups(gomock.Any()).Return(repoGroups, nil)
				return hresGroups, "/api/v1/groups"
			},
			statusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupGroupMock(t)

			hresGroups, path := tt.setup(s)

			var resBody []*Group
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresGroups, resBody)
		})
	}
}

/*
func TestGroupHandler_GetGroups(t *testing.T) {
	type fields struct {
		srv service.GroupService
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
			h := &GroupHandler{
				srv: tt.fields.srv,
			}
			if err := h.GetGroups(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("GroupHandler.GetGroups() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

*/

/*

func TestGroupHandler_GetGroup(t *testing.T) {
	type fields struct {
		srv service.GroupService
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
			h := &GroupHandler{
				srv: tt.fields.srv,
			}
			if err := h.GetGroup(tt.args._c); (err != nil) != tt.wantErr {
				t.Errorf("GroupHandler.GetGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_formatGetGroup(t *testing.T) {
	type args struct {
		group *domain.GroupDetail
	}
	tests := []struct {
		name string
		args args
		want GroupDetail
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatGetGroup(tt.args.group); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("formatGetGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newGroupMember(t *testing.T) {
	type args struct {
		user     User
		Duration YearWithSemesterDuration
	}
	tests := []struct {
		name string
		args args
		want GroupMember
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newGroupMember(tt.args.user, tt.args.Duration); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newGroupMember() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newGroupDetail(t *testing.T) {
	type args struct {
		group   Group
		desc    string
		leader  User
		link    string
		members []GroupMember
	}
	tests := []struct {
		name string
		args args
		want GroupDetail
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newGroupDetail(tt.args.group, tt.args.desc, tt.args.leader, tt.args.link, tt.args.members); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newGroupDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
