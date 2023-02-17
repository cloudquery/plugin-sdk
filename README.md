# CloudQuery Plugin SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/cloudquery/plugin-sdk#section-readme.svg)](https://pkg.go.dev/github.com/cloudquery/plugin-sdk#section-readme)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudquery/plugin-sdk)](https://goreportcard.com/report/github.com/cloudquery/plugin-sdk)
[![Unit tests](https://github.com/cloudquery/plugin-sdk/actions/workflows/unittest.yml/badge.svg)](https://github.com/cloudquery/plugin-sdk/actions/workflows/unittest.yml)

CloudQuery SDK enables building CloudQuery source and destination plugins.

Source plugins allows developers to extract information from third party APIs and enjoying built-in transformations, concurrency, logging, testing and database agnostic support via destination plugins.

Destinations plugins allows writing the data from any of the source plugins to an additional database, message queue, storage or any other destination without recompiling any of the source plugins.

The plugin SDK is imported as a dependency by CloudQuery plugins. When starting a new plugin, you should use the Scaffold tool.

## Getting Started & Documentation

* [CloudQuery Homepage](https://www.cloudquery.io)
* [CloudQuery Releases](https://github.com/cloudquery/cloudquery/releases?q=cli%2F&expanded=true)
* [Creating a new Plugin](https://www.cloudquery.io/docs/developers/creating-new-plugin) (Docs)
* [How to Write a CloudQuery Plugin](https://www.youtube.com/watch?v=3Ka_Ob8E6P8) (Video 🎥)

## Supported plugins

<https://www.cloudquery.io/plugins>

If you want us to add a new plugin or resource please open an [Issue](https://github.com/cloudquery/cloudquery/issues).
