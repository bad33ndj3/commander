# Commander

[![Go](https://github.com/bad33ndj3/commander/actions/workflows/go.yml/badge.svg)](https://github.com/bad33ndj3/commander/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bad33ndj3/commander)](https://goreportcard.com/report/github.com/bad33ndj3/commander)
[![codecov](https://codecov.io/gh/bad33ndj3/commander/branch/main/graph/badge.svg)](https://codecov.io/gh/bad33ndj3/commander)
[![GoDoc](https://godoc.org/github.com/bad33ndj3/commander?status.svg)](https://godoc.org/github.com/bad33ndj3/commander)

Commander is a simple, structured CLI framework built as an alternative to Make/Mage for Go projects. It was created out of a need for a more Go-idiomatic build tool that leverages Go's type system and provides better IDE support.

> ⚠️ **Note**: This is an experimental project and not production-ready. Use at your own risk.

## Why Another Build Tool?

While Make is powerful and Mage is Go-native, I found myself wanting:
- Better type safety for command arguments
- Structured command organization
- Built-in help system with automatic flag documentation
- More intuitive command grouping
- Native Go syntax without string parsing
