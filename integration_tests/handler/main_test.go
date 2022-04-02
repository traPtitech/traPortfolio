//go:build integration && db

package handler_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/util/config"
)

var (
	port    = config.Port()
	baseURL = fmt.Sprintf("http://localhost%s/api/v1", port)
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

		// TODO: DBのセットアップを行う
		// Ref: https://github.com/traPtitech/traPortfolio/pull/228

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
