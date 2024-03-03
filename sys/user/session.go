// Manage corpdesk user sessions
package user

import (
	"fmt"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
)

// Session model
type Session struct {
	SessionId     uint   `gorm:"primaryKey"`
	SessionGuid   string `json:"Session_guid"`
	SessionName   string `json:"Session_name"`
	ActivationKey string `json:"activation_key"`
	CompanyId     uint   `json:"company_id"`
	SessionTypeId uint   `json:"Session_type_id"`
}

func SessInit(cdToken string) {

	mErr := mc.Set(&memcache.Item{Key: "CD_SESS_ID", Value: []byte(cdToken)})
	if mErr != nil {
		fmt.Println("Error setting memache:", mErr)
		return
	}
	// Set the session ID as an environment variable
	err := os.Setenv("CD_SESS_ID", cdToken)
	if err != nil {
		fmt.Println("Error setting CD_SESS_ID environment variable:", err)
		return
	}
	fmt.Println("Session init success...cdToken:", cdToken)
}

func SessID() string {
	return "dummy_cd_token"
}

func SessIdGet() (string, error) {
	it, err := mc.Get("CD_SESS_ID")
	if err != nil {
		fmt.Println(err)
		return "", err
	} else {
		fmt.Printf("Key: %s, Value: %s\n", it.Key, it.Value)
		return string(it.Value), err
	}

	//////////////////////////////
	// return os.Getenv("CD_SESS_ID")
}

func SessIsValid() bool {
	sessId, err := SessIdGet()
	if err != nil {
		fmt.Println("Error getting session id:", err)
		return false
	}
	if sessId == "" {
		fmt.Println("session is invalid. Kindly login via 'cd-cli auth'")
		return false
	} else {
		return true
	}
}

func SessEnd() {
	// err := mc.Set(&memcache.Item{Key: "CD_SESS_ID", Value: []byte(cdToken)})
	// if err != nil {
	// 	fmt.Println("Error setting memache:", err)
	// 	return
	// }

	// Key of the item to delete
	// key := "CD_SESS_ID"

	if SessIsValid() {
		// // Delete the item from Memcached
		// err := mc.Delete("CD_SESS_ID")
		// if err != nil {
		// 	fmt.Println("Probelem encountered while ending session:", err)
		// }
		// Update: Set a new value for the existing key
		// key := "CD_SESS_ID"
		// newValue := []byte("")
		err := mc.Set(&memcache.Item{Key: "CD_SESS_ID", Value: []byte("")})
		if err != nil {
			fmt.Println("Error updating session cache.")
		}
		fmt.Println("Session ended. You can initiate and another session via cd-cli auth.")
	} else {
		fmt.Println("No valid session in progress.")
	}

}
