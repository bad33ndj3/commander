# Commander

[![Go](https://github.com/bad33ndj3/commander/actions/workflows/go.yml/badge.svg)](https://github.com/bad33ndj3/commander/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bad33ndj3/commander)](https://goreportcard.com/report/github.com/bad33ndj3/commander)
[![codecov](https://codecov.io/gh/bad33ndj3/commander/branch/main/graph/badge.svg)](https://codecov.io/gh/bad33ndj3/commander)
[![GoDoc](https://godoc.org/github.com/bad33ndj3/commander?status.svg)](https://godoc.org/github.com/bad33ndj3/commander)

Commander is a simple, structured CLI framework built as an alternative to Make/Mage for Go projects. It was created out of a need for a more Go-idiomatic build tool that leverages Go's type system and provides better IDE support.

> ⚠️ **Note**: This project is currently in an experimental stage and may not be production-ready. Features may change, and there might be bugs. Use it at your own risk.


I would advise you to use https://taskfile.dev/ since its a complete package.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Quick Start](#quick-start)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Defining Commands](#defining-commands)
    - [Example Argument Structs](#example-argument-structs)
    - [Writing Handler Functions](#writing-handler-functions)
  - [Setting Up Commander](#setting-up-commander)
  - [Running the Application](#running-the-application)
  - [Using the Help System](#using-the-help-system)
- [Struct Tag Options](#struct-tag-options)
- [Handler Function Rules](#handler-function-rules)
- [License](#license)
- [Conclusion](#conclusion)

## Features

- **Type-Safe Arguments**: Define command arguments using structs with tags.
- **Automatic Flag Parsing**: Struct fields are converted into CLI flags.
- **Command Categories**: Organize commands into logical groups.
- **Built-in Help System**: Automatically generated help messages with colorized output.
- **Standard Library Only**: No external dependencies required.

## Quick Start

Here's a minimal example to get started:

```go
package main

import (
    "context"
    "fmt"

    "github.com/bad33ndj3/commander"
)

func main() {
    cmdr := commander.New()

    rootCat := cmdr.AddCategory("Root")
    rootCat.Register(&commander.Command{
        Name:        "hello",
        Description: "Prints Hello World",
        Handler:     helloHandler,
    })

    if err := cmdr.Run(); err != nil {
        fmt.Println(err)
    }
}

func helloHandler(ctx context.Context) {
    fmt.Println("Hello, World!")
}
```

## Getting Started

### Installation

Include the `commander` package in your project by importing it.

```sh
go get github.com/bad33ndj3/commander
```

### Defining Commands

Commands are defined using handler functions and argument structs. Handler functions must accept `context.Context` as the first parameter and an optional struct for arguments.

#### Example Argument Structs

```go
type BuildArgs struct {
    Debug   bool   `flag:"debug" default:"false" usage:"Enable debug mode"`
    Output  string `flag:"output" default:"./bin" usage:"Output directory"`
    Version string `flag:"version" default:"dev" usage:"Build version"`
}

type TestArgs struct {
    Verbose bool   `flag:"verbose" default:"false" usage:"Enable verbose output"`
    Pattern string `flag:"pattern" default:"./..." usage:"Test pattern to run"`
}
```

### Setting Up Commander

Initialize `Commander`, add categories, and register commands:

```go
func main() {
    cmdr := commander.New()

    buildCat := cmdr.AddCategory("Build")
    testCat := cmdr.AddCategory("Test")

    buildCat.Register(&commander.Command{
        Name:        "build",
        Description: "Build the project",
        Handler:     buildHandler,
    })

    testCat.Register(&commander.Command{
        Name:        "test",
        Description: "Run tests",
        Handler:     testHandler,
    })

    if err := cmdr.Run(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

### Using the Help System

Display general help:

```sh
go run main.go help
```

Display help for a specific command:

```sh
go run main.go help build
```

## Struct Tag Options

- `flag:"name"`: Custom flag name (default is field name in lowercase).
- `default:"value"`: Default value for the flag.
- `usage:"description"`: Help text displayed in the help message.

## Handler Function Rules

1. **First Parameter**: Must be `context.Context`.
2. **Second Parameter**: Optional struct for command arguments.
3. **Supported Field Types**: `bool`, `int`, `string`.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Conclusion

The `commander` package simplifies the creation of CLI applications by using Go's native features. By defining commands with clear handler functions and argument structs, you can build intuitive and maintainable command-line tools efficiently.

For more examples and advanced usage, refer to the [package documentation](https://pkg.go.dev/github.com/bad33ndj3/commander) or the [examples directory](https://github.com/bad33ndj3/commander/tree/main/examples) in the repository.
