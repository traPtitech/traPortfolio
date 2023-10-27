package external

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/traPtitech/traPortfolio/util/config"
)

type apiClient struct {
	client *http.Client
	conf   *config.APIConfig
}

func newAPIClient(jar *cookiejar.Jar, conf *config.APIConfig) apiClient {
	return apiClient{
		client: &http.Client{Jar: jar},
		conf:   conf,
	}
}

func (c *apiClient) apiGet(path string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", c.conf.APIEndpoint, path), nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func newCookieJar(c *config.APIConfig, name string) (*cookiejar.Jar, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	cookies := []*http.Cookie{
		{
			Name:  name,
			Value: c.Cookie,
			Path:  "/",
		},
	}

	u, err := url.Parse(c.APIEndpoint)
	if err != nil {
		return nil, err
	}

	jar.SetCookies(u, cookies)

	return jar, nil
}
