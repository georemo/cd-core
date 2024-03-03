// corpdesk base module
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

/*
ExecPlug is for processing corpdesk json command
Input: json string request in ICdRequest format
Output 1: json output in ICdResponse format
Output 2: error
*/
func ExecPlug(req string) (string, error) {

	fmt.Println("b::ExecPlug()/Processing JSON...")

	err := json.Unmarshal([]byte(req), &jsonMap)
	if err == nil {
		fmt.Println("Successfull JSON encoding")
		fmt.Println(jsonMap)

		jReq.ctx = removeQt(jToStr("ctx"))
		jReq.m = removeQt(jToStr("m"))
		jReq.c = removeQt(jToStr("c"))
		jReq.a = removeQt(jToStr("a"))
		jReq.dat = removeQt(jToStr("dat"))

	} else {
		fmt.Println("Error:", err)
		return "", err
	}

	/////////////////////////////////////
	// Name of the plugin to load
	fmt.Println("Controller:", jReq.c)
	pluginName := "plugins/" + jReq.m + "/" + jReq.c + ".so" // Replace with the name of your plugin file
	fmt.Println("pluginName:", pluginName)

	fmt.Printf("Loading plugin %s", pluginName)
	p, err := plugin.Open(pluginName)
	if err != nil {
		// panic(err)
		return "", err

	}

	////////////////////////////////////
	// Lookup – Searches for Action symbol name in the plugin
	symbolAx, errAuth := p.Lookup(jReq.a)
	if errAuth != nil {
		// panic(errAuth)
		return "", errAuth

	}

	// symbol – Checks the function signature
	f, ok := symbolAx.(func(string) string)
	if !ok {
		// panic("Plugin has no " + jReq.a + " 'f(string)string' function")
		fmt.Printf("\nPlugin has no %s  'f(string)string' function\n", jReq.a)
		return "", &CdError{}
	}

	// Uses f() function to return results
	resp := f(jReq.dat)
	fmt.Printf("\nf() return is:%s\n", resp)

	return resp, nil
}

func Run(req string) string {

	Conn()

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
