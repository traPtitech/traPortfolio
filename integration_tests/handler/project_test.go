package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/random"
)

// GetProjects GET /projects
func TestGetProjects(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		want       interface{} // []handler.Project
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockProjects,
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("project_handler_get_projects")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Project.GetProjects), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// GetProject GET /projects/:projectID
func TestGetProject(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		projectID  uuid.UUID
		want       interface{} // handler.ProjectDetail | echo.HTTPError
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockProjectDetails[0].Id,
			mockdata.HMockProjectDetails[0],
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError("bad request: nil id"),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			testutils.HTTPError("not found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("project_handler_get_project")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Project.GetProject, tt.projectID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// CreateProject POST /projects
func TestCreateProjecct(t *testing.T) {
	var (
		name        = random.AlphaNumeric()
		link        = random.RandURLString()
		invalidLink = "invalid link"
		description = random.AlphaNumeric()
		duration    = handler.ConvertDuration(random.Duration())
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		reqBody    handler.CreateProjectJSONRequestBody
		want       interface{} // handler.Project | echo.HTTPError
	}{
		"201": {
			http.StatusCreated,
			handler.CreateProjectJSONRequestBody{
				Name:        name,
				Link:        &link,
				Description: description,
				Duration:    duration,
			},
			handler.Project{
				Id:       uuid.Nil, // OptRetrieveIDで取得する
				Name:     name,
				Duration: duration,
			},
		},
		"400 invalid name": {
			http.StatusBadRequest,
			handler.CreateProjectJSONRequestBody{
				Link:        &link,
				Description: description,
				Duration:    duration,
			},
			testutils.HTTPError("bad request: validate error"),
		},
		"400 invalid URL": {
			http.StatusBadRequest,
			handler.CreateProjectJSONRequestBody{
				Name:        name,
				Link:        &invalidLink,
				Description: description,
				Duration:    duration,
			},
			testutils.HTTPError("bad request: validate error"),
		},
		"400 invalid description": {
			http.StatusBadRequest,
			handler.CreateProjectJSONRequestBody{
				Name:     name,
				Link:     &link,
				Duration: duration,
			},
			testutils.HTTPError("bad request: validate error"),
		},
		"400 invalid duration": {
			http.StatusBadRequest,
			handler.CreateProjectJSONRequestBody{
				Name:        name,
				Link:        &link,
				Description: description,
			},
			testutils.HTTPError("bad request: validate error"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("project_handler_add_project")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodPost, e.URL(api.Project.CreateProject), &tt.reqBody)
			switch want := tt.want.(type) {
			case handler.Project:
				testutils.AssertResponse(t, tt.statusCode, tt.want, res, testutils.OptSyncID, testutils.OptRetrieveID(&want.Id))
			case error:
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

// EditProject PATCH /projects/:projectID
func TestEditProject(t *testing.T) {
	var (
		name        = random.AlphaNumeric()
		link        = random.RandURLString()
		description = random.AlphaNumeric()
		duration    = handler.ConvertDuration(random.Duration())
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		projectID  uuid.UUID
		reqBody    handler.EditProjectJSONRequestBody
		want       interface{} // nil | echo.HTTPError
	}{
		"204": {
			http.StatusNoContent,
			mockdata.HMockProjects[0].Id,
			handler.EditProjectJSONRequestBody{
				Name:        &name,
				Link:        &link,
				Description: &description,
				Duration:    &duration,
			},
			nil,
		},
		"204 without changes": {
			http.StatusNoContent,
			mockdata.HMockProjects[0].Id,
			handler.EditProjectJSONRequestBody{},
			nil,
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.EditProjectJSONRequestBody{},
			testutils.HTTPError("bad request: nil id"),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			handler.EditProjectJSONRequestBody{},
			testutils.HTTPError("not found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("project_handler_update_project")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.statusCode == http.StatusNoContent {
				// Get response before update
				var project handler.ProjectDetail
				res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Project.GetProject, tt.projectID), nil)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &project))
				// Update & Assert
				res = testutils.DoRequest(t, e, http.MethodPatch, e.URL(api.Project.EditProject, tt.projectID), &tt.reqBody)
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
			} else {
				res := testutils.DoRequest(t, e, http.MethodPatch, e.URL(api.Project.EditProject, tt.projectID), &tt.reqBody)
				testutils.AssertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

// GetProjectMembers GET /projects/:projectID/members
func TestGetProjectMembers(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		projectID  uuid.UUID
		want       interface{} // []handler.ProjectMember | echo.HTTPError
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockProjects[0].Id,
			[]handler.ProjectMember{
				mockdata.HMockProjectMembers[0],
				mockdata.HMockProjectMembers[1],
			},
		},
		"200 no members with existing projectID": {
			http.StatusOK,
			mockdata.HMockProjects[2].Id,
			[]handler.ProjectMember{},
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError("bad request: nil id"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("project_handler_get_project_members")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodGet, e.URL(api.Project.GetProjectMembers, tt.projectID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// AddProjectMembers POST /projects/:projectID/members
func TestAddProjectMembers(t *testing.T) {
	var (
		userID1   = mockdata.HMockUsers[0].Id
		duration1 = handler.ConvertDuration(random.Duration())
		userID2   = mockdata.HMockUsers[1].Id
		duration2 = handler.ConvertDuration(random.Duration())
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		projectID  uuid.UUID
		reqBody    handler.AddProjectMembersJSONRequestBody
		want       interface{} // nil | echo.HTTPError
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockProjects[0].Id,
			handler.AddProjectMembersJSONRequestBody{
				Members: []handler.MemberIDWithYearWithSemesterDuration{
					{
						Duration: duration1,
						UserId:   userID1,
					},
					{
						Duration: duration2,
						UserId:   userID2,
					},
				},
			},
			nil,
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.AddProjectMembersJSONRequestBody{},
			testutils.HTTPError("bad request: nil id"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("project_handler_add_member")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodPost, e.URL(api.Project.AddProjectMembers, tt.projectID), &tt.reqBody)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// DeleteProjectMembers DELETE /projects/:projectID/members
func TestDeleteProjectMembers(t *testing.T) {
	var (
		userID1 = mockdata.MockProjectMembers[0].ID
	)
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		projectID  uuid.UUID
		reqBody    handler.DeleteProjectMembersJSONRequestBody
		want       interface{} // nil | echo.HTTPError
	}{
		"204": {
			http.StatusNoContent,
			mockdata.HMockProjects[0].Id,
			handler.DeleteProjectMembersJSONRequestBody{
				Members: []uuid.UUID{userID1},
			},
			nil,
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.DeleteProjectMembersJSONRequestBody{
				Members: []uuid.UUID{userID1},
			},
			testutils.HTTPError("bad request: nil id"),
		},
		"400 invalid memberID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.DeleteProjectMembersJSONRequestBody{
				Members: []uuid.UUID{uuid.Nil},
			},
			testutils.HTTPError("bad request: nil id"),
		},
		"400 invalid members": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.DeleteProjectMembersJSONRequestBody{},
			testutils.HTTPError("bad request: nil id"),
		},
		"404 not found": {
			http.StatusNotFound,
			random.UUID(),
			handler.DeleteProjectMembersJSONRequestBody{
				Members: []uuid.UUID{userID1},
			},
			testutils.HTTPError("not found: not found"),
		},
	}

	e := echo.New()
	conf := testutils.GetConfigWithDBName("project_handler_delete_project")
	api, err := testutils.SetupRoutes(t, e, conf)
	assert.NoError(t, err)
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := testutils.DoRequest(t, e, http.MethodDelete, e.URL(api.Project.DeleteProjectMembers, tt.projectID), &tt.reqBody)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
