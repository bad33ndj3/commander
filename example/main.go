package main

import (
    "context"
    "fmt"
    "os"

    "github.com/bad33ndj3/commander"
)

// Command argument structs with tags
type StartEngineArgs struct {
    Quiet bool `flag:"quiet" default:"false" usage:"Start the engine in quiet mode"`
}

type ACArgs struct {
    Temperature int `flag:"temperature" default:"22" usage:"Temperature in Celsius"`
    FanSpeed    int `flag:"fanspeed" default:"3" usage:"Fan speed (1-5)"`
}

type HeatArgs struct {
    Temperature int `flag:"temperature" default:"20" usage:"Temperature in Celsius"`
}

func main() {
    cmdr := commander.New()

    engineCategory := cmdr.AddCategory("Engine")
    climateCategory := cmdr.AddCategory("Climate")
    infoCategory := cmdr.AddCategory("Information")

    // Engine controls
    engineCategory.Register(&commander.Command{
        Name:        "start",
        Description: "Starts the car engine",
        Handler:     startEngineHandler,
    })

    engineCategory.Register(&commander.Command{
        Name:        "stop",
        Description: "Stops the car engine",
        Handler:     stopEngineHandler,
    })

    // Climate controls
    climateCategory.Register(&commander.Command{
        Name:        "ac",
        Description: "Controls the air conditioning",
        Handler:     acHandler,
    })

    climateCategory.Register(&commander.Command{
        Name:        "heat",
        Description: "Controls the heating system",
        Handler:     heatHandler,
    })

    // Information commands
    infoCategory.Register(&commander.Command{
        Name:        "status",
        Description: "Displays the car's current status",
        Handler:     statusHandler,
    })

    infoCategory.Register(&commander.Command{
        Name:        "fuel",
        Description: "Shows fuel level",
        Handler:     fuelHandler,
    })

    if err := cmdr.Run(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

// Handler functions using argument structs
func startEngineHandler(ctx context.Context, args StartEngineArgs) {
    if args.Quiet {
        fmt.Println("Quietly starting the car engine...")
    } else {
        fmt.Println("Starting the car engine... VROOM!")
    }
}

func stopEngineHandler(ctx context.Context) {
    fmt.Println("Stopping the car engine...")
}

func acHandler(ctx context.Context, args ACArgs) {
    fmt.Printf("Setting AC temperature to %d°C with fan speed %d\n", args.Temperature, args.FanSpeed)
}

func heatHandler(ctx context.Context, args HeatArgs) {
    fmt.Printf("Setting heating temperature to %d°C\n", args.Temperature)
}

func statusHandler(ctx context.Context) {
    fmt.Println("Car Status:")
    fmt.Println("- Engine: Running")
    fmt.Println("- Speed: 0 km/h")
    fmt.Println("- Temperature: 22°C")
}

func fuelHandler(ctx context.Context) {
    fmt.Println("Fuel level: 75%")
}