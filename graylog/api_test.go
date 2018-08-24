package graylog

import (
	"testing"
)

func TestAPICreation(t *testing.T) {
	url := "MOCK_URL"
	token := "MOCK_TOKEN"
	a := NewAPI(url, token)
	if a.token != token {
		t.Logf("API Creation failed, token not set. Expected %s got %s", token, a.token)
		t.Fail()
	}
	if a.baseURL != url {
		t.Logf("API Creation failed, url not set. Expected %s got %s", url, a.baseURL)
		t.Fail()
	}
	t.Log("API created. Token set.")
}
