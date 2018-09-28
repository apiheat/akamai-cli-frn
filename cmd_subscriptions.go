package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	common "github.com/apiheat/akamai-cli-common"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdSubscriptions(c *cli.Context) error {
	return listSubscriptions(c)
}

func cmdUpdSubscriptions(c *cli.Context) error {
	return updateSubscriptions(c)
}

func listSubscriptions(c *cli.Context) error {
	data := listSubscriptionsData(c)

	if raw {
		println(data)

		return nil
	}

	result, err := subscriptionsParse(data)
	common.ErrorCheck(err)

	printSubscriptions(result)
	return nil
}

func listSubscriptionsData(c *cli.Context) string {
	urlStr := fmt.Sprintf("%s/subscriptions", URL)

	return fetchData(urlStr, "GET", nil)
}

func updateSubscriptions(c *cli.Context) error {
	var idsToAdd, idsToDelete, currentIDs, list []int

	urlStr := fmt.Sprintf("%s/subscriptions", URL)
	eMail := common.SetStringId(c, "Please provide user e-mail")

	if c.String("add") != "" {
		idsToAdd = common.StringToIntArr(c.String("add"))
	}
	if c.String("delete") != "" {
		idsToDelete = common.StringToIntArr(c.String("delete"))
	}

	dataCurrent := listSubscriptionsData(c)
	dataParsed, err := subscriptionsParse(dataCurrent)
	common.ErrorCheck(err)

	for _, s := range dataParsed.Subscriptions {
		currentIDs = append(currentIDs, s.ServiceID)
	}
	sort.Ints(currentIDs)

	currentIDs = append(currentIDs, idsToAdd...)
	sort.Ints(currentIDs)
	result := common.RemoveIntDuplicates(currentIDs)

	if c.String("delete") == "" {
		list = result
	} else {
		list = common.DeleteSlicefromSlice(result, idsToDelete)
	}

	body := createSubscriptionBody(list, eMail)
	data := fetchData(urlStr, "PUT", body)

	if raw {
		println(data)

		return nil
	}

	res, err := subscriptionsParse(data)
	common.ErrorCheck(err)

	printSubscriptions(res)

	return nil
}

func printSubscriptions(subscriptions SubscriptionsResp) {
	color.Set(color.FgGreen)
	fmt.Println("# Firewall Rules Notification Services you are subscribed to:")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, fmt.Sprint("# ID\tName\tDescription\tE-Mail\tSign up Date"))
	for _, f := range subscriptions.Subscriptions {
		fmt.Fprintln(w, fmt.Sprintf("%v\t%s\t%s\t%s\t%s", f.ServiceID, f.ServiceName, f.Description, f.Email, f.SignupDate))
	}
	w.Flush()
}
