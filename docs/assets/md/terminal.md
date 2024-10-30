
### Using Linux Terminal

```bash

curl -fsSL https://github.com/mguptahub/nanodns/releases/download/v1.1.3/nanodns-linux-amd64 -o /usr/local/bin/nanodns && chmod +x /usr/local/bin/nanodns

nanodns start


```

Help Command
```
nanodns --help
```

```
Usage: nanodns [command | options]

commands:
  start                              Run the binary as a daemon
  stop                               Stop the running daemon service
  status                             Show service status
  logs                               Show service logs

options:
  -v | --version                     Show the binary version
  -a | --action-logs                 Show the action logs. This works with the logs command

```
