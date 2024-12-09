package main

import (
	"fmt"
	"os"

	"github.com/bad33ndj3/commander"
	"github.com/bad33ndj3/commander/example/golang/tools/cmder/cmd"
)



func main() {
	cmdr := commander.New()

	// Create categories
	serverCat := cmdr.AddCategory("Server")
	goCat := cmdr.AddCategory("Go")

	// Server commands
	serverCat.Register(&commander.Command{
		Name:        "start",
		Description: "Start the development server",
		Handler:     cmd.StartServerCommand,
	})

	// Test commands
	goCat.Register(&commander.Command{
		Name:        "test",
		Description: "Run tests",
		Handler:     cmd.TestCommand,
	})

	// Lint commands
	goCat.Register(&commander.Command{
		Name:        "lint",
		Description: "Run linters",
		Handler:     cmd.LintCommand,
	})

	if err := cmdr.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
