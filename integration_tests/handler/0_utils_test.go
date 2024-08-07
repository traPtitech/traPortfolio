package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/internal/handler"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/external/mock_external_e2e"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/migration"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/repository"
	"github.com/traPtitech/traPortfolio/internal/pkgs/config"
	"github.com/traPtitech/traPortfolio/internal/pkgs/mockdata"
	"github.com/traPtitech/traPortfolio/internal/pkgs/random"
	"gorm.io/gorm"
)

func injectIntoAPIServer(t *testing.T, c *config.Config, db *gorm.DB) (handler.API, error) {
	t.Helper()

	// FIXME: モック前提のテストがあるためassert
	assert.False(t, c.IsProduction)

	// external API
	portalAPI := mock_external_e2e.NewMockPortalAPI()
	traQAPI := mock_external_e2e.NewMockTraQAPI()
	knoqAPI := mock_external_e2e.NewMockKnoqAPI()

	// repository
	userRepo := repository.NewUserRepository(db, portalAPI, traQAPI)
	projectRepo := repository.NewProjectRepository(db, portalAPI)
	eventRepo := repository.NewEventRepository(db, knoqAPI)
	contestRepo := repository.NewContestRepository(db, portalAPI)
	groupRepo := repository.NewGroupRepository(db)

	// service, handler, API
	api := handler.NewAPI(
		handler.NewPingHandler(),
		handler.NewUserHandler(userRepo, eventRepo),
		handler.NewProjectHandler(projectRepo),
		handler.NewEventHandler(eventRepo, userRepo),
		handler.NewContestHandler(contestRepo),
		handler.NewGroupHandler(groupRepo, userRepo),
	)

	return api, nil
}

func setupRoutes(t *testing.T, e *echo.Echo) *handler.API {
	t.Helper()

	db := SetupGormDB(t)
	err := mockdata.InsertSampleDataToDB(db)
	assert.NoError(t, err)

	api, err := injectIntoAPIServer(t, testConfig, db)
	assert.NoError(t, err)

	err = handler.Setup(false, e, api)
	assert.NoError(t, err)

	return &api
}

func SetupGormDB(t *testing.T) *gorm.DB {
	t.Helper()

	db := establishTestDBConnection(t)
	dropAll(t, db)
	init, err := migration.Migrate(db, migration.AllTables())
	assert.True(t, init)
	assert.NoError(t, err)

	return db
}

func establishTestDBConnection(t *testing.T) *gorm.DB {
	t.Helper()

	sqlConf := testConfig.DB
	sqlConf.Name = "portfolio_test_" + t.Name()

	_, err := testDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", sqlConf.Name))
	assert.NoError(t, err)

	db, err := repository.NewGormDB(sqlConf)
	assert.NoError(t, err)

	return db
}

func dropAll(t *testing.T, db *gorm.DB) {
	t.Helper()

	tables := []interface{}{"migrations"}
	tables = append(tables, migration.AllTables()...)

	err := db.Migrator().DropTable(tables...)
	assert.NoError(t, err)
}

func doRequest(t *testing.T, e *echo.Echo, method string, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		assert.NoError(t, err)

		bodyReader = bytes.NewReader(b)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	return rec
}

// OptRetrieveIDなどですぐに変更され得るUUIDであることの明示に使う
func dummyUUID(t *testing.T) uuid.UUID {
	t.Helper()
	return random.UUID()
}

func httpError(t *testing.T, message string) *echo.HTTPError {
	t.Helper()
	return &echo.HTTPError{
		Message: message,
	}
}

func assertResponse(t *testing.T, expectedStatusCode int, expectedBody interface{}, res *httptest.ResponseRecorder, opts ...option) {
	t.Helper()

	assert.Equal(t, expectedStatusCode, res.Code)

	for _, o := range opts {
		assert.NoError(t, o(t, &expectedBody, res))
	}

	actual := res.Body
	switch expected := expectedBody.(type) {
	case string:
		assert.Equal(t, expected, actual.String())
	case []byte:
		assert.Equal(t, expected, actual.Bytes())
	case nil:
		assert.Empty(t, actual.String())
	default:
		b, err := json.Marshal(expected)
		assert.NoError(t, err)
		assert.JSONEq(t, string(b), actual.String())
	}
}

// NOTE: expectedBodyPtr must be a pointer to expectedBody
type option func(t *testing.T, expectedBodyPtr interface{}, res *httptest.ResponseRecorder) error

// TODO: testifyからgo-cmpに乗り換えたらFilterValuesを使ってIDをignoreする
func optSyncID(t *testing.T, expectedBodyPtr interface{}, res *httptest.ResponseRecorder) error {
	t.Helper()

	m := struct {
		ID uuid.UUID `json:"id"`
	}{}

	if err := json.Unmarshal(res.Body.Bytes(), &m); err != nil {
		return err
	}

	v := reflect.ValueOf(expectedBodyPtr).Elem()
	tmp := reflect.New(v.Elem().Type()).Elem()
	tmp.Set(v.Elem())
	tmp.FieldByName("Id").Set(reflect.ValueOf(m.ID))
	v.Set(tmp)

	return nil
}

func optRetrieveID(idPtr *uuid.UUID) option {
	return func(t *testing.T, expectedBodyPtr interface{}, res *httptest.ResponseRecorder) error {
		t.Helper()

		m := struct {
			ID uuid.UUID `json:"id"`
		}{}

		if err := json.Unmarshal(res.Body.Bytes(), &m); err != nil {
			return err
		}

		*idPtr = m.ID

		return nil
	}
}
