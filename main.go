package main

import (
	"os"
	"sort"

	common "github.com/apiheat/akamai-cli-common"
	edgegrid "github.com/apiheat/go-edgegrid"
	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

var (
	apiClient       *edgegrid.Client
	appName, appVer string
)

// Constants
const (
	padding = 3
)

func main() {
	app := common.CreateNewApp(appName, "A CLI to interact with Akamai Firewall Rules Notifications", appVer)
	app.Flags = common.CreateFlags()

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "Get a list of [subcommand]]",
			Subcommands: []cli.Command{
				{
					Name:   "services",
					Usage:  "... services you are subscribed to",
					Action: cmdServices,
				},
				{
					Name:   "subscriptions",
					Usage:  "... subscriptions you are created for yourself and other users",
					Action: cmdSubscriptions,
				},
				{
					Name:  "cidr",
					Usage: "... cidr blocks for all services you are subscribed to",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "last-action",
							Value: "",
							Usage: "Return only CIDR blocks with a change status of add, update, or delete.",
						},
						cli.StringFlag{
							Name:  "effective-date",
							Value: "",
							Usage: "The ISO 8601 date(YYYY-MM-DD) the CIDR block starts serving traffic to your origin",
						},
						cli.StringFlag{
							Name:  "services",
							Value: "",
							Usage: "Return CIDR blocks  only for comma separated list of services",
						},
						cli.StringFlag{
							Name:  "output",
							Value: "json",
							Usage: "Type of output. json or table is supported.",
						},
						cli.BoolFlag{
							Name:  "only-addresses",
							Usage: "Show only CIDR Blocks addresses.",
						},
					},
					Action: cmdCidrs,
				},
			},
		},
		{
			Name:  "get",
			Usage: "Get a specific [subcommand]] `ID`",
			Subcommands: []cli.Command{
				{
					Name:   "service",
					Usage:  "... service",
					Action: cmdGetService,
				},
			},
		},
		{
			Name:  "update",
			Usage: "Update [subcommand]]",
			Subcommands: []cli.Command{
				{
					Name:  "subscriptions",
					Usage: "update subscriptions [parameters] User_Email",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "add",
							Value: "",
							Usage: "Specify comma(',') separated list of Service IDs to which you want to subscribe",
						},
						cli.StringFlag{
							Name:  "delete",
							Value: "",
							Usage: "Specify comma(',') separated list of Service IDs to which you want to unsubscribe",
						},
					},
					Action: cmdUpdSubscriptions,
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Before = func(c *cli.Context) error {
		var err error

		apiClient, err = common.EdgeClientInit(c.GlobalString("config"), c.GlobalString("section"), c.GlobalString("debug"))

		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
