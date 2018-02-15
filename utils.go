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

func printJSON(str interface{}) {
	jsonRes, _ := json.MarshalIndent(str, "", "  ")
	fmt.Printf("%+v\n", string(jsonRes))
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

func difference(slice1 []string, slice2 []string) []string {
	var diff []string
	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if s1 == s2 {
				found = true
				break
			}
		}
		// String not found. We add it to return slice
		if !found {
			diff = append(diff, s1)
		}
	}

	return diff
}

func strToIntArr(str string) (intArr []int) {
	for _, s := range strings.Split(str, ",") {
		num, _ := strconv.Atoi(s)
		intArr = append(intArr, num)
	}

	sort.Ints(intArr)
	return intArr
}

func remove(slice []int, s int) []int {
	fmt.Println(slice[:s])
	fmt.Println(slice[s+1:])
	return append(slice[:s], slice[s+1:]...)
}

func deleteSlicefromSlice(slice, delete []int) []int {
	//var out []int
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

func setEmail(c *cli.Context) string {
	var id string
	if c.NArg() == 0 {
		log.Fatal("Please provide user e-mail")
	}

	id = c.Args().Get(0)
	//verifyID(id)
	return id
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
