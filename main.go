/*
*
In this example:

	We define a simple interface PluginInterface that specifies the method Run. Plugins must implement this interface.
	The main function loads a plugin specified by pluginName using plugin.Open.
	We then look up the Run symbol in the plugin using plugin.Lookup.
	We assert that the symbol implements the expected interface and cast it to the appropriate type.
	Finally, we call the function in the plugin with input parameters and handle any errors.

You'll need to compile your plugin separately as a shared object file
(.so on Unix/Linux systems) and replace "example_plugin.so" with the name of your plugin file. Additionally,
ensure that your plugin implements the PluginInterface interface and exports a function named Run that matches
the signature specified in the interface.

By George Oremo
For EMP Services Ltd
22 Fef 2024
*/
package main

import (
	"fmt"
	"plugin"
)

// Define a simple interface that plugins must implement.
type PluginInterface interface {
	Run(input string) (string, error)
}

func main() {
	// Name of the plugin to load
	pluginName := "example_plugin.so" // Replace with the name of your plugin file

	// Load the plugin
	p, err := plugin.Open(pluginName)
	if err != nil {
		fmt.Println("Error opening plugin:", err)
		return
	}

	// Look up the symbol (function) in the plugin
	runSymbol, err := p.Lookup("Run")
	if err != nil {
		fmt.Println("Error finding symbol in plugin:", err)
		return
	}

	// Assert that the symbol implements the PluginInterface
	var pluginFunc func(input string) (string, error)
	pluginFunc, ok := runSymbol.(func(input string) (string, error))
	if !ok {
		fmt.Println("Error: Symbol does not implement expected interface.")
		return
	}

	// Call the function in the plugin with input parameters
	result, err := pluginFunc("input_parameter")
	if err != nil {
		fmt.Println("Error calling plugin function:", err)
		return
	}

	fmt.Println("Plugin function returned:", result)
}
