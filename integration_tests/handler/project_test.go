package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/internal/handler/schema"
	"github.com/traPtitech/traPortfolio/internal/pkgs/mockdata"
	"github.com/traPtitech/traPortfolio/internal/pkgs/random"
)

// GetProjects GET /projects
func TestGetProjects(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		want       interface{} // []schema.Project
	}{
		"200": {
			http.StatusOK,
			mockdata.HMockProjects,
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.Project.GetProjects), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// GetProject GET /projects/:projectID
func TestGetProject(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		statusCode int
		projectID  uuid.UUID
		want       interface{} // schema.ProjectDetail | echo.HTTPError
	}{
		"200": {
			http.StatusOK,
			mockdata.ProjectID1(),
			mockdata.HMockProjectDetails[0],
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			httpError(t, "Bad Request: nil id"),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.Project.GetProject, tt.projectID), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// CreateProject POST /projects
func TestCreateProject(t *testing.T) {
	var (
		name                    = random.AlphaNumeric()
		links                   = random.Array(random.RandURLString, 1, 3)
		invalidLink             = []string{"invalid link"}
		description             = random.AlphaNumeric()
		justCountName           = strings.Repeat("亜", 32)
		justCountDescription    = strings.Repeat("亜", 256)
		tooLongName             = strings.Repeat("亜", 33)
		tooLongDescriptionKanji = strings.Repeat("亜", 257)
		duration                = schema.ConvertDuration(random.Duration())
		conflictedProject       = random.CreateProjectArgs()
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		reqBody    schema.CreateProjectRequest
		want       interface{} // schema.Project | echo.HTTPError
	}{
		"201": {
			http.StatusCreated,
			schema.CreateProjectRequest{
				Name:        name,
				Links:       links,
				Description: description,
				Duration:    duration,
			},
			schema.Project{
				Id:       uuid.Nil, // OptRetrieveIDで取得する
				Name:     name,
				Duration: duration,
			},
		},
		"201 with kanji": {
			http.StatusCreated,
			schema.CreateProjectRequest{
				Name:        justCountName,
				Links:       links,
				Description: justCountDescription,
				Duration:    duration,
			},
			schema.Project{
				Id:       uuid.Nil,
				Name:     justCountName,
				Duration: duration,
			},
		},
		"400 invalid URL": {
			http.StatusBadRequest,
			schema.CreateProjectRequest{
				Name:        name,
				Links:       invalidLink,
				Description: description,
				Duration:    duration,
			},
			httpError(t, "Bad Request: validate error: links: (0: must be a valid URL.)."),
		},
		"400 too long description": {
			http.StatusBadRequest,
			schema.CreateProjectRequest{
				Name:        name,
				Links:       links,
				Description: tooLongDescriptionKanji,
				Duration:    duration,
			},
			httpError(t, "Bad Request: validate error: description: the length must be between 1 and 256."),
		},
		"400 too long name": {
			http.StatusBadRequest,
			schema.CreateProjectRequest{
				Name:        tooLongName,
				Links:       links,
				Description: description,
				Duration:    duration,
			},
			httpError(t, "Bad Request: validate error: name: the length must be between 1 and 32."),
		},
		"400 empty name": {
			http.StatusBadRequest,
			schema.CreateProjectRequest{
				Links:       links,
				Description: description,
				Duration:    duration,
			},
			httpError(t, "Bad Request: validate error: name: cannot be blank."),
		},
		"400 empty description": {
			http.StatusBadRequest,
			schema.CreateProjectRequest{
				Name:     name,
				Links:    links,
				Duration: duration,
			},
			httpError(t, "Bad Request: validate error: description: cannot be blank."),
		},
		"400 empty duration": {
			http.StatusBadRequest,
			schema.CreateProjectRequest{
				Name:        name,
				Links:       links,
				Description: description,
			},
			httpError(t, "Bad Request: argument error"),
		},
		"400 project already exists": {
			http.StatusBadRequest,
			schema.CreateProjectRequest{
				Name:        conflictedProject.Name,
				Links:       links,
				Description: description,
			},
			httpError(t, "Bad Request: argument error"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_ = doRequest(t, e, http.MethodPost, e.URL(api.Project.CreateProject), &conflictedProject)
			res := doRequest(t, e, http.MethodPost, e.URL(api.Project.CreateProject), &tt.reqBody)
			switch want := tt.want.(type) {
			case schema.Project:
				assertResponse(t, tt.statusCode, tt.want, res, optSyncID, optRetrieveID(&want.Id))
			case error:
				assertResponse(t, tt.statusCode, tt.want, res)
			}
		})
	}
}

// EditProject PATCH /projects/:projectID
func TestEditProject(t *testing.T) {
	var (
		name                    = random.AlphaNumeric()
		links                   = random.Array(random.RandURLString, 1, 3)
		description             = random.AlphaNumeric()
		justCountName           = strings.Repeat("亜", 32)
		justCountDescription    = strings.Repeat("亜", 256)
		tooLongName             = strings.Repeat("亜", 33)
		tooLongDescriptionKanji = strings.Repeat("亜", 257)
		duration                = schema.ConvertDuration(random.Duration())
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		projectID  uuid.UUID
		reqBody    schema.EditProjectRequest
		want       interface{} // nil | echo.HTTPError
	}{
		"204": {
			http.StatusNoContent,
			mockdata.ProjectID1(),
			schema.EditProjectRequest{
				Name:        &name,
				Links:       &links,
				Description: &description,
				Duration:    &duration,
			},
			nil,
		},
		"204 with kanji": {
			http.StatusNoContent,
			mockdata.ProjectID1(),
			schema.EditProjectRequest{
				Name:        &justCountName,
				Links:       &links,
				Description: &justCountDescription,
				Duration:    &duration,
			},
			nil,
		},
		"204 without changes": {
			http.StatusNoContent,
			mockdata.ProjectID2(),
			schema.EditProjectRequest{},
			nil,
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			schema.EditProjectRequest{},
			httpError(t, "Bad Request: nil id"),
		},
		"400 invalid Name": {
			http.StatusBadRequest,
			mockdata.ProjectID1(),
			schema.EditProjectRequest{
				Name: &tooLongName,
			},
			httpError(t, "Bad Request: validate error: name: the length must be between 1 and 32."),
		},
		"400 invalid Description": {
			http.StatusBadRequest,
			mockdata.ProjectID1(),
			schema.EditProjectRequest{
				Description: &tooLongDescriptionKanji,
			},
			httpError(t, "Bad Request: validate error: description: the length must be between 1 and 256."),
		},
		"404": {
			http.StatusNotFound,
			random.UUID(),
			schema.EditProjectRequest{},
			httpError(t, "Not Found: not found"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.statusCode == http.StatusNoContent {
				// Get response before update
				var project schema.ProjectDetail
				res := doRequest(t, e, http.MethodGet, e.URL(api.Project.GetProject, tt.projectID), nil)
				assert.Equal(t, http.StatusOK, res.Code)
				assert.NoError(t, json.Unmarshal(res.Body.Bytes(), &project))
				// Update & Assert
				res = doRequest(t, e, http.MethodPatch, e.URL(api.Project.EditProject, tt.projectID), &tt.reqBody)
				assertResponse(t, tt.statusCode, tt.want, res)
			} else {
				res := doRequest(t, e, http.MethodPatch, e.URL(api.Project.EditProject, tt.projectID), &tt.reqBody)
				assertResponse(t, tt.statusCode, tt.want, res)
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
		want       interface{} // []schema.ProjectMember | echo.HTTPError
	}{
		"200": {
			http.StatusOK,
			mockdata.ProjectID1(),
			[]schema.ProjectMember{
				mockdata.HMockProjectMembers[0],
				mockdata.HMockProjectMembers[1],
			},
		},
		"200 no members with existing projectID": {
			http.StatusOK,
			mockdata.ProjectID3(),
			[]schema.ProjectMember{},
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			httpError(t, "Bad Request: nil id"),
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodGet, e.URL(api.Project.GetProjectMembers, tt.projectID), nil)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}

// EditProjectMembers PUT /projects/:projectID/members
func TestEditProjectMembers(t *testing.T) {
	var (
		userID1 = mockdata.UserID1()
		userID2 = mockdata.UserID2()
	)

	t.Parallel()
	tests := map[string]struct {
		statusCode int
		projectID  uuid.UUID
		reqBody    schema.EditProjectMembersRequest
		want       interface{} // nil | echo.HTTPError
	}{
		"204": {
			http.StatusNoContent,
			mockdata.ProjectID3(),
			schema.EditProjectMembersRequest{
				Members: []schema.MemberIDWithYearWithSemesterDuration{
					{
						Duration: schema.YearWithSemesterDuration{
							Since: schema.YearWithSemester{
								Year:     2021,
								Semester: 0,
							},
							Until: &schema.YearWithSemester{
								Year:     2021,
								Semester: 1,
							},
						},
						UserId: userID1,
					},
					{
						Duration: schema.YearWithSemesterDuration{
							Since: schema.YearWithSemester{
								Year:     2021,
								Semester: 1,
							},
							Until: &schema.YearWithSemester{
								Year:     2022,
								Semester: 1,
							},
						},
						UserId: userID2,
					},
				},
			},
			nil,
		},
		"400 invalid projectID": {
			http.StatusBadRequest,
			uuid.Nil,
			schema.EditProjectMembersRequest{},
			httpError(t, "Bad Request: nil id"),
		},
		"400 exceeded duration user exists": {
			http.StatusBadRequest,
			mockdata.ProjectID1(),
			schema.EditProjectMembersRequest{
				Members: []schema.MemberIDWithYearWithSemesterDuration{
					{
						Duration: schema.YearWithSemesterDuration{
							Since: schema.YearWithSemester{
								Year:     2021,
								Semester: 0,
							},
							Until: &schema.YearWithSemester{
								Year:     2024,
								Semester: 1,
							},
						},
						UserId: userID1,
					},
				},
			},
			httpError(t, "Bad Request: argument error: exceeded duration user(project: {Since:{Year:2021 Semester:0} Until:{v:{Year:2021 Semester:1} valid:true}}, member: {Since:{Year:2021 Semester:0} Until:{v:{Year:2024 Semester:1} valid:true}})"), // TODO: improve message
		},
	}

	e := echo.New()
	api := setupRoutes(t, e)
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := doRequest(t, e, http.MethodPut, e.URL(api.Project.EditProjectMembers, tt.projectID), &tt.reqBody)
			assertResponse(t, tt.statusCode, tt.want, res)
		})
	}
}
