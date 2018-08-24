package graylog

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestAPICanCreateDashboard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			// test incoming request path
			if req.URL.Path != "/dashboards" {
				t.Errorf("unexpected request path %s", req.URL.Path)
				return
			}
			// test incoming params
			body, _ := ioutil.ReadAll(req.Body)
			params := strings.TrimSpace(string(body))
			if params != `{"title": "test-stg", "description": "test-stg"}` {
				t.Errorf("unexpected params '%v'", params)
				return
			}
			// send result
			result := map[string]string{
				"dashboard_id": "12345",
			}
			err := json.NewEncoder(resp).Encode(&result)
			if err != nil {
				t.Fatal(err)
				return
			}
		},
	))
	defer server.Close()

	token := "MOCK_TOKEN"
	var wg sync.WaitGroup
	wg.Add(1)
	a := NewAPI(server.URL, token)

	var dash Dashboard
	appName := "test"
	env := "stg"
	dash.CreateDashboard(&a, &wg, appName, env)

	testID := "12345"
	if dash.ID != testID {
		t.Logf("Expected: %s, got: %s", testID, dash.ID)
		t.Fail()
	}
	t.Log("Call finished")
}
