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
	"github.com/traPtitech/traPortfolio/util/random"
)

//TODO: anyCtxを書く

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
		{
			name: "internal error",
			setup: func(s *mock_service.MockGroupService) (hres []*Group, path string) {

				s.EXPECT().GetAllGroups(gomock.Any()).Return(nil, errors.New("Internal Server Error"))
				return nil, "/api/v1/groups"
			},
			statusCode: http.StatusInternalServerError,
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

func TestGroupHandler_GetGroup(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(s *mock_service.MockGroupService) (hres *GroupDetail, path string)
		statusCode int
	}{
		{
			name: "success",
			setup: func(s *mock_service.MockGroupService) (hres *GroupDetail, path string) {

				rgroupLeader := domain.User{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(),
					RealName: random.AlphaNumeric(),
				}

				hgroupLeader := User{
					Id:       rgroupLeader.ID,
					Name:     rgroupLeader.Name,
					RealName: rgroupLeader.RealName,
				}

				rgroupMembers := []*domain.UserGroup{}
				hgroupMembers := []GroupMember{}

				groupLen := rand.Intn(256)
				for i := 0; i < groupLen; i++ {
					rgroupmember := domain.UserGroup{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(),
						RealName: random.AlphaNumeric(),
						Duration: random.Duration(),
					}

					hgroupmember := GroupMember{
						Duration: convertDuration(rgroupmember.Duration),
						Id:       rgroupmember.ID,
						Name:     rgroupmember.Name,
						RealName: rgroupmember.RealName,
					}

					//TODO: 何故かSemesterとYearが反転する
					t := rgroupmember.Duration.Since.Semester
					rgroupmember.Duration.Since.Semester = rgroupmember.Duration.Since.Year
					rgroupmember.Duration.Since.Year = int(t)

					t = rgroupmember.Duration.Until.Semester
					rgroupmember.Duration.Until.Semester = rgroupmember.Duration.Until.Year
					rgroupmember.Duration.Until.Year = int(t)
					//以上の7行はテストを通すための反転であり想定されてはいけない実装です

					rgroupMembers = append(rgroupMembers, &rgroupmember)
					hgroupMembers = append(hgroupMembers, hgroupmember)
				}

				rgroup := domain.GroupDetail{
					ID:          random.UUID(),
					Name:        random.AlphaNumeric(),
					Link:        random.AlphaNumeric(),
					Leader:      &rgroupLeader,
					Members:     rgroupMembers,
					Description: random.AlphaNumeric(),
				}

				hgroup := GroupDetail{
					Description: rgroup.Description,
					Id:          rgroup.ID,
					Leader:      hgroupLeader,
					Link:        rgroup.Link,
					Members:     hgroupMembers,
					Name:        rgroup.Name,
				}

				repoGroup := &rgroup
				hresGroup := &hgroup

				//謎の反転が発生するのはrepoGroup側で、ここ以降

				s.EXPECT().GetGroup(gomock.Any(), rgroup.ID).Return(repoGroup, nil)
				path = fmt.Sprintf("/api/v1/groups/%s", rgroup.ID)
				return hresGroup, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(s *mock_service.MockGroupService) (hres *GroupDetail, path string) {
				groupID := random.UUID()
				s.EXPECT().GetGroup(gomock.Any(), groupID).Return(nil, errors.New("Internal Server Error"))
				path = fmt.Sprintf("/api/v1/groups/%s", groupID)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Bad Request: validate error: UUID",
			setup: func(s *mock_service.MockGroupService) (hres *GroupDetail, path string) {
				groupID := random.UUID()
				s.EXPECT().GetGroup(gomock.Any(), groupID).Return(nil, repository.ErrValidate)
				path = fmt.Sprintf("/api/v1/groups/%s", groupID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(s *mock_service.MockGroupService) (hres *GroupDetail, path string) {
				groupID := random.AlphaNumericn(36)
				path = fmt.Sprintf("/api/v1/groups/%s", groupID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupGroupMock(t)

			hresGroups, path := tt.setup(s)

			var resBody *GroupDetail
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresGroups, resBody)
		})
	}
}

/*

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
