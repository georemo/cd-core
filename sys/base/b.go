package base

import (
	"encoding/json"
	"fmt"
	"os"
	"plugin"
)

var jsonMap map[string]interface{}
var jReq ICdRequest

func jToStr(field string) string {
	f := jsonMap[field]
	// fmt.Println("ctx:", ctx)
	biteF, err := json.Marshal(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ""
	} else {
		// fmt.Println("biteF:", biteF)
		return string(biteF[:])
	}

}

func removeQt(s string) string {
	return s[1 : len(s)-1]
}

func Run(req string) string {

	fmt.Println("b::Run()/Processing JSON...")

	r := json.Unmarshal([]byte(req), &jsonMap)
	if r == nil {
		fmt.Println("Successfull JSON encoding")
		fmt.Println(jsonMap)

		jReq.ctx = removeQt(jToStr("ctx"))
		jReq.m = removeQt(jToStr("m"))
		jReq.c = removeQt(jToStr("c"))
		jReq.a = removeQt(jToStr("a"))
		jReq.dat = removeQt(jToStr("dat"))

	} else {
		fmt.Println("Error:", r)
	}

	////////////////////////////////////
	// Name of the plugin to load
	fmt.Println("Controller:", jReq.c)
	pluginName := jReq.c + ".so" // Replace with the name of your plugin file
	fmt.Println("pluginName:", pluginName)

	// Load the plugin
	p, err := plugin.Open(pluginName)
	if err != nil {
		fmt.Println("Error opening plugin:", err)
		return "{}"
	}

	// Look up the symbol (function) in the plugin
	runSymbol, err := p.Lookup(jReq.a)
	if err != nil {
		fmt.Println("Error finding symbol in plugin:", err)
		return "{}"
	}

	// Assert that the symbol implements the PluginInterface
	var pluginFunc func(string) (string, error)
	pluginFunc, ok := runSymbol.(func(string) (string, error))
	if !ok {
		fmt.Println("Error: Symbol does not implement expected interface.")
		return "{}"
	}

	// Call the function in the plugin with input parameters
	resp, err := pluginFunc(jReq.dat)
	if err != nil {
		fmt.Println("Error calling plugin function:", err)
		return "{}"
	}

	fmt.Println("Plugin function returned:", resp)
	return resp
}
