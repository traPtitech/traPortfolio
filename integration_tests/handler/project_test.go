package handler

import (
	"encoding/json"
	"net/http"
	"strings"
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
			mockdata.ProjectID1(),
			mockdata.HMockProjectDetails[0],
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError("Bad Request: nil id"),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			testutils.HTTPError("Not Found: not found"),
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
func TestCreateProject(t *testing.T) {
	var (
		name                    = random.AlphaNumeric()
		link                    = random.RandURLString()
		invalidLink             = "invalid link"
		description             = random.AlphaNumeric()
		justCountName           = strings.Repeat("亜", 32)
		justCountDescription    = strings.Repeat("亜", 256)
		tooLongName             = strings.Repeat("亜", 33)
		tooLongDescriptionKanji = strings.Repeat("亜", 257)
		duration                = handler.ConvertDuration(random.Duration())
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
		"201 with kanji": {
			http.StatusCreated,
			handler.CreateProjectJSONRequestBody{
				Name:        justCountName,
				Link:        &link,
				Description: justCountDescription,
				Duration:    duration,
			},
			handler.Project{
				Id:       uuid.Nil,
				Name:     justCountName,
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
			testutils.HTTPError("Bad Request: validate error: link: must be a valid URL."),
		},
		"400 too long description": {
			http.StatusBadRequest,
			handler.CreateProjectJSONRequestBody{
				Name:        name,
				Link:        &link,
				Description: tooLongDescriptionKanji,
				Duration:    duration,
			},
			testutils.HTTPError("Bad Request: validate error: description: must be less than or equal to 256 characters."),
		},
		"400 too long name": {
			http.StatusBadRequest,
			handler.CreateProjectJSONRequestBody{
				Name:        tooLongName,
				Link:        &link,
				Description: description,
				Duration:    duration,
			},
			testutils.HTTPError("Bad Request: validate error: name: must be less than or equal to 32 characters."),
		},
		"400 empty name": {
			http.StatusBadRequest,
			handler.CreateProjectJSONRequestBody{
				Link:        &link,
				Description: description,
				Duration:    duration,
			},
			testutils.HTTPError("Bad Request: validate error: name: cannot be blank."),
		},
		"400 empty description": {
			http.StatusBadRequest,
			handler.CreateProjectJSONRequestBody{
				Name:     name,
				Link:     &link,
				Duration: duration,
			},
			testutils.HTTPError("Bad Request: validate error: description: cannot be blank."),
		},
		"400 empty duration": {
			http.StatusBadRequest,
			handler.CreateProjectJSONRequestBody{
				Name:        name,
				Link:        &link,
				Description: description,
			},
			testutils.HTTPError("Bad Request: argument error"),
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
		justCountName           = strings.Repeat("亜", 32)
		justCountDescription    = strings.Repeat("亜", 256)
		tooLongName             = strings.Repeat("亜", 33)
		tooLongDescriptionKanji = strings.Repeat("亜", 257)
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
			mockdata.ProjectID1(),
			handler.EditProjectJSONRequestBody{
				Name:        &name,
				Link:        &link,
				Description: &description,
				Duration:    &duration,
			},
			nil,
		},
		"204 with kanji": {
			http.StatusNoContent,
			mockdata.ProjectID1(),
			handler.EditProjectJSONRequestBody{
				Name:        &justCountName,
				Link:        &link,
				Description: &justCountDescription,
				Duration:    &duration,
			},
			nil,
		},
		"204 without changes": {
			http.StatusNoContent,
			mockdata.ProjectID2(),
			handler.EditProjectJSONRequestBody{},
			nil,
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			handler.EditProjectJSONRequestBody{},
			testutils.HTTPError("Bad Request: nil id"),
		},
		"400 invalid Name": {
			http.StatusBadRequest,
			mockdata.ProjectID1(),
			handler.EditProjectJSONRequestBody{
				Name: &tooLongName,
			},
			testutils.HTTPError("Bad Request: validate error: name: must be less than or equal to 32 characters."),
		},
		"400 invalid Description": {
			http.StatusBadRequest,
			mockdata.ProjectID1(),
			handler.EditProjectJSONRequestBody{
				Description: &tooLongDescriptionKanji,
			},
			testutils.HTTPError("Bad Request: validate error: description: must be less than or equal to 256 characters."),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			handler.EditProjectJSONRequestBody{},
			testutils.HTTPError("Not Found: not found"),
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
			mockdata.ProjectID1(),
			[]handler.ProjectMember{
				mockdata.HMockProjectMembers[0],
				mockdata.HMockProjectMembers[1],
			},
		},
		"200 no members with existing projectID": {
			http.StatusOK,
			mockdata.ProjectID3(),
			[]handler.ProjectMember{},
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			testutils.HTTPError("Bad Request: nil id"),
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
		userID1   = mockdata.UserID1()
		duration1 = handler.ConvertDuration(random.Duration())
		userID2   = mockdata.UserID2()
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
			mockdata.ProjectID1(),
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
			testutils.HTTPError("Bad Request: nil id"),
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
			mockdata.ProjectID1(),
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
			testutils.HTTPError("Bad Request: nil id"),
		},
		"400 invalid memberID": {
			http.StatusBadRequest,
			random.UUID(),
			handler.DeleteProjectMembersJSONRequestBody{
				Members: []uuid.UUID{uuid.Nil},
			},
			testutils.HTTPError("Bad Request: validate error: members: (0: must be a valid UUID v4.)."),
		},
		"400 empty members": {
			http.StatusBadRequest,
			random.UUID(),
			handler.DeleteProjectMembersJSONRequestBody{},
			testutils.HTTPError("Bad Request: validate error: members: cannot be blank."),
		},
		"400 empty memberIDs": {
			http.StatusBadRequest,
			random.UUID(),
			handler.DeleteProjectMembersJSONRequestBody{
				Members: []uuid.UUID{},
			},
			testutils.HTTPError("Bad Request: validate error: members: cannot be blank."),
		},
		"404 not found": {
			http.StatusNotFound,
			random.UUID(),
			handler.DeleteProjectMembersJSONRequestBody{
				Members: []uuid.UUID{userID1},
			},
			testutils.HTTPError("Not Found: not found"),
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
