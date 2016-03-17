# Slack Bot That Says GitHub Opened Pull Requests

This bot gets opened pull requests from your github repository and post their contents to   
you need golang develop environments.

## How to use

### build
```
go build -o PullReqBot -ldflags "-X main.token=<your github api token> -X main.endpoint=<github root endpoint> -X main.username=<github name> -X main.hookurl=<slack hookurl> main.channel=<slack channel>" main.go
```

### execute
```
$ ./PullReqBot
```

### auto execute
You can use jenkins or cron ... etc.
