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

				s.EXPECT().GetAllGroups(anyCtx{}).Return(repoGroups, nil)
				return hresGroups, "/api/v1/groups"
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(s *mock_service.MockGroupService) (hres []*Group, path string) {

				s.EXPECT().GetAllGroups(anyCtx{}).Return(nil, errors.New("Internal Server Error"))
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

				rgroupAdmins := []*domain.User{}
				hgroupAdmins := []User{}

				adminLen := rand.Intn(256)
				for i := 0; i < adminLen; i++ {
					rgroupAdmin := domain.User{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(),
						RealName: random.AlphaNumeric(),
					}

					hgroupAdmin := User{
						Id:       rgroupAdmin.ID,
						Name:     rgroupAdmin.Name,
						RealName: rgroupAdmin.RealName,
					}

					rgroupAdmins = append(rgroupAdmins, &rgroupAdmin)
					hgroupAdmins = append(hgroupAdmins, hgroupAdmin)
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
						Duration: ConvertDuration(rgroupmember.Duration),
						Id:       rgroupmember.ID,
						Name:     rgroupmember.Name,
						RealName: rgroupmember.RealName,
					}

					rgroupMembers = append(rgroupMembers, &rgroupmember)
					hgroupMembers = append(hgroupMembers, hgroupmember)
				}

				rgroup := domain.GroupDetail{
					ID:          random.UUID(),
					Name:        random.AlphaNumeric(),
					Link:        random.AlphaNumeric(),
					Admin:       rgroupAdmins,
					Members:     rgroupMembers,
					Description: random.AlphaNumeric(),
				}

				hgroup := GroupDetail{
					Description: rgroup.Description,
					Id:          rgroup.ID,
					Admin:       hgroupAdmins,
					Link:        rgroup.Link,
					Members:     hgroupMembers,
					Name:        rgroup.Name,
				}

				s.EXPECT().GetGroup(anyCtx{}, rgroup.ID).Return(&rgroup, nil)
				path = fmt.Sprintf("/api/v1/groups/%s", rgroup.ID)
				return &hgroup, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(s *mock_service.MockGroupService) (hres *GroupDetail, path string) {
				groupID := random.UUID()
				s.EXPECT().GetGroup(anyCtx{}, groupID).Return(nil, errors.New("Internal Server Error"))
				path = fmt.Sprintf("/api/v1/groups/%s", groupID)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "forbidden",
			setup: func(s *mock_service.MockGroupService) (hres *GroupDetail, path string) {
				groupID := random.UUID()
				s.EXPECT().GetGroup(anyCtx{}, groupID).Return(nil, repository.ErrForbidden)
				path = fmt.Sprintf("/api/v1/groups/%s", groupID)
				return nil, path
			},
			statusCode: http.StatusForbidden,
		},
		{
			name: "not found",
			setup: func(s *mock_service.MockGroupService) (hres *GroupDetail, path string) {
				groupID := random.UUID()
				s.EXPECT().GetGroup(anyCtx{}, groupID).Return(nil, repository.ErrNotFound)
				path = fmt.Sprintf("/api/v1/groups/%s", groupID)
				return nil, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(_ *mock_service.MockGroupService) (hres *GroupDetail, path string) {
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
