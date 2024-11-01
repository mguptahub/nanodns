
### Using Linux Terminal

Install using the script

```bash
curl -fsSL https://nanodns.mguptahub.com/install.sh | sh -s -- --install
```

Start using the script

```bash
# Check the values in /usr/local/share/nanodns.env before starting
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
