# Akamai CLI for Firewall Rules Notifications

The Akamai Firewall Rules Notification Kit is a set of go libraries that wraps Akamai's {OPEN} APIs to let you manage who receives notifications about changes Akamai makes to IP addresses. You can subscribe or unsubscribe users to notifications for specific services, retrieve subscription and service information, and get CIDR block information with which to update your firewall rules.

Should you miss something we *gladly accept patches* :)

CLI uses custom [Akamai API client](https://github.com/apiheat/go-edgegrid)

## Configuration & Installation

### Credentials

Set up your credential files as described in the [authorization](https://developer.akamai.com/introduction/Prov_Creds.html) and [credentials](https://developer.akamai.com/introduction/Conf_Client.html) sections of the getting started guide on developer.akamai.com.

Tools expect proper format of sections in edgerc file which example is shown below

*NOTE:* Default file location is *~/.edgerc*

```
[default]
client_secret = XXXXXXXXXXXX
host = XXXXXXXXXXXX
access_token = XXXXXXXXXXXX
client_token = XXXXXXXXXXXX
```

In order to change section which is being actively used you can

* change it via `--config parameter` of the tool itself
* change it via env variable `export AKAMAI_EDGERC_CONFIG=/Users/jsmitsh/.edgerc`

In order to change section which is being actively used you can

* change it via `--section parameter` of the tool itself
* change it via env variable `export AKAMAI_EDGERC_SECTION=mycustomsection`

>*NOTE:* Make sure your API client do have appropriate scopes enabled

### Installation

The tool can be used as a stand-alone binary or in conjuction with [Akamai CLI](https://developer.akamai.com/cli).

#### Akamai-cli ( recommended )

Execute the following from console

```shell
> akamai install https://github.com/apiheat/akamai-cli-frn
```

#### Stand-alone

As part of automated releases/builds you can download latest version from the project release page

## Usage

```shell
NAME:
   akamai-frn - A CLI to interact with Akamai Firewall Rules Notifications

USAGE:
   akamai-frn [global options] command [command options] [arguments...]

VERSION:
   X.X.X

AUTHORS:
   Petr Artamonov
   Rafal Pieniazek

COMMANDS:
     get       Get a specific [subcommand]] `ID`
     list, ls  Get a list of [subcommand]]
     update    Update [subcommand]]
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE   Location of the credentials FILE (default: "/Users/USER_NAME/.edgerc") [$AKAMAI_EDGERC_CONFIG]
   --debug value            Debug Level [$AKAMAI_EDGERC_DEBUGLEVEL]
   --section NAME, -s NAME  NAME of section to use from credentials file (default: "default") [$AKAMAI_EDGERC_SECTION]
   --help, -h               show help
   --version, -v            print the version
```

### CIDR Blocks Commands
You can list all CIDR Blocks associated with services to which you subscribed

```shell
> akamai frn list cidr
# Firewall Rules Notification CIDR Blocks you are subscribed to:
# ID    Service Name (ID)                               CIDR      Port     Active       Last Action
303     Secure Edge Staging Network (8)                 XXXXXX    80,443   2007-10-13   update
306     Secure Edge Staging Network (8)                 XXXXXX    80,443   2007-10-13   update
1241    Edge Staging Network (7)                        XXXXXX    80,443   2008-11-25   update
1601    Secure Edge Staging Network (8)                 XXXXXX    80,443   2009-04-21   add
...
```

You can filter results by Akamai API supported flags:

* --last-action:    Return only CIDR blocks with a change status of add, update, or delete.
* --effective-date: The ISO 8601 date(YYYY-MM-DD) the CIDR block starts serving traffic to your origin

Also you can filter services in output by `name` with `--services` parameter by provided comma separated string with names

```shell
akamai frn list cidr --services "SiteShield + Secure Edge Staging Network,Edge Staging Network"
# Firewall Rules Notification CIDR Blocks you are subscribed to:
# Showing CIDR Blocks only for: SiteShield + Secure Edge Staging Network
# ID    Service Name (ID)                               CIDR               Port     Active       Last Action
13231   SiteShield + Secure Edge Staging Network (32)   XXXXXXXX           80,443   2016-11-10   add
14540   Edge Staging Network (7)                        XXXXXXXX           80,443   2017-11-13   add
...
```

Some times you may need to get only CIDRs as output to process then with xargs or in any other way. You can do that with `--only-addresses` flag

### Update commands

You can subscribe to any service you want and unsubscribe too with `update` command

There are 2 flags:

* `--add`    comma(',') separated list of Service IDs to which you want to subscribe
* `--delete` comma(',') separated list of Service IDs to which you want to unsubscribe

Please take a note that the list of required services created in the following way:

1. We get your `current` subscriptions from Akamai
1. We append `add` list to `current`
1. We sort and uniq the `result` list
1. We remove all elements present in `delete` list from `result` list
1. We send to Akamai the `result` list

If user had subscriptions for services with ID 6 and 7, then after running the following command he will be subscribed too same list

```shell
> akamai frn --add "3,6"--delete "3" user@e-mail.com
```

## Development

In order to develop the tool with us do the following:

1. Fork repository
1. Clone it to your folder ( within *GO* path )
1. Ensure you can restore dependencies by running

   ```shell
   dep ensure
   ```

1. Make necessary changes
1. Make sure solution builds properly ( feel free to add tests )

   ```shell
   go build -ldflags="-s -w -X main.appVer=1.2.3 -X main.appName=$(basename `pwd`)"
   ```
