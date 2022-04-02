//go:build integration && db

package handler_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
)

var (
	port = func() int {
		if p, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
			return p
		}

		return 1323
	}()
	baseURL = fmt.Sprintf("http://localhost:%d/api/v1", port)
)

func mustMarshal(t *testing.T, v interface{}) []byte {
	t.Helper()

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func TestMain(m *testing.M) {
	s := infrastructure.NewSQLConfig("root", "password", "localhost", "portfolio", 3307)
	t := infrastructure.NewTraQConfig("", "", true)
	p := infrastructure.NewPortalConfig("", "", true)
	k := infrastructure.NewKnoqConfig("", "", true)
	api, err := infrastructure.InjectAPIServer(&s, &t, &p, &k)
	if err != nil {
		log.Fatal(err)
	}

	go func(api handler.API) {
		e := echo.New()
		handler.Setup(e, api)

		log.Fatal(e.Start(fmt.Sprintf(":%d", port)))
	}(api)

	for i := 0; ; i++ {
		log.Printf("waiting for server to start... (%ds)", i)

		if i > 10 {
			log.Fatal("failed to connect to server")
		}

		if _, err := http.Get(baseURL + "/ping"); err == nil {
			break
		}

		time.Sleep(time.Second)
	}

	m.Run()
}
