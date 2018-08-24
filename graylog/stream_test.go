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

func TestAPICanCreateStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			// test incoming request path
			if req.URL.Path != "/streams" && req.URL.Path != "/streams/12345/resume" && req.URL.Path != "/streams/12345/rules" {
				t.Fatalf("unexpected request path %s", req.URL.Path)
				return
			}
			// test incoming params
			body, _ := ioutil.ReadAll(req.Body)
			params := strings.TrimSpace(string(body))
			if params != `{"title": "test-stg", "description": "test-stg", "remove_matches_from_default_stream": "false", "index_set_id": "59bba5962bb9363ee119338c"}` &&
				params != `{"type": "1", "field": "service", "value": "test", "description": ""}` &&
				params != `{"type": "1", "field": "env", "value": "stg", "description": ""}` &&
				params != "" {
				t.Errorf("unexpected params '%v'", params)
				return
			}
			// send result
			result := map[string]string{
				"stream_id": "12345",
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

	var stream Stream
	appName := "test"
	env := "stg"
	stream.CreateStream(&a, &wg, appName, env)

	testID := "12345"
	if stream.ID != testID {
		t.Logf("Expected: %s, got: %s", testID, stream.ID)
		t.Fail()
	}
	t.Log("Call finished")
}
