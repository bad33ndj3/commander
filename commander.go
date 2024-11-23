package commander

import (
	"context"
	"errors"
	"flag"
	"fmt"
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

// Commander is the main CLI application handler
type Commander struct {
	categories map[string]*Category
}

// Category groups related commands
type Category struct {
	Name     string
	commands map[string]*Command
}

// Command represents a single CLI command
type Command struct {
	Name        string
	Description string
	Handler     interface{}
	flags       *flag.FlagSet
}

// HandlerFunc is a type constraint for command handlers
type HandlerFunc interface {
	~func(context.Context) | ~func(context.Context, any)
}

func (c *Commander) helpHandler(ctx context.Context, cmdName string) {
	if cmdName == "" {
		c.PrintUsage()
		return
	}

	cmd, err := c.findCommand(cmdName)
	if err != nil {
		fmt.Printf("%s%s%s\n", colorRed, err.Error(), colorReset)
		c.PrintUsage()
		return
	}

	for catName, cat := range c.categories {
		if _, exists := cat.commands[cmdName]; exists {
			fmt.Printf("\n%s%sHelp for command '%s%s%s' in category '%s%s%s':%s\n",
				colorBold, colorCyan,
				colorGreen, cmdName, colorReset,
				colorYellow, catName, colorReset,
				colorReset)
			fmt.Printf("%s%sDescription:%s %s\n",
				colorBold, colorPurple, colorReset,
				cmd.Description)
			return
		}
	}
}

func (c *Commander) PrintUsage() {
	fmt.Printf("%s%s🚀 Available Commands:%s\n",
		colorBold, colorCyan, colorReset)

	for _, cat := range c.categories {
		fmt.Printf("\n%s📁 %s%s\n",
			colorYellow, cat.Name, colorReset)

		for name, cmd := range cat.commands {
			fmt.Printf("  %s%-12s%s %s\n",
				colorGreen, name,
				colorReset, cmd.Description)
		}
	}

	fmt.Printf("\n%s%s💡 Usage:%s\n",
		colorBold, colorPurple, colorReset)
	fmt.Printf("  %scommand [flags]%s\n",
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

	if err := fs.Parse(os.Args[2:]); err != nil {
		return nil, err
	}

	args[1] = argValue
	return args, nil
}

func (c *Commander) handleStringArg(args []reflect.Value) ([]reflect.Value, error) {
	var arg string
	if len(os.Args) > 2 {
		arg = os.Args[2]
	}
	args[1] = reflect.ValueOf(arg)
	return args, nil
}

func (c *Commander) callHandler(handler interface{}, args []reflect.Value) error {
	reflect.ValueOf(handler).Call(args)
	return nil
}

// New creates a new Commander instance with built-in help command
func New() *Commander {
	cmdr := &Commander{
		categories: make(map[string]*Category),
	}
	
	helpCat := cmdr.AddCategory("Help")
	helpCat.Register(&Command{
		Name:        "help",
		Description: "Show help information for commands",
		Handler:     cmdr.helpHandler,
	})
	
	return cmdr
}

// AddCategory creates a new command category
func (c *Commander) AddCategory(name string) *Category {
	cat := &Category{
		Name:     name,
		commands: make(map[string]*Command),
	}
	c.categories[name] = cat
	return cat
}

// Register adds a command to a category
func (cat *Category) Register(cmd *Command) {
	if cat.commands == nil {
		cat.commands = make(map[string]*Command)
	}
	cat.commands[cmd.Name] = cmd
}

// Run executes the CLI application
func (c *Commander) Run() error {
	if len(os.Args) < 2 {
		c.PrintUsage()
		return ErrNoSubcommand
	}

	cmd, err := c.findCommand(os.Args[1])
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

