/*
*
Entry point for cd system and applications

By George Oremo
For EMP Services Ltd
22 Fef 2024
*/
package cd

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/tcp-x/cd-core/sys/base"
)

var mc = memcache.New("localhost:11211")

func exec(data string) {
	// base.Exec(data) // Call the function with the parameter
}

func run(data base.ICdRequest) {

}
