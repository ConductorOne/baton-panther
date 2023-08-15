# baton-panther

`baton-panther` is a connector for Panther built using the [Baton SDK](https://github.com/conductorone/baton-sdk). It communicates with the Panther API to sync data about users and roles.

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Getting Started

## Prerequisites

1. Panther API token that can be created in `Settings -> API Tokens` with permissions:
 - Read API Token Info
 - Read Panther Settings Info
 - Read User Info
2. API URL used to fetch the data from Panther. It is located in `Settings -> API Tokens`.


## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-panther
baton-panther
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_TOKEN=pantherApiToken BATON_URL=pantherApiUrl ghcr.io/conductorone/baton-panther:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-panther/cmd/baton-panther@main

BATON_TOKEN=pantherApiToken BATON_URL=pantherApiUrl
baton resources
```

# Data Model

`baton-panther` pulls down information about the following Panther resources:
- Users
- Roles

# Contributing, Support, and Issues

We started Baton because we were tired of taking screenshots and manually building spreadsheets. We welcome contributions, and ideas, no matter how small -- our goal is to make identity and permissions sprawl less painful for everyone. If you have questions, problems, or ideas: Please open a Github Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-panther` Command Line Usage

```
baton-panther

Usage:
  baton-panther [flags]
  baton-panther [command]

Available Commands:
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --client-id string              The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string          The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                   The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                          help for baton-panther
      --log-format string             The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string              The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --token string                  API token used to authenticate to the Panther API. ($BATON_TOKEN)
      --url string                    API url of your panther account. ($BATON_URL)
  -v, --version                       version for baton-panther

Use "baton-panther [command] --help" for more information about a command.

```
