package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	common "github.com/apiheat/akamai-cli-common"
	edgegrid "github.com/apiheat/go-edgegrid"
	log "github.com/sirupsen/logrus"
	"github.com/thedevsaddam/gojsonq"

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
	var services string
	if filter == "?" {
		filter = ""
	}
	data, _, err := apiClient.FRN.ListCIDRBlocks(filter)
	common.ErrorCheck(err)

	if c.String("services") != "" {
		services = c.String("services")
	}

	if c.Bool("only-addresses") {
		printCidrs(data, services, c.String("output"))

		return nil
	}

	printData(data, services, c.String("output"))

	return nil
}

func jQ(cidrs *edgegrid.FRNCidrs) *gojsonq.JSONQ {
	b, err := json.Marshal(cidrs)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return gojsonq.New().FromString(string(b))
}

func printData(cidrs *edgegrid.FRNCidrs, search, outputFormat string) {
	searchSlice := searchServices(search)

	if outputFormat == "json" {
		jq := jQ(cidrs)
		res := jq.Where("description", "in", searchSlice)
		common.OutputJSON(res.Get())

		return
	}

	color.Set(color.FgGreen)
	fmt.Println("# Firewall Rules Notification CIDR Blocks you are subscribed to:")
	color.Unset()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, fmt.Sprint("# ID\tService Name (ID)\tCIDR\tPort\tActive\tLast Action"))
	for _, f := range *cidrs {
		if common.IsStringInSlice(f.Description, searchSlice) {
			fmt.Fprintln(w, fmt.Sprintf("%v\t%s (%v)\t%s\t%s\t%s\t%s",
				f.CidrID, f.Description, f.ServiceID, f.Cidr+f.CidrMask, f.Port, f.EffectiveDate, f.LastAction))
		}
	}
	w.Flush()
}

func printCidrs(cidrs *edgegrid.FRNCidrs, search, outputFormat string) {
	searchSlice := searchServices(search)

	if outputFormat == "json" {
		jq := jQ(cidrs)
		res := jq.Where("description", "in", searchSlice).Pluck("cidr")
		common.OutputJSON(res)

		return
	}

	color.Set(color.FgGreen)
	fmt.Println("# Firewall Rules Notification CIDR Blocks you are subscribed to:")
	color.Unset()

	ips := make([]string, len(*cidrs))
	for _, f := range *cidrs {
		if common.IsStringInSlice(f.Description, searchSlice) {
			ips = append(ips, f.Cidr+f.CidrMask)
		}
	}
	sort.Strings(ips)
	cidr := common.RemoveStringDuplicates(ips)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, fmt.Sprint("# Total number of CIDR Blocks: ", len(cidr)))
	for _, ip := range cidr {
		fmt.Fprintln(w, fmt.Sprintf("%s", ip))
	}

	w.Flush()

	return
}

func searchServices(searchStr string) (searchSlice []string) {
	if searchStr != "" {
		color.Set(color.FgYellow)
		fmt.Printf("# Showing CIDR Blocks only for: %s\n", searchStr)
		color.Unset()

		searchSlice = strings.Split(searchStr, ",")
	}

	return searchSlice
}
