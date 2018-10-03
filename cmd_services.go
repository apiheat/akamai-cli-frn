package main

import (
	common "github.com/apiheat/akamai-cli-common"
	"github.com/urfave/cli"
)

func cmdServices(c *cli.Context) error {
	return listServices(c)
}

func cmdGetService(c *cli.Context) error {
	return getService(c)
}

func listServices(c *cli.Context) error {
	data, _, err := apiClient.FRN.ListServices()
	common.ErrorCheck(err)

	common.OutputJSON(data)
	return nil
}

func getService(c *cli.Context) error {
	id := common.SetIntID(c, "Please provide ID for a service")

	data, _, err := apiClient.FRN.ListService(id)
	common.ErrorCheck(err)

	common.OutputJSON(data)

	return nil
}
