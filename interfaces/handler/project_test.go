package handler

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/traPtitech/traPortfolio/interfaces/handler/schema"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/optional"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/service/mock_service"
	"github.com/traPtitech/traPortfolio/util/random"
)

var (
	errInternal = errors.New("internal error")
)

func setupProjectMock(t *testing.T) (*mock_service.MockProjectService, API) {
	t.Helper()

	ctrl := gomock.NewController(t)
	s := mock_service.NewMockProjectService(ctrl)
	api := NewAPI(nil, nil, NewProjectHandler(s), nil, nil, nil)

	return s, api
}

func makeCreateProjectRequest(t *testing.T, description string, since schema.YearWithSemester, until *schema.YearWithSemester, name string, link string) *schema.CreateProjectJSONRequestBody {
	t.Helper()
	return &schema.CreateProjectJSONRequestBody{
		Description: description,
		Duration: schema.YearWithSemesterDuration{
			Since: since,
			Until: until,
		},
		Name: name,
		Link: &link,
	}
}

// func makeAddProjectMembersRequest(members []ProjectMember) *AddProjectMembersJSONRequestBody {
// 	ret := &AddProjectMembersJSONRequestBody{}
// 	for _, v := range members {
// 		ret.Members = append(ret.Members, MemberIDWithYearWithSemesterDuration{
// 			Duration: YearWithSemesterDuration{
// 				Since: v.Duration.Since,
// 				Until: v.Duration.Until,
// 			},
// 			UserId: v.Id,
// 		})
// 	}
// 	return ret
// }

func TestProjectHandler_GetProjects(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(s *mock_service.MockProjectService) ([]*schema.Project, string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockProjectService) ([]*schema.Project, string) {
				duration := random.Duration()
				repo := []*domain.Project{
					{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(),
						Duration: duration,
					},
					{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(),
						Duration: duration,
					},
				}

				var reqBody []*schema.Project
				for _, v := range repo {
					reqBody = append(reqBody, &schema.Project{
						Duration: schema.YearWithSemesterDuration{
							Since: schema.YearWithSemester{
								Year:     v.Duration.Since.Year,
								Semester: schema.Semester(v.Duration.Since.Semester),
							},
							Until: &schema.YearWithSemester{
								Year:     v.Duration.Until.Year,
								Semester: schema.Semester(v.Duration.Until.Semester),
							},
						},
						Id:   v.ID,
						Name: v.Name,
					})
				}

				s.EXPECT().GetProjects(anyCtx{}).Return(repo, nil)
				return reqBody, "/api/v1/projects"
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Internal Error",
			setup: func(s *mock_service.MockProjectService) ([]*schema.Project, string) {
				s.EXPECT().GetProjects(anyCtx{}).Return(nil, errInternal)
				return nil, "/api/v1/projects"
			},
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupProjectMock(t)

			expectedHres, path := tt.setup(s)

			hres := []*schema.Project(nil)
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &hres)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, expectedHres, hres)
		})
	}
}

func TestProjectHandler_GetProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(s *mock_service.MockProjectService) (*schema.ProjectDetail, string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockProjectService) (*schema.ProjectDetail, string) {
				duration := random.Duration()
				projectID := random.UUID()
				repo := domain.ProjectDetail{
					Project: domain.Project{
						ID:       projectID,
						Name:     random.AlphaNumeric(),
						Duration: duration,
					},
					Description: random.AlphaNumeric(),
					Link:        random.RandURLString(),
					Members: []*domain.UserWithDuration{
						{
							User:     *domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
							Duration: random.Duration(),
						},
					},
				}

				var members []schema.ProjectMember
				for _, v := range repo.Members {
					members = append(members, schema.ProjectMember{
						Duration: schema.YearWithSemesterDuration{
							Since: schema.YearWithSemester{
								Year:     v.Duration.Since.Year,
								Semester: schema.Semester(v.Duration.Since.Semester),
							},
							Until: &schema.YearWithSemester{
								Year:     v.Duration.Until.Year,
								Semester: schema.Semester(v.Duration.Until.Semester),
							},
						},
						Id:       v.User.ID,
						Name:     v.User.Name,
						RealName: v.User.RealName(),
					})
				}
				reqBody := &schema.ProjectDetail{
					Description: repo.Description,
					Duration: schema.YearWithSemesterDuration{
						Since: schema.YearWithSemester{
							Semester: schema.Semester(repo.Duration.Since.Semester),
							Year:     repo.Duration.Since.Year,
						},
						Until: &schema.YearWithSemester{
							Semester: schema.Semester(repo.Duration.Until.Semester),
							Year:     repo.Duration.Until.Year,
						},
					},
					Id:      repo.ID,
					Link:    repo.Link,
					Members: members,
					Name:    repo.Name,
				}

				s.EXPECT().GetProject(anyCtx{}, projectID).Return(&repo, nil)
				return reqBody, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Bad Request: Validate error: invalid projectID",
			setup: func(s *mock_service.MockProjectService) (*schema.ProjectDetail, string) {
				return nil, fmt.Sprintf("/api/v1/projects/%s", invalidID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Internal Error",
			setup: func(s *mock_service.MockProjectService) (*schema.ProjectDetail, string) {
				projectID := random.UUID()
				s.EXPECT().GetProject(anyCtx{}, projectID).Return(nil, errInternal)
				return nil, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "Validation Error",
			setup: func(s *mock_service.MockProjectService) (*schema.ProjectDetail, string) {
				projectID := random.AlphaNumericn(36)
				return nil, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Not Found Error",
			setup: func(s *mock_service.MockProjectService) (*schema.ProjectDetail, string) {
				projectID := random.UUID()
				s.EXPECT().GetProject(anyCtx{}, projectID).Return(nil, repository.ErrNotFound)
				return nil, fmt.Sprintf("/api/v1/projects/%s", projectID)
			},
			statusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupProjectMock(t)

			expectedHres, path := tt.setup(s)

			var hres *schema.ProjectDetail
			statusCode, _ := doRequest(t, api, http.MethodGet, path, nil, &hres)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, expectedHres, hres)
		})
	}
}

func TestProjectHandler_CreateProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(s *mock_service.MockProjectService) (reqBody *schema.CreateProjectJSONRequestBody, expectedResBody schema.Project, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockProjectService) (reqBody *schema.CreateProjectJSONRequestBody, expectedResBody schema.Project, path string) {
				duration := random.Duration()
				reqBody = makeCreateProjectRequest(
					t,
					random.AlphaNumeric(),
					schema.ConvertDuration(duration).Since,
					schema.ConvertDuration(duration).Until,
					random.AlphaNumeric(),
					random.RandURLString(),
				)
				args := repository.CreateProjectArgs{
					Name:          reqBody.Name,
					Description:   reqBody.Description,
					Link:          optional.FromPtr(reqBody.Link),
					SinceYear:     reqBody.Duration.Since.Year,
					SinceSemester: int(reqBody.Duration.Since.Semester),
					UntilYear:     reqBody.Duration.Until.Year,
					UntilSemester: int(reqBody.Duration.Until.Semester),
				}
				want := domain.ProjectDetail{
					Project: domain.Project{
						ID:   random.UUID(),
						Name: args.Name,
						Duration: domain.YearWithSemesterDuration{
							Since: domain.YearWithSemester{
								Year:     args.SinceYear,
								Semester: args.SinceSemester,
							},
							Until: domain.YearWithSemester{
								Year:     args.UntilYear,
								Semester: args.UntilSemester,
							},
						},
					},
					Description: args.Description,
					Link:        args.Link.ValueOrZero(),
					Members:     nil,
				}
				expectedResBody = schema.Project{
					Duration: schema.ConvertDuration(want.Duration),
					Id:       want.ID,
					Name:     want.Name,
				}
				s.EXPECT().CreateProject(anyCtx{}, &args).Return(&want, nil)
				return reqBody, expectedResBody, "/api/v1/projects"
			},
			statusCode: http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupProjectMock(t)

			reqBody, res, path := tt.setup(s)

			var resBody schema.Project
			statusCode, _ := doRequest(t, api, http.MethodPost, path, reqBody, &resBody)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
			assert.Equal(t, resBody, res)
		})
	}
}

func TestProjectHandler_AddProjectMembers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setup      func(s *mock_service.MockProjectService) (reqBody *schema.AddProjectMembersJSONRequestBody, path string)
		statusCode int
	}{
		{
			name: "Success",
			setup: func(s *mock_service.MockProjectService) (*schema.AddProjectMembersJSONRequestBody, string) {
				projectID := random.UUID()
				project := domain.ProjectDetail{
					Project: domain.Project{
						ID:   projectID,
						Name: random.AlphaNumeric(),
						Duration: domain.YearWithSemesterDuration{
							Since: domain.YearWithSemester{
								Year:     2020,
								Semester: 0,
							},
							Until: domain.YearWithSemester{
								Year:     2022,
								Semester: 1,
							},
						},
					},
					Description: random.AlphaNumeric(),
					Link:        random.RandURLString(),
					Members:     []*domain.UserWithDuration{},
				}
				userID := random.UUID()
				userDuration := domain.YearWithSemesterDuration{
					Since: domain.YearWithSemester{
						Year:     2021,
						Semester: 0,
					},
					Until: domain.YearWithSemester{
						Year:     2022,
						Semester: 1,
					},
				}
				reqBody := &schema.AddProjectMembersJSONRequestBody{
					Members: []schema.MemberIDWithYearWithSemesterDuration{
						{
							Duration: schema.YearWithSemesterDuration{
								Since: schema.YearWithSemester{
									Semester: schema.Semester(userDuration.Since.Semester),
									Year:     userDuration.Since.Year,
								},
								Until: &schema.YearWithSemester{
									Semester: schema.Semester(userDuration.Until.Semester),
									Year:     userDuration.Until.Year,
								},
							},
							UserId: userID,
						},
					},
				}
				memberReq := []*repository.CreateProjectMemberArgs{
					{
						UserID:        userID,
						SinceYear:     userDuration.Since.Year,
						SinceSemester: userDuration.Since.Semester,
						UntilYear:     userDuration.Until.Year,
						UntilSemester: userDuration.Until.Semester,
					},
				}
				s.EXPECT().GetProject(anyCtx{}, projectID).Return(&project, nil)
				s.EXPECT().AddProjectMembers(anyCtx{}, projectID, memberReq).Return(nil)
				return reqBody, fmt.Sprintf("/api/v1/projects/%s/members", projectID)
			},
			statusCode: http.StatusNoContent,
		},
		{
			name: "BadRequest: Invalid Project ID",
			setup: func(s *mock_service.MockProjectService) (reqBody *schema.AddProjectMembersJSONRequestBody, path string) {
				projectID := random.UUID()
				return nil, fmt.Sprintf("/api/v1/projects/%s/members", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid request body: member is empty",
			setup: func(s *mock_service.MockProjectService) (reqBody *schema.AddProjectMembersJSONRequestBody, path string) {
				projectID := random.UUID()
				return &schema.AddProjectMembersJSONRequestBody{}, fmt.Sprintf("/api/v1/projects/%s/members", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid request body: memberID is invalid",
			setup: func(s *mock_service.MockProjectService) (reqBody *schema.AddProjectMembersJSONRequestBody, path string) {
				projectID := random.UUID()
				duration := random.Duration()
				return &schema.AddProjectMembersJSONRequestBody{
					Members: []schema.MemberIDWithYearWithSemesterDuration{
						{
							Duration: schema.YearWithSemesterDuration{
								Since: schema.YearWithSemester{
									Semester: schema.Semester(duration.Since.Semester),
									Year:     duration.Since.Year,
								},
								Until: &schema.YearWithSemester{
									Semester: schema.Semester(duration.Until.Semester),
									Year:     duration.Until.Year,
								},
							},
							UserId: uuid.Nil,
						},
					},
				}, fmt.Sprintf("/api/v1/projects/%s/members", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid request body: member is already exists",
			setup: func(s *mock_service.MockProjectService) (*schema.AddProjectMembersJSONRequestBody, string) {
				userID := random.UUID()
				projectID := random.UUID()
				duration := random.Duration()
				reqBody := &schema.AddProjectMembersJSONRequestBody{
					Members: []schema.MemberIDWithYearWithSemesterDuration{
						{
							Duration: schema.YearWithSemesterDuration{
								Since: schema.YearWithSemester{
									Semester: schema.Semester(duration.Since.Semester),
									Year:     duration.Since.Year,
								},
								Until: &schema.YearWithSemester{
									Semester: schema.Semester(duration.Until.Semester),
									Year:     duration.Until.Year,
								},
							},
							UserId: userID,
						},
					},
				}
				s.EXPECT().GetProject(anyCtx{}, projectID).Return(&domain.ProjectDetail{}, nil)
				s.EXPECT().AddProjectMembers(anyCtx{}, projectID, []*repository.CreateProjectMemberArgs{
					{
						UserID:        userID,
						SinceYear:     int(duration.Since.Year),
						SinceSemester: int(duration.Since.Semester),
						UntilYear:     int(duration.Until.Year),
						UntilSemester: int(duration.Until.Semester),
					},
				}).Return(repository.ErrInvalidArg)
				return reqBody, fmt.Sprintf("/api/v1/projects/%s/members", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "BadRequest: invalid request body: bad duration user exists",
			setup: func(s *mock_service.MockProjectService) (*schema.AddProjectMembersJSONRequestBody, string) {
				userID := random.UUID()
				projectID := random.UUID()
				userDuration := domain.YearWithSemesterDuration{
					Since: domain.YearWithSemester{
						Year:     2021,
						Semester: 0,
					},
					Until: domain.YearWithSemester{
						Year:     2023,
						Semester: 1,
					},
				}
				projectDuration := domain.YearWithSemesterDuration{
					Since: domain.YearWithSemester{
						Year:     2022,
						Semester: 0,
					},
					Until: domain.YearWithSemester{
						Year:     2023,
						Semester: 1,
					},
				}
				reqBody := &schema.AddProjectMembersJSONRequestBody{
					Members: []schema.MemberIDWithYearWithSemesterDuration{
						{
							Duration: schema.YearWithSemesterDuration{
								Since: schema.YearWithSemester{
									Semester: schema.Semester(userDuration.Since.Semester),
									Year:     userDuration.Since.Year,
								},
								Until: &schema.YearWithSemester{
									Semester: schema.Semester(userDuration.Until.Semester),
									Year:     userDuration.Until.Year,
								},
							},
							UserId: userID,
						},
					},
				}
				project := domain.ProjectDetail{
					Project: domain.Project{
						ID:       projectID,
						Name:     random.AlphaNumeric(),
						Duration: projectDuration,
					},
					Description: random.AlphaNumeric(),
					Link:        random.RandURLString(),
					Members: []*domain.UserWithDuration{
						{
							User:     *domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
							Duration: random.Duration(),
						},
					},
				}
				s.EXPECT().GetProject(anyCtx{}, projectID).Return(&project, nil)
				return reqBody, fmt.Sprintf("/api/v1/projects/%s/members", projectID)
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			s, api := setupProjectMock(t)

			reqBody, path := tt.setup(s)

			statusCode, _ := doRequest(t, api, http.MethodPost, path, reqBody, nil)

			// Assertion
			assert.Equal(t, tt.statusCode, statusCode)
		})
	}
}
