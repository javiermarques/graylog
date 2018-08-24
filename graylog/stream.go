package graylog

import (
	"encoding/json"
	"fmt"
	"sync"
)

//Stream graylog definition
type Stream struct {
	ID string `json:"stream_id"`
}

//CreateStream creates the stream in graylog, enables the stream and add some default rules to it
func (s *Stream) CreateStream(a *API, wg *sync.WaitGroup, appName string, env string) {
	name := fmt.Sprintf("%s-%s", appName, env)
	str := fmt.Sprintf(`{"title": "%[1]s", "description": "%[1]s", "remove_matches_from_default_stream": "false", "index_set_id": "59bba5962bb9363ee119338c"}`, name)
	var jsonStr = []byte(str)
	resp := a.sendRequest("/streams", jsonStr, "POST")
	json.Unmarshal(resp, &s)

	s.enableStream(a)
	s.addRules(a, appName, env)

	defer wg.Done()
}

func (s *Stream) enableStream(a *API) {
	path := fmt.Sprintf("/streams/%s/resume", s.ID)
	_ = a.sendRequest(path, nil, "POST")
}

func (s *Stream) addRules(a *API, appName string, env string) {
	path := fmt.Sprintf("/streams/%s/rules", s.ID)

	str := fmt.Sprintf(`{"type": "1", "field": "service", "value": "%s", "description": ""}`, appName)
	var jsonStr = []byte(str)
	_ = a.sendRequest(path, jsonStr, "POST")

	str = fmt.Sprintf(`{"type": "1", "field": "env", "value": "%s", "description": ""}`, env)
	jsonStr = []byte(str)
	_ = a.sendRequest(path, jsonStr, "POST")
}
