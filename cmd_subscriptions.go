package main

import (
	"sort"

	common "github.com/apiheat/akamai-cli-common"
	"github.com/urfave/cli"
)

func cmdSubscriptions(c *cli.Context) error {
	return listSubscriptions(c)
}

func cmdUpdSubscriptions(c *cli.Context) error {
	return updateSubscriptions(c)
}

func listSubscriptions(c *cli.Context) error {
	data, _, err := apiClient.FRN.ListSubscriptions()
	common.ErrorCheck(err)

	common.OutputJSON(data)
	return nil
}

func updateSubscriptions(c *cli.Context) error {
	var idsToAdd, idsToDelete, currentIDs, list []int

	eMail := common.SetStringId(c, "Please provide user e-mail")

	if c.String("add") != "" {
		idsToAdd = common.StringToIntArr(c.String("add"))
	}
	if c.String("delete") != "" {
		idsToDelete = common.StringToIntArr(c.String("delete"))
	}

	dataCurrent, _, err := apiClient.FRN.ListSubscriptions()
	common.ErrorCheck(err)

	for _, s := range dataCurrent.Subscriptions {
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

	data, _, err := apiClient.FRN.UpdateSubscriptions(list, eMail)
	common.ErrorCheck(err)

	common.OutputJSON(data)
	return nil
}
