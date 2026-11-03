## CLI

Command-line interface for managing the `usernet` control plane.

### Build

We use [Task](https://taskfile.dev/) (`go-task`) to automate build commands.

1. Install Task:

```bash
# Via Go or with your packet manager
go install github.com/go-task/task/v3/cmd/task@latest
```

2. Build the binary:

```bash
task build
```

### Usage

```bash
# View available commands and flags
./usernet --help

# Apply network topology configuration
./usernet apply -f config.yaml

# Run one or more scenario files
./usernet run scenario_1.yaml scenario_2.yaml

# Apply configuration and run scenarios in one step
./usernet run -f config.yaml scenario_1.yaml
```
