package graylog

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestRoleGetPermissions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			// test incoming request path
			if req.URL.Path != "/roles/team" {
				t.Fatalf("unexpected request path %s", req.URL.Path)
				return
			}
			// test incoming params
			body, _ := ioutil.ReadAll(req.Body)
			params := strings.TrimSpace(string(body))
			if params != "" {
				t.Errorf("unexpected params '%v'", params)
				return
			}
			// send result
			_, err := resp.Write([]byte(`{"name": "team", "description": "team", "permissions": ["dashboards:read:12345", "streams:read:12345", "streams:edit:12345"], "read_only": "false"}`))
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

	var role Role
	teamName := "team"
	role.GetPermissions(&a, &wg, teamName)

	if role.Name != teamName {
		t.Logf("Expected: %s, got: %s", teamName, role.Name)
		t.Fail()
	}
	t.Log("Call finished")
}

func TestCreateRole(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			// test incoming request path
			if req.URL.Path != "/roles/team" {
				t.Fatalf("unexpected request path %s", req.URL.Path)
				return
			}
			// test incoming params
			body, _ := ioutil.ReadAll(req.Body)
			params := strings.TrimSpace(string(body))
			if params != `{"name":"team","permissions":["streams:edit:12345","streams:read:12345","dashboards:edit:12345","dashboards:read:12345"],"description":"team","read_only":false}` &&
				params != "" {
				t.Errorf("unexpected params '%v'", params)
				return
			}
			// send result
			_, err := resp.Write([]byte(`{"name": "team", "description": "team", "permissions": [], "read_only": "false"}`))
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

	var role Role
	teamName := "team"
	role.GetPermissions(&a, &wg, teamName)
	d := Dashboard{"12345"}
	s := Stream{"12345"}
	role.UpdateRole(&a, s, d)

	t.Log("Call finished")
}

func Contains(sl []string, v string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
