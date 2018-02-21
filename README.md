# Akamai CLI for Firewall Rules Notifications
*NOTE:* This tool is intended to be installed via the Akamai CLI package manager, which can be retrieved from the releases page of the [Akamai CLI](https://github.com/akamai/cli) tool.

<!--ts-->
   * [Akamai CLI for Firewall Rules Notifications](#akamai-cli-for-firewall-rules-notifications)
      * [Local Install, if you choose not to use the akamai package manager](#local-install-if-you-choose-not-to-use-the-akamai-package-manager)
      * [Credentials](#credentials)
      * [Overview](#overview)
      * [Main Command Usage](#main-command-usage)
         * [Raw output. (JSON)](#raw-output-json)
         * [Get command](#get-command)
         * [List commands](#list-commands)
            * [Services](#services)
            * [Subscriptions](#subscriptions)
            * [CIDR Blocks](#cidr-blocks)
         * [Update commands](#update-commands)

<!-- Added by: partamonov, at:  -->

<!--te-->

### Local Install, if you choose not to use the akamai package manager
If you want to compile it from source, you will need Go 1.9 or later, and the [Glide](https://glide.sh) package manager installed:
1. Fetch the package:
   `go get https://github.com/partamonov/akamai-cli-frn`
1. Change to the package directory:
   `cd $GOPATH/src/github.com/partamonov/akamai-cli-frn`
1. Install dependencies using Glide:
   `glide install`
1. Compile the binary:
   `go build -ldflags="-s -w -X main.version=X.X.X" -o akamai-frn`

### Credentials
In order to use this configuration, you need to:
* Set up your credential files as described in the [authorization](https://developer.akamai.com/introduction/Prov_Creds.html) and [credentials](https://developer.akamai.com/introduction/Conf_Client.html) sections of the getting started guide on developer.akamai.com.

Expects `default` section in .edgerc, can be changed via --section parameter

```
[default]
client_secret = XXXXXXXXXXXX
host = XXXXXXXXXXXX
access_token = XXXXXXXXXXXX
client_token = XXXXXXXXXXXX
```

## Overview
The Akamai Firewall Rules Notification Kit is a set of go libraries that wraps Akamai's {OPEN} APIs to let you manage who receives notifications about changes Akamai makes to IP addresses. You can subscribe or unsubscribe users to notifications for specific services, retrieve subscription and service information, and get CIDR block information with which to update your firewall rules.

## Main Command Usage
```shell
NAME:
   akamai frn - A CLI to interact with Akamai Firewall Rules Notifications

USAGE:
   akamai frn [global options] command [command options] [arguments...]

VERSION:
   X.X.X

AUTHORS:
   Petr Artamonov
   Rafal Pieniazek

COMMANDS:
     get       Get a specific [subcommand]] `ID`
     list, ls  Get a list of [subcommand]]
     update    Update [subcommand]] `ID`
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE   Location of the credentials FILE (default: "/Users/partamonov/.edgerc") [$AKAMAI_EDGERC]
   --no-color               Disable color output
   --raw                    Show raw output. It will be JSON format
   --section NAME, -s NAME  NAME of section to use from credentials file (default: "default") [$AKAMAI_EDGERC_SECTION]
   --help, -h               show help
   --version, -v            print the version
```

### Raw output. (JSON)
Any command can run with `--raw` parameter. Output will be as received from Akamai in JSON format

```shell
akamai frn --section custom-name --raw get service 3
{
  "serviceId" : 3,
  "serviceName" : "LOG_DELIVERY",
  "description" : "Log Delivery"
}
```

### Get command
Get can be used only to get service specific information

```shell
> akamai-frn get service <ID>
```

Example:

```shell
> akamai frn --section custom-name get service 3
# Firewall Rules Notification Services you are subscribed to:
# ID   Name           Description
3      LOG_DELIVERY   Log Delivery
```

### List commands

#### Services

You can list all available services for subscription

```shell
> akamai frn list services
# Firewall Rules Notification Services you are subscribed to:
# ID   Name                                       Description
1      FIRSTPOINT                                 Global Traffic Management
3      LOG_DELIVERY                               Log Delivery
4      SITE_SNAPSHOT                              Site Snapshot
...
```

#### Subscriptions
You can list all subscription assigned to your `API` account.

```shell
akamai frn list subscriptions
# Firewall Rules Notification Services you are subscribed to:
# ID   Name                               Description                                E-Mail Sign up Date
7      ESN                                Edge Staging Network                       XXXX   2018-02-16
8      SESN                               Secure Edge Staging Network                XXXX   2018-02-16

```

#### CIDR Blocks
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

```shell
> akamai-cli-frn list cidr --services "Secure Edge Staging Network,SiteShield + Secure Edge Staging Network" --only-addresses
1.2.3.4/32
5.6.7.8/32
```

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

```shell
If user had subscriptions for services with ID 6 and 7, then after running the following command he will be subscribed too same list
> akamai frn --add "3,6"--delete "3" user@e-mail.com
```
