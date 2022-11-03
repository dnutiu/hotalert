# Hot Alert

## Introduction
Hotalert is a command line tool that for task execution and configuration. Tasks and alerts are defined
in yaml files and the program parses the files, executes the tasks and emits alerts when the tasks conditions are met.

For example if you want to send a notification to your mobile phone when a keyword is found on a website you can use the
`webhook_discord` alerter and the `scrape_web` task.

Example:

```yaml
tasks:
  - options:
      url: [...]
      keywords: ["Episode 10", "Episode 11"]
    timeout: 10
    alerter: "webhook_discord"
    task: "scrape_web"
alerts:
  webhook_discord:
    webhook: https://discord.com/api/webhooks/[...]
    # $keywords can be used as a placeholder in the message, and it will be replaced with the actual keywords.
    message: "Hi, the keyword(s) `$keywords` was found on [...]"
```

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