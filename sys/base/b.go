// corpdesk base module
package base

import (
	"encoding/json"
	"fmt"
	"os"
	"plugin"

	"github.com/bradfitz/gomemcache/memcache"
	"gorm.io/datatypes"
)

var db = Conn()
var mc = memcache.New("localhost:11211")
var logger Logger
var respMsg = ""

// ///////////////////////////////////////
// Make a new CdResponse type that is a typed collection of fields
// (Title and Status), both of which are of type string
type CdResponse struct {
	AppState CdAppState
	Data     RespData
}

type CdAppState struct {
	Success bool
	Info    string
	Sess    datatypes.JSON
	Cache   string
	SConfig string
}

type RespData struct {
	Data           interface{}
	RowsAffected   int
	NumberOfResult int
}

type CdRequest struct {
	Ctx string
	M   string
	C   string
	A   string
	Dat FValDat
}

type FValDat struct {
	F_vals FValItem
	Token  string
}

type FValItem struct {
	Data interface{}
}

type ServiceInput struct {
	ServiceModel interface{}
	ModelName    string
	DocName      string
	Cmd          Cmd
	DSource      int
}

type Cmd struct {
	Action string
	Query  json.RawMessage
}

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

	// Assert that the symbol implements the PluginInterface--
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

// Function to convert JSON string to ICdRequest
func JSONToICdRequest(jsonString string) (ICdRequestExp, error) {
	var reqData ICdRequestExp
	err := json.Unmarshal([]byte(jsonString), &reqData)
	return reqData, err
}

func CdCreate(req CdRequest, servInput ServiceInput) CdResponse {
	result := db.Table(servInput.ModelName).Create(&servInput.ServiceModel)
	if result.Error != nil {
		respMsg = "Could not create" + servInput.ModelName
		logger.LogInfo("Base::CbCreate()/respMsg:" + respMsg)
		logger.LogError("Base::CbCreate():" + fmt.Sprint(result.Error))
		var appState = CdAppState{false, respMsg, nil, "", ""}
		var appData = RespData{Data: nil, RowsAffected: 0, NumberOfResult: 0}
		resp := CdResponse{AppState: appState, Data: appData}
		return resp
	}
	// Convert query result to JSON
	r, err := result.Rows()
	if err != nil {
		var appState = CdAppState{false, respMsg, nil, "", ""}
		var appData = RespData{Data: nil, RowsAffected: 0, NumberOfResult: 0}
		resp := CdResponse{AppState: appState, Data: appData}
		return resp
	}
	jsonResult, err := json.Marshal(r)
	if err != nil {
		var appState = CdAppState{false, respMsg, nil, "", ""}
		var appData = RespData{Data: nil, RowsAffected: 0, NumberOfResult: 0}
		resp := CdResponse{AppState: appState, Data: appData}
		return resp
	}
	resp := CdResponse{}
	var appState = CdAppState{true, respMsg, jsonResult, "", ""}
	var appData = RespData{Data: jsonResult, RowsAffected: 0, NumberOfResult: 1}
	resp = CdResponse{AppState: appState, Data: appData}
	return resp
}
