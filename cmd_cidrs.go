package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdCidrs(c *cli.Context) error {
	filterStr := "?"
	if c.String("last-action") != "" {
		filterStr = filterStr + "lastAction=" + c.String("last-action")
		if c.String("effective-date") != "" {
			filterStr = filterStr + "&"
		}
	}

	if c.String("effective-date") != "" {
		// TODO: Check if date is in ISO 8601 format
		r := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)
		matches := r.FindAllString(c.String("effective-date"), -1)

		if len(matches) == 0 {
			errStr := fmt.Sprintf("Date format is not correct. Supported YYYY-MM-DD. You provided: %s", c.String("effective-date"))
			log.Fatalf(errStr)
		}

		filterStr = filterStr + "effectiveDateGt=" + c.String("effective-date")
	}

	return listCidrs(c, filterStr)
}

func listCidrs(c *cli.Context, filter string) error {
	urlStr := fmt.Sprintf("%s/cidr-blocks", URL)
	if filter != "?" {
		urlStr = fmt.Sprintf("%s/cidr-blocks%s", URL, filter)
	}

	data := fetchData(urlStr, "GET", nil)

	result, err := cidrsParse(data)
	errorCheck(err)

	printCidrs(result)

	return nil
}

func printCidrs(cidrs Cidrs) {
	color.Set(color.FgGreen)
	fmt.Println("Firewall Rules Notification CIDR Blocks you are subscribed to:")
	color.Unset()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, fmt.Sprint("ID\tService Name (ID)\tCIDR\tPort\tActive\tLast Action"))
	for _, f := range cidrs {
		fmt.Fprintln(w, fmt.Sprintf("%v\t%s (%v)\t%s\t%s\t%s\t%s",
			f.CidrID, f.Description, f.ServiceID, f.Cidr+f.CidrMask, f.Port, f.EffectiveDate, f.LastAction))
	}
	w.Flush()
}
