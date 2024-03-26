// Manage corpdesk user sessions
package user

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
	"github.com/tcp-x/cd-core/sys/base"
)

var logIndex = 0

// Session model
// type Session struct {
// 	SessionId     uint `gorm:"primaryKey"`
// 	CurrentUserId uint
// 	CdToken       string
// 	Active        bool
// 	Ttl           uint
// 	AccTime       time.Time
// 	StartTime     time.Time
// 	DeviceNetId   datatypes.JSON
// 	ConsumerGuid  string
// }

type Session struct {
	SessionId     uint      `gorm:"session_id,primaryKey"`
	CurrentUserId uint      `gorm:"current_user_id"`
	CdToken       string    `gorm:"cd_token"`
	Active        bool      `gorm:"active"`
	Ttl           uint      `gorm:"ttl"`
	ConsumerGuid  string    `gorm:"consumer_guid"`
	AccTime       time.Time `gorm:"acc_time"`
}

func CreateSess(req base.CdRequest) base.CdResponse {
	logger.LogInfo(strconv.Itoa(logIndex) + ". Starting UserModule::Session::SessNew()...")
	// var sess Session
	// sess.AccTime = time.Now()
	// sess.StartTime = time.Now()
	// // sessResult := db.Create(&sess)
	// sessResult := db.Table("session").Create(&sess)
	// if sessResult.Error != nil {
	// 	fmt.Println("Error creating session:", sessResult.Error)
	// 	return 0, sessResult.Error
	// }
	// logger.LogInfo("UserModule::Session::SessCreate()/result:" + fmt.Sprint(sessResult))

	// Create a new session instance
	cdToken := uuid.New()
	sess := Session{
		CdToken:       cdToken.String(),
		Active:        true,
		ConsumerGuid:  "",
		CurrentUserId: 1000,
		AccTime:       time.Now(),
	}
	servInput := base.ServiceInput{}
	servInput.ModelName = "session"
	servInput.ServiceModel = sess
	resp := base.CdCreate(req, servInput)
	req.Dat.Token = sess.CdToken
	// sessReq := base.CdRequest{}
	// sessReq.Ctx = "sys"
	// sessReq.M = "user"
	// sessReq.C = "session"
	// sessReq.A = "GetSessByToken"
	// sessDat := sessReq.Dat.F_vals.Data.(Session)
	// sessDat.CdToken = sess.CdToken
	// sessReq.Dat.Token = sess.CdToken
	sessReq := SetSessionRequest(sess, "GetSessByToken")
	sessReq.Dat.Token = fmt.Sprint(cdToken)
	tokenResp := GetSessByToken(sessReq)
	logger.LogInfo("UserModule::Session::CreateSess()/tokenResp:" + fmt.Sprint(tokenResp))
	logger.LogInfo("UserModule::Session::CreateSess()/tokenResp.Data:" + fmt.Sprint(tokenResp.Data))
	sessData, err := json.Marshal(tokenResp.Data)
	if err != nil {
		logger.LogInfo("UserModule::User::Auth()/Error creating sesson:" + err.Error())
		var appState = base.CdAppState{false, err.Error(), nil, "", ""}
		var appData = base.RespData{Data: []User{anon}, RowsAffected: 0, NumberOfResult: 1}
		resp := base.CdResponse{AppState: appState, Data: appData}
		return resp
	}
	// logger.LogInfo("UserModule::Session::CreateSess()/sessData:" + fmt.Sprint(sessData))
	resp.AppState.Sess = sessData
	resp.Data = tokenResp.Data
	return resp
}

func SetSessionRequest(newUser Session, action string) base.CdRequest {
	// newUser := User{"karl", "secret", "karl@emp.net"}
	fvalItem := base.FValItem{newUser}
	fvalDat := base.FValDat{fvalItem, ""}
	return base.CdRequest{"sys", "user", "session", action, fvalDat}
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
	// create a new session
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

/*
`active` tinyint(1) DEFAULT NULL,

	`ttl` int DEFAULT NULL,
	`acc_time` datetime(4) DEFAULT NULL,
	`start_time` datetime DEFAULT NULL,
	`device_net_id` json DEFAULT NULL,
	`consumer_guid` varchar(40) DEFAULT NULL,
*/
func GetSessByToken(req base.CdRequest) base.CdResponse {
	logger.LogInfo(strconv.Itoa(logIndex) + ". Starting UserModule::Session::GetSessByToken()...")
	var session Session
	sessionResult := db.Table("session").
		Select("session_id", "current_user_id", "cd_token", "active", "acc_time", "start_time", "consumer_guid").
		Where("cd_token = ?", req.Dat.Token).
		Scan(&session)
	if sessionResult.Error != nil {
		logger.LogInfo("UserModule::Session::GetSessByToken()/fetching sesson data:" + fmt.Sprint(sessionResult.Error))
		var appState = base.CdAppState{false, fmt.Sprint(sessionResult.Error), nil, "", ""}
		var appData = base.RespData{Data: []User{anon}, RowsAffected: 0, NumberOfResult: 1}
		resp := base.CdResponse{AppState: appState, Data: appData}
		return resp
	}
	logger.LogInfo("UserModule::Session::GetSessByToken()/sessionResult:" + fmt.Sprint(sessionResult))
	// rows, err := userResult.Rows()
	// if err != nil {
	// 	logger.LogInfo("UserModule::Session::GetSessByToken()/scanning sesson data:" + fmt.Sprint(sessionResult.Error))
	// 	var appState = base.CdAppState{false, fmt.Sprint(sessionResult.Error), nil, "", ""}
	// 	var appData = base.RespData{Data: []User{anon}, RowsAffected: 0, NumberOfResult: 1}
	// 	resp := base.CdResponse{AppState: appState, Data: appData}
	// 	return resp
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	db.ScanRows(rows, &session)
	// 	fmt.Println(session)
	// 	// if users.UserName == username {
	// 	// 	u = users
	// 	// }
	// 	// if users.UserName == "anon" {
	// 	// 	anon = users
	// 	// }
	// }

	var appState = base.CdAppState{true, "", nil, "", ""}
	// appState.Sess = sid
	var appData = base.RespData{Data: []Session{session}, RowsAffected: 0, NumberOfResult: 1}
	resp := base.CdResponse{AppState: appState, Data: appData}
	return resp
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
