package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

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
	printSubscriptions(listSubscriptionsData(c))
	return nil
}

func listSubscriptionsData(c *cli.Context) SubscriptionsResp {
	urlStr := fmt.Sprintf("%s/subscriptions", URL)
	data := fetchData(urlStr, "GET", nil)

	result, err := subscriptionsParse(data)
	errorCheck(err)

	return result
}

func updateSubscriptions(c *cli.Context) error {
	urlStr := fmt.Sprintf("%s/subscriptions", URL)
	eMail := setEmail(c)

	var idsToAdd, idsToDelete, currentIDs, list []int
	if c.String("add") != "" {
		idsToAdd = strToIntArr(c.String("add"))
	}
	if c.String("delete") != "" {
		idsToDelete = strToIntArr(c.String("delete"))
	}

	for _, s := range listSubscriptionsData(c).Subscriptions {
		currentIDs = append(currentIDs, s.ServiceID)
	}
	sort.Ints(currentIDs)

	currentIDs = append(currentIDs, idsToAdd...)
	sort.Ints(currentIDs)
	result := removeDuplicates(currentIDs)

	if c.String("delete") == "" {
		list = result
	} else {
		list = deleteSlicefromSlice(result, idsToDelete)
	}

	body := createSubscriptionBody(list, eMail)
	data := fetchData(urlStr, "PUT", body)

	res, err := subscriptionsParse(data)
	errorCheck(err)

	printSubscriptions(res)

	return nil
}

func printSubscriptions(subscriptions SubscriptionsResp) {
	color.Set(color.FgGreen)
	fmt.Println("Firewall Rules Notification Services you are subscribed to:")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, fmt.Sprint("ID\tName\tDescription\tE-Mail\tSign up Date"))
	for _, f := range subscriptions.Subscriptions {
		fmt.Fprintln(w, fmt.Sprintf("%v\t%s\t%s\t%s\t%s", f.ServiceID, f.ServiceName, f.Description, f.Email, f.SignupDate))
	}
	w.Flush()
}
