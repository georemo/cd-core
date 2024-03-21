package user

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/tcp-x/cd-core/sys/base"
	"golang.org/x/crypto/bcrypt"
)

var db = base.Conn()
var mc = memcache.New("localhost:11211")
var logger base.Logger

type IUser interface{}
type FVals struct {
	data User
}

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
	Sess    string
	Cache   string
	SConfig string
}

type RespData struct {
	Data           []User
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
	Data User
}

////////////////////////////////////////////

// User model
type User struct {
	UserId        uint      `gorm:"primaryKey"`
	UserGuid      string    `json:"user_guid"`
	UserName      string    `gorm: "column: user_name`
	Password      string    `json:"password"`
	Email         string    `json:"email"`
	DocId         uint      `json:"doc_id"`
	Mobile        string    `json:"mobile"`
	Gender        bool      `json:"gender"`
	BirthDate     time.Time `json:"birth_date"`
	PostalAddr    string    `json:"Postal_addr"`
	FName         string    `json:"f_name"`
	MName         string    `json:"m_name"`
	LName         string    `json:"l_name"`
	NationalId    uint      `json:"natonal_id"`
	PassportId    uint      `json:"passport_id"`
	UserEnabled   bool      `json:"user_enabled"`
	ZipCode       string    `json:"zip_code"`
	ActivationKey string    `json:"activation_key"`
	CompanyId     uint      `json:"company_id"`
	UserTypeId    uint      `json:"user_type_id"`
}

/*
*
  - {
    "ctx": "Sys",
    "m": "User",
    "c": "User",
    "a": "Login",
    "dat": {
    "f_vals": [
    {
    "data": {
    "userName": "jondoo",
    "password": "iiii",
    "consumerGuid": "B0B3DA99-1859-A499-90F6-1E3F69575DCD"
    }
    }
    ],
    "token": ""
    },
    "args": null
    }
  - @param req
  - @param res
*/
func Auth(req CdRequest) CdResponse {
	var users []User
	// var records []User
	// usr := base.Get("user", records, db)

	// get user and anon data
	// 1. convert req to struct
	// reqStruct, err := base.JSONToICdRequest(req)
	// if err != nil {
	// 	logger.LogError(err.Error())
	// 	return "", nil
	// }

	// // Accessing fields of MyStruct
	// logger.LogInfo("Module:" + reqStruct.M)
	// logger.LogInfo("Dat:" + reqStruct.Dat)

	// fv, err := fVals(reqStruct.Dat)
	// if err != nil {
	// 	log.Fatal("Error:", err)
	// 	return "", nil
	// }

	authenticated, err := AuthenticateUser(req.Dat.F_vals.Data.UserName, req.Dat.F_vals.Data.UserName)
	if err != nil {
		log.Fatal("Error authenticating user:", err)
		var appState = CdAppState{true, err.Error(), "", "", ""}
		var appData = RespData{Data: users, RowsAffected: 0, NumberOfResult: 1}
		resp := CdResponse{AppState: appState, Data: appData}
		return resp
	}

	respMsg := ""
	if authenticated {
		respMsg = "User authenticated successfully"
		fmt.Println(respMsg)
	} else {
		respMsg = "User authenticated successfully"
		fmt.Println(respMsg)
	}

	// connect to db and check validity of password
	// Auth input should have username and password

	// test if /tcp-x/user/session is accessible
	sid := SessID()
	fmt.Println("cd-user/Auth(): SessionID:", sid)

	// resp := "{name:User, version:0.0.7 publisher: \"EMP Services Ltd\"}"
	var appState = CdAppState{authenticated, respMsg, "", "", ""}
	var appData = RespData{Data: users, RowsAffected: 0, NumberOfResult: 1}
	resp := CdResponse{AppState: appState, Data: appData}
	return resp
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a hashed password with its plaintext version
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// AuthenticateUser authenticates a user by username and password
func AuthenticateUser(username, password string) (bool, error) {
	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return false, result.Error
	}
	return CheckPasswordHash(password, user.Password), nil
}

func fVals(fvals string) (FVals, error) {
	fvalStruct, err := JSONToFVals(fvals)
	if err != nil {
		fmt.Println("Error:", err)
		return fvalStruct, err
	}
	return fvalStruct, nil
}

// Function to convert JSON string to fVals
func JSONToFVals(jsonString string) (FVals, error) {
	var fVals FVals
	err := json.Unmarshal([]byte(jsonString), &fVals)
	return fVals, err
}
