package graylog

import (
	"encoding/json"
	"fmt"
	"sync"
)

//Dashboard Graylog representation
type Dashboard struct {
	ID string `json:"dashboard_id"`
}

//CreateDashboard - Async create a dashboard on graylog
func (d *Dashboard) CreateDashboard(a *API, wg *sync.WaitGroup, appName string, env string) {
	name := fmt.Sprintf("%s-%s", appName, env)
	str := fmt.Sprintf(`{"title": "%[1]s", "description": "%[1]s"}`, name)
	var jsonStr = []byte(str)
	resp := a.sendRequest("/dashboards", jsonStr, "POST")
	json.Unmarshal(resp, &d)
	if debug {
		fmt.Println(d.ID)
	}
	defer wg.Done()
}
