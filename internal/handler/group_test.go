package handler

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/handler/schema"
	"github.com/traPtitech/traPortfolio/internal/pkgs/random"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository/mock_repository"
	"go.uber.org/mock/gomock"
)

func setupGroupMock(t *testing.T) (MockRepository, API) {
	t.Helper()

	ctrl := gomock.NewController(t)
	user := mock_repository.NewMockUserRepository(ctrl)
	group := mock_repository.NewMockGroupRepository(ctrl)
	mr := MockRepository{user: user, group: group}
	api := NewAPI(nil, nil, nil, nil, nil, NewGroupHandler(group, user))

	return mr, api
}

func TestGroupHandler_GetGroups(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (hres []*schema.Group, path string)
		statusCode int
	}{
		{
			name: "success",
			setup: func(mr MockRepository) (hres []*schema.Group, path string) {
				casenum := 2
				repoGroups := []*domain.Group{}
				hresGroups := []*schema.Group{}

				for range casenum {
					rgroup := domain.Group{
						ID:   random.UUID(),
						Name: random.AlphaNumeric(),
					}

					hgroup := schema.Group{
						Id:   rgroup.ID,
						Name: rgroup.Name,
					}

					repoGroups = append(repoGroups, &rgroup)
					hresGroups = append(hresGroups, &hgroup)
				}

				mr.group.EXPECT().GetGroups(anyCtx{}).Return(repoGroups, nil)
				return hresGroups, "/api/v1/groups"
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(mr MockRepository) (hres []*schema.Group, path string) {
				mr.group.EXPECT().GetGroups(anyCtx{}).Return(nil, errors.New("Internal Server Error"))
				return nil, "/api/v1/groups"
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			s, api := setupGroupMock(t)

			hresGroups, path := tt.setup(s)

			var resBody []*schema.Group
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresGroups, resBody)
		})
	}
}

func TestGroupHandler_GetGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(mr MockRepository) (hres *schema.GroupDetail, path string)
		statusCode int
	}{
		{
			name: "success",
			setup: func(mr MockRepository) (hres *schema.GroupDetail, path string) {
				rgroupAdmins := []*domain.User{}
				hgroupAdmins := []schema.User{}

				adminLen := rand.IntN(256)
				for range adminLen {
					rgroupAdmin := domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool())

					hgroupAdmin := schema.User{
						Id:       rgroupAdmin.ID,
						Name:     rgroupAdmin.Name,
						RealName: rgroupAdmin.RealName(),
					}

					rgroupAdmins = append(rgroupAdmins, rgroupAdmin)
					hgroupAdmins = append(hgroupAdmins, hgroupAdmin)
				}

				rgroupMembers := []*domain.UserWithDuration{}
				hgroupMembers := []schema.GroupMember{}

				groupLen := rand.IntN(256)
				for range groupLen {
					rgroupmember := domain.UserWithDuration{
						User:     *domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
						Duration: random.Duration(),
					}

					hgroupmember := schema.GroupMember{
						Duration: schema.ConvertDuration(rgroupmember.Duration),
						Id:       rgroupmember.User.ID,
						Name:     rgroupmember.User.Name,
						RealName: rgroupmember.User.RealName(),
					}

					rgroupMembers = append(rgroupMembers, &rgroupmember)
					hgroupMembers = append(hgroupMembers, hgroupmember)
				}

				rgroup := domain.GroupDetail{
					ID:          random.UUID(),
					Name:        random.AlphaNumeric(),
					Links:       random.Array(random.RandURLString, 1, 3),
					Admin:       rgroupAdmins,
					Members:     rgroupMembers,
					Description: random.AlphaNumeric(),
				}

				hgroup := schema.GroupDetail{
					Description: rgroup.Description,
					Id:          rgroup.ID,
					Admin:       hgroupAdmins,
					Links:       rgroup.Links,
					Members:     hgroupMembers,
					Name:        rgroup.Name,
				}

				mr.group.EXPECT().GetGroup(anyCtx{}, rgroup.ID).Return(&rgroup, nil)
				mr.user.EXPECT().GetUsers(anyCtx{}, &repository.GetUsersArgs{}).Return(rgroupAdmins, nil)
				path = fmt.Sprintf("/api/v1/groups/%s", rgroup.ID)
				return &hgroup, path
			},
			statusCode: http.StatusOK,
		},
		{
			name: "internal error",
			setup: func(mr MockRepository) (hres *schema.GroupDetail, path string) {
				groupID := random.UUID()
				mr.group.EXPECT().GetGroup(anyCtx{}, groupID).Return(nil, errors.New("Internal Server Error"))
				path = fmt.Sprintf("/api/v1/groups/%s", groupID)
				return nil, path
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "forbidden",
			setup: func(mr MockRepository) (hres *schema.GroupDetail, path string) {
				groupID := random.UUID()
				mr.group.EXPECT().GetGroup(anyCtx{}, groupID).Return(nil, repository.ErrForbidden)
				path = fmt.Sprintf("/api/v1/groups/%s", groupID)
				return nil, path
			},
			statusCode: http.StatusForbidden,
		},
		{
			name: "not found",
			setup: func(mr MockRepository) (hres *schema.GroupDetail, path string) {
				groupID := random.UUID()
				mr.group.EXPECT().GetGroup(anyCtx{}, groupID).Return(nil, repository.ErrNotFound)
				path = fmt.Sprintf("/api/v1/groups/%s", groupID)
				return nil, path
			},
			statusCode: http.StatusNotFound,
		},
		{
			name: "Bad Request: validate error nonUUID",
			setup: func(_ MockRepository) (hres *schema.GroupDetail, path string) {
				groupID := random.AlphaNumericN(36)
				path = fmt.Sprintf("/api/v1/groups/%s", groupID)
				return nil, path
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			mr, api := setupGroupMock(t)

			hresGroups, path := tt.setup(mr)

			var resBody *schema.GroupDetail
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, hresGroups, resBody)
		})
	}
}
