/*
*
Entry point for cd system and applications

By George Oremo
For EMP Services Ltd
22 Fef 2024
*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/tcp-x/cd-core/sys/base"
	"github.com/tcp-x/cd-core/sys/user"
)

var mc = memcache.New("localhost:11211")
var logger base.Logger
var logIndex = 0

func exec(data string) {
	// base.Exec(data) // Call the function with the parameter
}

func run(data base.ICdRequest) {

}

func setCdRequest(newUser user.User) user.CdRequest {
	// newUser := User{"karl", "secret", "karl@emp.net"}
	fvalItem := user.FValItem{newUser}
	fvalDat := user.FValDat{fvalItem, ""}
	return user.CdRequest{"Sys", "UserModule", "UserController", "Create", fvalDat}
}

func main() {
	// userHandle := new(UserController)
	var err error
	var cdResp user.CdResponse

	// -------------------------------
	// AUTH USER
	// -------------------------------
	req := setCdRequest(user.User{0, "", "karl", "secretx", "", 0, "", false, time.Now(), "", "", "", "", 0, 0, false, "", "", 0, 0})
	// err = userHandle.EditPassword(req6, &cdResp)
	cdResp = user.Auth(req)
	if err != nil {
		// log.Fatal("Issue authenticating User: ", err)
		logger.LogInfo(strconv.Itoa(logIndex) + ". Response: " + string(err.Error()))
	}
	logIndex++
	log.Println(strconv.Itoa(logIndex)+".", cdResp)

	respJStr, err := json.Marshal(cdResp)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(string(respJStr))
	// var r string = string(respJStr)
	logIndex++
	logger.LogInfo(strconv.Itoa(logIndex) + ". Response: " + string(respJStr))

}
