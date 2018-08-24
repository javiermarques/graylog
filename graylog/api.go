package graylog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var debug bool

//API definition
//Graylog token needed for auth
type API struct {
	baseURL string
	token   string
}

//NewAPI API object constructor
func NewAPI(url string, t string) API {
	a := new(API)
	a.baseURL = url
	a.token = t
	return *a
}

func (a *API) sendRequest(endpoint string, body []byte, t string) []byte {
	url := a.baseURL + endpoint
	req, err := http.NewRequest(t, url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return nil
	}

	req.SetBasicAuth(a.token, "token")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	if debug {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		fmt.Println("response Body:", string(respBody))
	}

	return respBody
}
