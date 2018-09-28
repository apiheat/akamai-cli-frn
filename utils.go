package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	common "github.com/apiheat/akamai-cli-common"
)

func servicesParse(in string) (services Services, err error) {
	if err = json.Unmarshal([]byte(in), &services); err != nil {
		return
	}
	return
}

func cidrsParse(in string) (services Cidrs, err error) {
	if err = json.Unmarshal([]byte(in), &services); err != nil {
		return
	}
	return
}

func subscriptionsParse(in string) (subscriptions SubscriptionsResp, err error) {
	if err = json.Unmarshal([]byte(in), &subscriptions); err != nil {
		return
	}
	return
}

func createSubscriptionBody(services []int, email string) io.Reader {
	var obj SubscriptionsResp
	for _, s := range services {
		service := Subscription{ServiceID: s, Email: email}
		obj.Subscriptions = append(obj.Subscriptions, service)
	}

	json, _ := json.Marshal(obj)

	return strings.NewReader(string(json))
}

func fetchData(urlPath, method string, body io.Reader) (result string) {
	req, err := client.NewRequest(edgeConfig, method, urlPath, body)
	common.ErrorCheck(err)

	resp, err := client.Do(edgeConfig, req)
	common.ErrorCheck(err)

	defer resp.Body.Close()
	byt, _ := ioutil.ReadAll(resp.Body)

	return string(byt)
}
