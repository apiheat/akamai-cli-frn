package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdServices(c *cli.Context) error {
	return listServices(c)
}

func cmdGetService(c *cli.Context) error {
	return getService(c)
}

func listServices(c *cli.Context) error {
	urlStr := fmt.Sprintf("%s/services", URL)
	data := fetchData(urlStr, "GET", nil)

	if raw {
		println(data)

		return nil
	}

	result, err := servicesParse(data)
	errorCheck(err)

	printServices(result)

	return nil
}

func getService(c *cli.Context) error {
	id := setID(c)

	urlStr := fmt.Sprintf("%s/services/%s", URL, id)
	data := fetchData(urlStr, "GET", nil)

	if raw {
		println(data)

		return nil
	}

	result, err := servicesParse("[" + data + "]")
	errorCheck(err)

	printServices(result)

	return nil
}

func printServices(services Services) {
	color.Set(color.FgGreen)
	fmt.Println("# Firewall Rules Notification Services you are subscribed to:")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, fmt.Sprint("# ID\tName\tDescription"))
	for _, f := range services {
		fmt.Fprintln(w, fmt.Sprintf("%v\t%s\t%s", f.ServiceID, f.ServiceName, f.Description))
	}
	w.Flush()
}
