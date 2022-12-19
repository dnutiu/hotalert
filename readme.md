# hotalert

## Introduction
hotalert is a command line tool for task execution and alerting. Tasks and alerts are defined in yaml files.

A sample use case is to run it on PC or a Raspberry PI and get alerts on your mobile device when certain keywords apear in an URL.

## Installation

**Snapstore**

Get the snap: https://snapcraft.io/hotalert/preview

### Example
If you want to send a notification to your mobile phone when a keyword is found on a website you can use the
`webhook_discord` alerter and the `web_scrape` task function.

Example:

```yaml
tasks:
  - options:
      url: [...]
      keywords: ["Episode 10", "Episode 11"]
    timeout: 10
    alerter: "webhook_discord"
    function: "web_scrape"
alerts:
  webhook_discord:
    webhook: https://discord.com/api/webhooks/[...]
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword(s) `$keywords` was found on [...]"
```

Viewing the help:

```bash
Hotalert is a command line tool that for task execution and configuration. Tasks and alerts are defined 
in yaml files and the program parses the files, executes the tasks and emits alerts when the tasks conditions are met.

Usage:
  hotalert [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  directory   execute each yaml file from a directory
  file        execute tasks from a single file
  help        Help about any command

Flags:
  -h, --help   help for hotalert

Use "hotalert [command] --help" for more information about a command.
```

Running the program
```bash
./hotalert file test_file.yaml
2022-12-19T22:12:22.435+0200	INFO	alert/discord.go:82	Alert posted:
BEGIN
Hi, the keyword(s) `Episode 10,Episode 11` was found on [...]
END
2022-12-19T22:12:22.435+0200	INFO	cmd/file.go:60	Done
```

Output: 

![Discord preview](/docs/discord_alert.png)

### Available task functions

#### web_scrape

The web_scrape task scrapes a web page by issuing a GET request and parses the response to look for keywords.

**Options**:
- url (string) - The url to scrape.
- keywords (array[string]) - A list of keyword strings to look for.

### Development

To build the program for Linux under Linux use the following command:

```bash
GOOS=linux GOARCH=arm64 go build
```

If you're using Windows use:
```bash
$env:GOOS = "linux"
$env:GOARCH = "arm"
$env:GOARM = 5 

go build -o hotalert .
```
