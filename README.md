# Akamai CLI for Firewall Rules Notifications
*NOTE:* This tool is intended to be installed via the Akamai CLI package manager, which can be retrieved from the releases page of the [Akamai CLI](https://github.com/akamai/cli) tool.

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
   --section NAME, -s NAME  NAME of section to use from credentials file (default: "default") [$AKAMAI_EDGERC_SECTION]
   --help, -h               show help
   --version, -v            print the version
```

### Get commands
```shell
NAME:
   akamai-frn get - Get a specific [subcommand]] `ID`

USAGE:
   akamai-frn get command [command options] [arguments...]

COMMANDS:
     service  ... service `ID`

OPTIONS:
   --help, -h  show help
```

Example:
```shell
> akamai-cli-frn get service 3
```

### List commands
```shell
NAME:
   akamai-frn list - Get a list of [subcommand]]

USAGE:
   akamai-frn list command [command options] [arguments...]

COMMANDS:
     services       ... services you are subscribed to
     subscriptions  ... subscriptions you are created for yourself and other users
     cidr           ... cidr blocks for all services you are subscribed to

OPTIONS:
   --help, -h  show help
```

Example:
```shell
> akamai-cli-frn list cidr --services "Secure Edge Staging Network,SiteShield + Secure Edge Staging Network" --only-addresses
```

### Update commands
```shell
NAME:
   akamai-frn update - Update [subcommand]] `ID`

USAGE:
   akamai-frn update command [command options] [arguments...]

COMMANDS:
     subscriptions  ... subscribe or unsubscribe users to services

OPTIONS:
   --help, -h  show help
```

Example:
```shell
> akamai-cli-frn update subscriptions --add "7,8,32" --delete "6,3" test@user.com
```