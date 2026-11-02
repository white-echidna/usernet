## CLI

Command-line interface for managing the `usernet` control plane.

### Build

```bash
go build -o usernet ./cmd/usernet
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
