# Contributing to CloudQuery

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

This is the repository for CloudQuery SDK. If you are looking for CloudQuery CLI and plugins take a look at our [monorepo](https://github.com/cloudquery/cloudquery)

## Links

- [Homepage](https://cloudquery.io)
- [Documentation](https://docs.cloudquery.io)
- [CLI and Plugins Mono Repo](https://github.com/cloudquery/cloudquery)
- [Discord](https://cloudquery.io/discord)

## Development

### gRPC

CloudQuery has a pluggable architecture and uses [gRPC](https://grpc.io/) to communicate between source plugins, CLI and destination plugins. To develop a new plugin for CloudQuery, you don’t need to understand the inner workings of gRPC as those are abstracted via the [plugin-sdk](#cloudquery-plugin-sdk-repository).

If you want to make any changes to the protocol between plugins and the CLI you will need to install all [go-gRPC prerequisites](https://grpc.io/docs/languages/go/quickstart/#prerequisites).

All protobuf files and auto-generated Go gRPC server/client are located under [./internal/pb](./internal/pb/).

To regenerate new Go gRPC client and server run `make gen-proto`.

To provide a better API which abstracts the protobuf structs we maintain our own clients at [./clients](./clients) and servers at [./plugins/](./plugins/) so make sure to adjust those acordingly.

### Packages

- [serve](./serve) command line APIs to start serving plugins.
- [plugins](./plugins/) main plugin APIs (source/dest).
- [schema](./schema/) tables, columns and supported types
- [codegen](./codegen) apis to generate CloudQuery tables from Go structs
- [faker](./faker) useful api to fake structs for source plugin mock tests

### Tests

Run `make test` to run all unit-tests

### Lint

We use `golangci-lint` as our linter. Configuration available in [./golangci.yml] to run lint `make lint`
