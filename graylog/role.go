package graylog

import (
	"encoding/json"
	"fmt"
	"sync"
)

//Role graylog
type Role struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Description string   `json:"description"`
	ReadOnly    bool     `json:"read_only"`
}

var teamRole string

//GetPermissions - syncs the struct with the permissions from graylog
func (r *Role) GetPermissions(a *API, wg *sync.WaitGroup, team string) {
	teamRole = team
	endpoint := fmt.Sprintf("/roles/%s", teamRole)
	resp := a.sendRequest(endpoint, nil, "GET")
	json.Unmarshal(resp, &r)
	defer wg.Done()
}

//UpdateRole - adds the new dashboard and stream to the team
func (r Role) UpdateRole(a *API, s Stream, d Dashboard) {
	r.Permissions = append(r.Permissions, fmt.Sprintf("streams:edit:%s", s.ID), fmt.Sprintf("streams:read:%s", s.ID), fmt.Sprintf("dashboards:edit:%s", d.ID), fmt.Sprintf("dashboards:read:%s", d.ID))

	jsonStr, _ := json.Marshal(&r)
	if debug {
		fmt.Printf("JSON new: %s", jsonStr)
	}

	endpoint := fmt.Sprintf("/roles/%s", teamRole)
	_ = a.sendRequest(endpoint, jsonStr, "PUT")
}

//CleanRole - Aux func used to clean the role permissions, mainly for test usage
func (r Role) CleanRole(a API) {
	var jsonStr = []byte(`{"name":"testkk","description":"testkk","permissions":[],"read_only":false}`)
	_ = a.sendRequest("/roles/testkk", jsonStr, "PUT")
}
