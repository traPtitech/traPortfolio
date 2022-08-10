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
		want       interface{} // []handler.Project | echo.HTTPError
	}{
		"200": {
			http.StatusOK,
			[]handler.Project{
				mockdata.HMockProjects[0],
				mockdata.HMockProjects[1],
			},
		},
		"404": {
			http.StatusBadRequest,
			testutils.HTTPError("not found: not found"),
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
			mockdata.HMockProject.Id,
			mockdata.HMockProject,
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError("bad request: invalid project id"),
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
func TestAddProjecct(t *testing.T) {
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
		want       interface{}
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
				Id:       uuid.Nil,
				Name:     name,
				Duration: duration,
			},
		},
		"400 invalid URL": {
			http.StatusBadRequest,
			handler.CreateProjectJSONRequestBody{
				Name:        name,
				Link:        &invalidLink,
				Description: description,
				Duration:    duration,
			},
			testutils.HTTPError("bad request: invalid url"),
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
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
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
		want       interface{} // nil | error
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
			mockdata.HMockProjects[1].Id,
			handler.EditProjectJSONRequestBody{},
			nil,
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.EditProjectJSONRequestBody{},
			testutils.HTTPError("bad request: invalid project id"),
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
		want       interface{}
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockProjects[0].Id,
			mockdata.HMockProjectMembers,
		},
		"200 no members with existing projectID": {
			http.StatusOK,
			mockdata.HMockProjects[1].Id,
			[]handler.User{},
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError("bad request: invalid project id"),
		},
		"404 no project with not-existing projectID": {
			http.StatusNotFound,
			random.UUID(),
			testutils.HTTPError("not found: not found"),
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
		userID1   = random.UUID()
		duration1 = handler.ConvertDuration(random.Duration())
		userID2   = random.UUID()
		duration2 = handler.ConvertDuration(random.Duration())
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		projectID  uuid.UUID
		reqBody    handler.AddProjectMembersJSONRequestBody
		want       interface{}
	}{
		"201": {
			http.StatusCreated,
			mockdata.HMockProjects[0].Id,
			handler.AddProjectMembersJSONRequestBody{
				Members: []handler.MemberIDWithYearWithSemesterDuration{{duration1, userID1}, {duration2, userID2}},
			},
			nil,
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.AddProjectMembersJSONRequestBody{},
			testutils.HTTPError("bad request: invalid project id"),
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
		link        = random.RandURLString()
		description = random.AlphaNumeric()
		duration    = handler.ConvertDuration(random.Duration())
	)
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		projectID  uuid.UUID
		want       interface{}
	}{
		"204": {
			http.StatusNoContent,
			mockdata.HMockProjects[0].Id,
			nil,
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError("bad request: invalid project id"),
		},
		"404 project not found": {
			http.StatusNotFound,
			random.UUID(),
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
			reqBody := handler.CreateProjectJSONRequestBody{
				Name:        random.AlphaNumeric(),
				Link:        &link,
				Description: description,
				Duration:    duration,
			}
			res := testutils.DoRequest(t, e, http.MethodPost, e.URL(api.Project.CreateProject, tt.projectID), &reqBody)
			testutils.AssertResponse(t, http.StatusCreated, handler.Project{
				Duration: reqBody.Duration,
				Id:       tt.projectID,
				Name:     reqBody.Name,
			}, res, testutils.OptSyncID, testutils.OptRetrieveID(&tt.projectID))
			res = testutils.DoRequest(t, e, http.MethodDelete, e.URL(api.Project.DeleteProjectMembers, tt.projectID), nil)
			testutils.AssertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
