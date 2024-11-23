package commander

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

// Add new struct tag constants
const (
	flagTag     = "flag"
	defaultTag  = "default"
	usageTag    = "usage"
)

// Add color codes
const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
)

func New() *Commander {
	cmdr := &Commander{
		Categories: make(map[string]*Category),
	}
	
	// Add built-in help category and command
	helpCat := cmdr.AddCategory("Help")
	helpCat.Register(&Command{
		Name:        "help",
		Description: "Show help information for commands",
		Handler:     cmdr.helpHandler,
	})
	
	return cmdr
}

type Commander struct {
	Categories map[string]*Category
}

// Helper function for colored text
func colored(color string, text string) string {
	return color + text + colorReset
}

// Helper function for bold colored text
func coloredBold(color string, text string) string {
	return color + colorBold + text + colorReset
}

func (c *Commander) helpHandler(ctx context.Context, command string) {
	if command == "" {
		c.PrintUsage()
		return
	}

	// Search for specific command help
	for _, category := range c.Categories {
		if cmd, exists := category.Commands[command]; exists {
			fmt.Printf("\n%s '%s' %s '%s':\n",
				coloredBold(colorCyan, "Help for command"),
				colored(colorGreen, command),
				colored(colorGray, "in category"),
				colored(colorYellow, category.Name))
			fmt.Printf("%s %s\n",
				coloredBold(colorPurple, "Description:"),
				colored(colorGray, cmd.Description))
			
			if handlerType := reflect.TypeOf(cmd.Handler); handlerType.NumIn() > 1 {
				fmt.Printf("\n%s\n", coloredBold(colorBlue, "Flags:"))
				c.printFlagUsage(cmd.Handler)
			}
			return
		}
	}
	
	fmt.Printf("%s: %s\n",
		colored(colorRed, "Unknown command"),
		colored(colorYellow, command))
	c.PrintUsage()
}

func (c *Commander) registerFlags(fs *flag.FlagSet, handler interface{}) {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		return
	}

	for i := 1; i < handlerType.NumIn(); i++ {
		paramType := handlerType.In(i)
		paramName := strings.ToLower(paramType.Name())
		
		// Skip if parameter name is empty
		if paramName == "" {
			paramName = fmt.Sprintf("param%d", i)
		}

		switch paramType.Kind() {
		case reflect.Bool:
			fs.Bool(paramName, false, fmt.Sprintf("Flag for %s", paramName))
		case reflect.Int:
			fs.Int(paramName, 0, fmt.Sprintf("Flag for %s", paramName))
		case reflect.String:
			fs.String(paramName, "", fmt.Sprintf("Flag for %s", paramName))
		// Add other types as needed
		}
	}
}

func (c *Commander) printFlagUsage(handler interface{}) {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func || handlerType.NumIn() <= 1 {
		return
	}

	argsType := handlerType.In(1)
	if argsType.Kind() != reflect.Struct {
		return
	}

	// Get max flag name length for alignment
	maxFlagLen := 0
	for i := 0; i < argsType.NumField(); i++ {
		field := argsType.Field(i)
		flagName := field.Tag.Get(flagTag)
		if flagName == "" {
			flagName = strings.ToLower(field.Name)
		}
		if len(flagName) > maxFlagLen {
			maxFlagLen = len(flagName)
		}
	}

	// Print flags
	for i := 0; i < argsType.NumField(); i++ {
		field := argsType.Field(i)
		
		flagName := field.Tag.Get(flagTag)
		if flagName == "" {
			flagName = strings.ToLower(field.Name)
		}
		usage := field.Tag.Get(usageTag)
		defaultValue := field.Tag.Get(defaultTag)

		var typeHint string
		switch field.Type.Kind() {
		case reflect.Bool:
			typeHint = "bool"
		case reflect.Int:
			typeHint = "number"
		case reflect.String:
			typeHint = "string"
		default:
			typeHint = field.Type.String()
		}

		// Print flag with aligned columns and colors
		padding := strings.Repeat(" ", maxFlagLen-len(flagName))
		fmt.Printf("    %s%s", 
			colored(colorCyan, "--"+flagName),
			padding)
		
		// Print type hint if not bool
		if typeHint != "bool" {
			fmt.Printf(" %s", colored(colorPurple, "<"+typeHint+">"))
		}
		
		// Print usage and default value on same line
		if defaultValue != "" {
			fmt.Printf("  %s %s\n",
				colored(colorGray, usage),
				colored(colorYellow, "(default: "+defaultValue+")"))
		} else {
			fmt.Printf("  %s\n", colored(colorGray, usage))
		}
	}
}

func (c *Commander) PrintUsage() {
	fmt.Printf("%s\n", coloredBold(colorCyan, "🚀 Available Commands:"))
	fmt.Println(colored(colorGray, "=================="))

	// Get max command name length for alignment
	maxNameLen := 0
	for _, category := range c.Categories {
		for name := range category.Commands {
			if len(name) > maxNameLen {
				maxNameLen = len(name)
			}
		}
	}

	// Print each category and its commands
	for _, category := range c.Categories {
		fmt.Printf("%s %s\n", colored(colorYellow, "📁"), coloredBold(colorYellow, category.Name))
		
		for name, cmd := range category.Commands {
			// Print command name and description aligned
			padding := strings.Repeat(" ", maxNameLen-len(name))
			fmt.Printf("  %s%s  %s\n", 
				colored(colorGreen, name),
				padding,
				colored(colorGray, cmd.Description))
			
			// Print flags if the command has any
			if handlerType := reflect.TypeOf(cmd.Handler); handlerType.NumIn() > 1 {
				c.printFlagUsage(cmd.Handler)
			}
		}
		fmt.Println() // Single line between categories
	}

	fmt.Printf("%s\n", coloredBold(colorPurple, "💡 Usage:"))
	fmt.Printf("  %s\n", colored(colorBlue, "command [flags]"))
	fmt.Printf("  %s  %s\n", 
		colored(colorBlue, "help <command>"),
		colored(colorGray, "Show help for command"))
}

func (c *Commander) AddCategory(name string) *Category {
	category := &Category{
		Name:     name,
		
		Commands: make(map[string]*Command),
	}
	c.Categories[name] = category
	return category
}

type Category struct {
	Name     string
	Commands map[string]*Command
}

func (cat *Category) Register(cmd *Command) {
	if cat.Commands == nil {
		cat.Commands = make(map[string]*Command)
	}
	cat.Commands[cmd.Name] = cmd
}

type Command struct {
	Name        string
	Description string
	Handler     interface{}
	Flags       *flag.FlagSet
}

func (c *Commander) Run() error {
	if len(os.Args) < 2 {
		c.PrintUsage()
		return errors.New("no subcommand provided")
	}

	subcommandName := os.Args[1]
	var cmd *Command
	var found bool

	// Search for the command in all categories
	for _, category := range c.Categories {
		if category.Commands[subcommandName] != nil {
			cmd = category.Commands[subcommandName]
			found = true
			break
		}
	}

	if !found {
		c.PrintUsage()
		return fmt.Errorf("unknown subcommand: %s", subcommandName)
	}

	// Create a new FlagSet for this command
	fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)
	cmd.Flags = fs

	handlerType := reflect.TypeOf(cmd.Handler)
	if handlerType.Kind() != reflect.Func {
		return errors.New("handler is not a function")
	}

	if handlerType.NumIn() < 1 {
		return errors.New("handler function must accept at least one argument (context.Context)")
	}

	// Check that the first parameter is context.Context
	if !isContextType(handlerType.In(0)) {
		return errors.New("first parameter must be context.Context")
	}

	// Prepare argument values
	argValues := make([]reflect.Value, handlerType.NumIn())
	ctx := context.Background()
	argValues[0] = reflect.ValueOf(ctx)

	// Handle arguments based on the second parameter type
	if handlerType.NumIn() > 1 {
		secondParamType := handlerType.In(1)
		
		switch secondParamType.Kind() {
		case reflect.Struct:
			// Handle struct arguments
			argsValue := reflect.New(secondParamType)
			c.registerStructFlags(fs, argsValue.Elem())
			
			if err := fs.Parse(os.Args[2:]); err != nil {
				return err
			}
			
			argValues[1] = argsValue.Elem()
			
		case reflect.String:
			// Handle string argument (for help command)
			helpArg := ""
			if len(os.Args) > 2 {
				helpArg = os.Args[2]
			}
			argValues[1] = reflect.ValueOf(helpArg)
			
		default:
			return fmt.Errorf("unsupported parameter type: %v", secondParamType)
		}
	}

	// Call the handler
	reflect.ValueOf(cmd.Handler).Call(argValues)
	return nil
}

func (c *Commander) registerStructFlags(fs *flag.FlagSet, structValue reflect.Value) {
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)
		
		// Get flag name from tag or use field name
		flagName := field.Tag.Get(flagTag)
		if flagName == "" {
			flagName = strings.ToLower(field.Name)
		}

		// Get usage from tag
		usage := field.Tag.Get(usageTag)
		if usage == "" {
			usage = fmt.Sprintf("Flag for %s", field.Name)
		}

		// Get default value from tag
		defaultValue := field.Tag.Get(defaultTag)

		switch field.Type.Kind() {
		case reflect.Bool:
			defaultBool := false
			if defaultValue != "" {
				defaultBool = defaultValue == "true"
			}
			fs.BoolVar(fieldValue.Addr().Interface().(*bool), flagName, defaultBool, usage)

		case reflect.Int:
			defaultInt := 0
			if defaultValue != "" {
				if parsed, err := strconv.Atoi(defaultValue); err == nil {
					defaultInt = parsed
				}
			}
			fs.IntVar(fieldValue.Addr().Interface().(*int), flagName, defaultInt, usage)

		case reflect.String:
			fs.StringVar(fieldValue.Addr().Interface().(*string), flagName, defaultValue, usage)

		// Add other types as needed
		}
	}
}

// Helper functions
func isContextType(t reflect.Type) bool {
	return t.Implements(reflect.TypeOf((*context.Context)(nil)).Elem())
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float64:
		return v.Float() == 0.0
	case reflect.String:
		return v.String() == ""
	default:
		return false
	}
}

// Helper function to extract parameter names from function name
func getFunctionParamNames(handler interface{}) []string {
	// Get function pointer
	ptr := reflect.ValueOf(handler).Pointer()
	fn := runtime.FuncForPC(ptr)
	if fn == nil {
		return nil
	}

	// Get the function type
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		return nil
	}

	// Skip the context parameter (index 0) and collect remaining parameter names
	paramNames := make([]string, 0, handlerType.NumIn()-1)
	for i := 1; i < handlerType.NumIn(); i++ {
		paramType := handlerType.In(i)
		paramName := strings.ToLower(paramType.Name()) // This will be empty for built-in types

		// If it's a built-in type, get the parameter name from the type
		if paramName == "" {
			paramName = strings.ToLower(paramType.String())
		}

		paramNames = append(paramNames, paramName)
	}

	return paramNames
}
