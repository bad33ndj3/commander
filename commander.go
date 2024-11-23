package commander

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

// Common errors
var (
	ErrNoSubcommand   = errors.New("no subcommand provided")
	ErrInvalidHandler = errors.New("handler must be a function accepting context.Context")
	ErrUnknownCommand = errors.New("unknown command")
)

// Style constants
const (
	flagTag    = "flag"
	defaultTag = "default"
	usageTag   = "usage"
)

// ANSI color codes
type colorCode string

const (
	colorReset  colorCode = "\033[0m"
	colorBold   colorCode = "\033[1m"
	colorRed    colorCode = "\033[31m"
	colorGreen  colorCode = "\033[32m"
	colorYellow colorCode = "\033[33m"
	colorBlue   colorCode = "\033[34m"
	colorPurple colorCode = "\033[35m"
	colorCyan   colorCode = "\033[36m"
	colorGray   colorCode = "\033[37m"
)

// Commander is the main CLI application handler.
// It manages categories of commands and provides help functionality.
type Commander struct {
	categories map[string]*Category
	args       []string
	output     io.Writer
}

// Category groups related commands under a common theme or functionality.
// For example, "Network" category might contain commands like "ping" and "curl".
type Category struct {
	Name     string
	commands map[string]*Command
}

// Command represents a single CLI command with its handler and metadata.
// Each command belongs to a Category and can accept structured arguments.
type Command struct {
	Name        string
	Description string
	Handler     interface{}
	flags       *flag.FlagSet
}

// HandlerFunc defines the type constraint for valid command handlers.
// Handlers must accept a context.Context as their first parameter and
// optionally a struct for arguments.
type HandlerFunc interface {
	~func(context.Context) | ~func(context.Context, any)
}

func (c *Commander) printf(format string, a ...any) {
	fmt.Fprintf(c.output, format, a...)
}

func (c *Commander) helpHandler(ctx context.Context, cmdName string) {
	if cmdName == "" {
		c.PrintUsage()
		return
	}

	cmd, err := c.findCommand(cmdName)
	if err != nil {
		c.printf("%s%s%s\n", colorRed, err.Error(), colorReset)
		c.PrintUsage()
		return
	}

	for catName, cat := range c.categories {
		if _, exists := cat.commands[cmdName]; exists {
			c.printf("\n%s%sHelp for command '%s%s%s' in category '%s%s%s':%s\n",
				colorBold, colorCyan,
				colorGreen, cmdName, colorReset,
				colorYellow, catName, colorReset,
				colorReset)
			c.printf("%s%sDescription:%s %s\n",
				colorBold, colorPurple, colorReset,
				cmd.Description)
			return
		}
	}
}

func (c *Commander) PrintUsage() {
	c.printf("%s%s🚀 Available Commands:%s\n",
		colorBold, colorCyan, colorReset)

	for _, cat := range c.categories {
		c.printf("\n%s📁 %s%s\n",
			colorYellow, cat.Name, colorReset)

		for name, cmd := range cat.commands {
			// Print command name and description
			c.printf("  %s%-12s%s %s\n",
				colorGreen, name,
				colorReset, cmd.Description)

			// Show flags if command has a struct argument
			handlerType := reflect.TypeOf(cmd.Handler)
			if handlerType.NumIn() > 1 {
				secondArg := handlerType.In(1)
				if secondArg.Kind() == reflect.Struct {
					// Print flags for each struct field
					for i := 0; i < secondArg.NumField(); i++ {
						field := secondArg.Field(i)
						flagName := field.Tag.Get(flagTag)
						if flagName == "" {
							flagName = strings.ToLower(field.Name)
						}
						usage := field.Tag.Get(usageTag)
						defaultValue := field.Tag.Get(defaultTag)

						// Format flag help
						c.printf("    %s--%s%s",
							colorCyan, flagName, colorReset)

						if field.Type.Kind() != reflect.Bool {
							c.printf(" <%s>", field.Type.String())
						}

						if usage != "" {
							c.printf("  %s", usage)
						}

						if defaultValue != "" {
							c.printf(" %s(default: %s)%s",
								colorYellow, defaultValue, colorReset)
						}
						c.printf("\n")
					}
				}
			}
		}
	}

	c.printf("\n%s%s💡 Usage:%s\n",
		colorBold, colorPurple, colorReset)
	c.printf("  %scommand [flags]%s\n",
		colorBlue, colorReset)
}

func isContextType(t reflect.Type) bool {
	return t.Implements(reflect.TypeOf((*context.Context)(nil)).Elem())
}

func (c *Commander) handleStructArgs(cmd *Command, fs *flag.FlagSet, args []reflect.Value) ([]reflect.Value, error) {
	handlerType := reflect.TypeOf(cmd.Handler)
	secondArg := handlerType.In(1)
	argValue := reflect.New(secondArg).Elem()

	// Register flags for struct fields
	for i := 0; i < secondArg.NumField(); i++ {
		field := secondArg.Field(i)
		flagName := strings.ToLower(field.Name)
		if tag := field.Tag.Get(flagTag); tag != "" {
			flagName = tag
		}

		usage := field.Tag.Get(usageTag)
		defaultValue := field.Tag.Get(defaultTag)

		fieldValue := argValue.Field(i)
		switch field.Type.Kind() {
		case reflect.Bool:
			defaultBool := defaultValue == "true"
			fs.BoolVar(fieldValue.Addr().Interface().(*bool), flagName, defaultBool, usage)
		case reflect.Int:
			var defaultInt int
			if defaultValue != "" {
				fmt.Sscanf(defaultValue, "%d", &defaultInt)
			}
			fs.IntVar(fieldValue.Addr().Interface().(*int), flagName, defaultInt, usage)
		case reflect.String:
			fs.StringVar(fieldValue.Addr().Interface().(*string), flagName, defaultValue, usage)
		}
	}

	if err := fs.Parse(c.args[2:]); err != nil {
		return nil, err
	}

	args[1] = argValue
	return args, nil
}

func (c *Commander) handleStringArg(args []reflect.Value) ([]reflect.Value, error) {
	var arg string
	if len(c.args) > 2 {
		arg = c.args[2]
	}
	args[1] = reflect.ValueOf(arg)
	return args, nil
}

func (c *Commander) callHandler(handler interface{}, args []reflect.Value) error {
	reflect.ValueOf(handler).Call(args)
	return nil
}

// New creates a new Commander instance with built-in help command.
// It uses os.Args for command-line arguments and os.Stdout for output.
func New() *Commander {
	return NewWithArgs(os.Args)
}

// NewWithArgs creates a new Commander instance with custom arguments.
// This is useful for testing or when you want to parse arguments from a different source.
func NewWithArgs(args []string) *Commander {
	cmdr := &Commander{
		categories: make(map[string]*Category),
		args:       args,
		output:     os.Stdout,
	}

	helpCat := cmdr.AddCategory("Help")
	helpCat.Register(&Command{
		Name:        "help",
		Description: "Show help information for commands",
		Handler:     cmdr.helpHandler,
	})

	return cmdr
}

// AddCategory creates a new command category with the given name.
// Categories help organize commands into logical groups.
// Returns a pointer to the new Category for method chaining.
func (c *Commander) AddCategory(name string) *Category {
	cat := &Category{
		Name:     name,
		commands: make(map[string]*Command),
	}
	c.categories[name] = cat
	return cat
}

// Register adds a command to a category.
// The command handler must follow the HandlerFunc interface constraints.
func (cat *Category) Register(cmd *Command) {
	if cat.commands == nil {
		cat.commands = make(map[string]*Command)
	}
	cat.commands[cmd.Name] = cmd
}

// Run executes the CLI application by parsing arguments and running the appropriate command.
// Returns an error if:
// - No subcommand is provided
// - The subcommand is unknown
// - The command handler is invalid
// - Flag parsing fails
func (c *Commander) Run() error {
	if len(c.args) < 2 {
		c.PrintUsage()
		return ErrNoSubcommand
	}

	cmd, err := c.findCommand(c.args[1])
	if err != nil {
		c.PrintUsage()
		return err
	}

	return c.executeCommand(cmd)
}

// findCommand locates a command by name across all categories
func (c *Commander) findCommand(name string) (*Command, error) {
	for _, cat := range c.categories {
		if cmd, exists := cat.commands[name]; exists {
			return cmd, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrUnknownCommand, name)
}

// executeCommand runs a single command with its arguments
func (c *Commander) executeCommand(cmd *Command) error {
	handlerType := reflect.TypeOf(cmd.Handler)
	if !isValidHandler(handlerType) {
		return ErrInvalidHandler
	}

	fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)
	cmd.flags = fs

	args, err := c.prepareArgs(cmd, fs)
	if err != nil {
		return err
	}

	return c.callHandler(cmd.Handler, args)
}

// isValidHandler checks if a function matches the HandlerFunc constraint
func isValidHandler(t reflect.Type) bool {
	return t.Kind() == reflect.Func &&
		t.NumIn() >= 1 &&
		isContextType(t.In(0))
}

// prepareArgs sets up command arguments and parses flags
func (c *Commander) prepareArgs(cmd *Command, fs *flag.FlagSet) ([]reflect.Value, error) {
	handlerType := reflect.TypeOf(cmd.Handler)
	args := make([]reflect.Value, handlerType.NumIn())
	args[0] = reflect.ValueOf(context.Background())

	if handlerType.NumIn() > 1 {
		secondArg := handlerType.In(1)
		switch secondArg.Kind() {
		case reflect.Struct:
			return c.handleStructArgs(cmd, fs, args)
		case reflect.String:
			return c.handleStringArg(args)
		default:
			return nil, fmt.Errorf("unsupported argument type: %v", secondArg)
		}
	}

	return args, nil
}

// SetOutput sets the writer where the commander will write its output.
// This is useful for testing or redirecting output to a different destination.
func (c *Commander) SetOutput(w io.Writer) {
	c.output = w
}
