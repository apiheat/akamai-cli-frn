package main

import (
	"os"
	"sort"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

var (
	colorOn, raw              bool
	version, appName          string
	configSection, configFile string
	edgeConfig                edgegrid.Config
)

// Constants
const (
	URL     = "/firewall-rules-manager/v1"
	padding = 3
)

// Services data representation
type Services []struct {
	ServiceID   int    `json:"serviceId"`
	ServiceName string `json:"serviceName"`
	Description string `json:"description"`
}

// Cidrs data representation
type Cidrs []struct {
	CidrID        int         `json:"cidrId"`
	ServiceID     int         `json:"serviceId"`
	ServiceName   string      `json:"serviceName"`
	Description   string      `json:"description"`
	Cidr          string      `json:"cidr"`
	CidrMask      string      `json:"cidrMask"`
	Port          string      `json:"port"`
	CreationDate  string      `json:"creationDate"`
	EffectiveDate string      `json:"effectiveDate"`
	ChangeDate    interface{} `json:"changeDate"`
	MinIP         string      `json:"minIp"`
	MaxIP         string      `json:"maxIp"`
	LastAction    string      `json:"lastAction"`
}

// SubscriptionsResp data representation
type SubscriptionsResp struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

// Subscription data representation
type Subscription struct {
	ServiceID   int    `json:"serviceId"`
	ServiceName string `json:"serviceName,omitempty"`
	Description string `json:"description,omitempty"`
	Email       string `json:"email"`
	SignupDate  string `json:"signupDate,omitempty"`
}

func main() {
	_, inCLI := os.LookupEnv("AKAMAI_CLI")

	appName = "akamai-frn"
	if inCLI {
		appName = "akamai frn"
	}

	app := cli.NewApp()
	app.Name = appName
	app.HelpName = appName
	app.Usage = "A CLI to interact with Akamai Firewall Rules Notifications"
	app.Version = version
	app.Copyright = ""
	app.Authors = []cli.Author{
		{
			Name: "Petr Artamonov",
		},
		{
			Name: "Rafal Pieniazek",
		},
	}

	dir, _ := homedir.Dir()
	dir += string(os.PathSeparator) + ".edgerc"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "section, s",
			Value:       "default",
			Usage:       "`NAME` of section to use from credentials file",
			Destination: &configSection,
			EnvVar:      "AKAMAI_EDGERC_SECTION",
		},
		cli.StringFlag{
			Name:        "config, c",
			Value:       dir,
			Usage:       "Location of the credentials `FILE`",
			Destination: &configFile,
			EnvVar:      "AKAMAI_EDGERC",
		},
		cli.BoolFlag{
			Name:        "no-color",
			Usage:       "Disable color output",
			Destination: &colorOn,
		},
		cli.BoolFlag{
			Name:        "raw",
			Usage:       "Show raw output. It will be JSON format",
			Destination: &raw,
		},
	}

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
		if c.Bool("no-color") {
			color.NoColor = true
		}

		edgeConfig = config(configFile, configSection)
		return nil
	}

	app.Run(os.Args)
}
