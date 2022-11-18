# hotalert

## Introduction
hotalert is a command line tool for task execution and alerting. Tasks and alerts are defined in yaml files.

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