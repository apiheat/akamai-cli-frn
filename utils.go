package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"

	client "github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/urfave/cli"
)

func errorCheck(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func printJSON(str interface{}) {
	jsonRes, _ := json.MarshalIndent(str, "", "  ")
	fmt.Printf("%+v\n", string(jsonRes))
}

func setEmail(c *cli.Context) string {
	var id string
	if c.NArg() == 0 {
		log.Fatal("Please provide user e-mail")
	}

	id = c.Args().Get(0)
	//verifyID(id)
	return id
}

func setID(c *cli.Context) string {
	var id string
	if c.NArg() == 0 {
		log.Fatal("Please provide ID for map")
	}

	id = c.Args().Get(0)
	verifyID(id)
	return id
}

func verifyID(id string) {
	if _, err := strconv.Atoi(id); err != nil {
		errStr := fmt.Sprintf("SiteShield Map ID should be number, you provided: %q\n", id)
		log.Fatal(errStr)
	}
}

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
	errorCheck(err)

	resp, err := client.Do(edgeConfig, req)
	errorCheck(err)

	defer resp.Body.Close()
	byt, _ := ioutil.ReadAll(resp.Body)

	return string(byt)
}

func strToIntArr(str string) (intArr []int) {
	for _, s := range strings.Split(str, ",") {
		num, _ := strconv.Atoi(s)
		intArr = append(intArr, num)
	}

	sort.Ints(intArr)
	return intArr
}

func stringInSlice(a string, list []string) bool {
	// We need that to not filter for empty list
	if len(list) > 0 {
		for _, b := range list {
			if b == a {
				return true
			}
		}
		return false
	}
	return true
}

func deleteSlicefromSlice(slice, delete []int) []int {
	for _, d := range delete {
		for i := len(slice) - 1; i >= 0; i-- {
			if slice[i] == d {
				slice = append(slice[:i], slice[i+1:]...)
			}
		}
	}

	return removeDuplicates(slice)
}

func removeDuplicates(elements []int) []int {
	// Use map to record duplicates as we find them.
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
