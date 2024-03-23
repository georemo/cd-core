package user

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
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
	UserName      string    `json:"username"`
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
	logger.LogInfo("Module version:v0.0.56")
	logger.LogInfo("Starting UserModule::User::Auth()...")
	var users []User

	logger.LogInfo("UserModule::User::Auth()/req.Dat.F_vals.Data.UserName:" + req.Dat.F_vals.Data.UserName)
	logger.LogInfo("UserModule::User::Auth()/req.Dat.F_vals.Data.Password:" + req.Dat.F_vals.Data.Password)
	authenticated, err := AuthenticateUser(req.Dat.F_vals.Data.UserName, req.Dat.F_vals.Data.Password)
	logger.LogInfo("UserModule::User::Auth()/authenticated:" + fmt.Sprint(authenticated))
	if err != nil {
		logger.LogInfo("UserModule::User::Auth()/Error authenticating user:" + err.Error())
		var appState = CdAppState{false, err.Error(), "", "", ""}
		var appData = RespData{Data: users, RowsAffected: 0, NumberOfResult: 1}
		resp := CdResponse{AppState: appState, Data: appData}
		return resp
	}

	respMsg := ""
	sid := ""
	if authenticated {
		respMsg = "User authenticated successfully"
		logger.LogInfo("UserModule::User::Auth()/respMsg:" + respMsg)
		sidInt, err := CreateSess(req)
		if err != nil {
			logger.LogInfo("UserModule::User::Auth()/Error creating sesson:" + err.Error())
			var appState = CdAppState{false, err.Error(), "", "", ""}
			var appData = RespData{Data: users, RowsAffected: 0, NumberOfResult: 1}
			resp := CdResponse{AppState: appState, Data: appData}
			return resp
		}
		sid = strconv.Itoa(sidInt)
	} else {
		respMsg = "User authentication failed"
		logger.LogInfo("UserModule::User::Auth()/respMsg:" + respMsg)
	}

	fmt.Println("cd-user/Auth(): SessionID:", sid)

	var appState = CdAppState{authenticated, respMsg, "", "", ""}
	appState.Sess = sid
	var appData = RespData{Data: users, RowsAffected: 0, NumberOfResult: 1}
	resp := CdResponse{AppState: appState, Data: appData}
	return resp
}

// AuthenticateUser authenticates a user by username and password
func AuthenticateUser(username, password string) (bool, error) {
	logger.LogInfo("Starting UserModule::User::AuthenticateUser()...")
	logger.LogInfo("UserModule::User::AuthenticateUser()/username:" + username)
	logger.LogInfo("UserModule::User::AuthenticateUser()/password:" + password)
	var users User

	// Strategy 1: select only the queried user and expect one result only
	// result := db.Where("username = ?", username).First(&user)
	//
	/*
		// Strategy 2: get requested user and 'anon' data/ anon data is used in case of failure
		db.Where("role = ?", "admin").Or("role = ?", "super_admin").Find(&users)
		// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';
	*/
	// result := db.Where("username = ?", username).Or("username = ?", "anon").Find(&users)
	result := db.Table("user").Select("user_id", "user_name", "password").Where("username = ?", username).Scan(&users)
	if result.Error != nil {
		return false, result.Error
	}
	logger.LogInfo("UserModule::User::AuthenticateUser()/result:" + fmt.Sprint(result))
	logger.LogInfo("UserModule::User::AuthenticateUser()/user.Password:" + fmt.Sprint(users.Password))
	return CheckPasswordHash(password, users.Password), nil
}

// CheckPasswordHash compares a hashed password with its plaintext version
func CheckPasswordHash(password, hash string) bool {
	logger.LogInfo("Starting UserModule::User::CheckPasswordHash()...")
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		logger.LogWarning("UserModule::User::CheckPasswordHash(): password verification failed!")
		log.Println(err)
		return false
	}
	logger.LogInfo("UserModule::User::CheckPasswordHash(): password verification success!")
	return true
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CreateUser(req CdRequest) CdResponse {
	logger.LogInfo("Starting UserModule::User::CreateUser()...")
	var user User
	user.UserName = req.Dat.F_vals.Data.UserName
	user.UserGuid = fmt.Sprint(uuid.New())
	userResult := db.Create(&user)
	if userResult.Error != nil {
		logger.LogInfo("UserModule::User::CreateUser()/Error creating user:" + fmt.Sprint(userResult.Error))
		var appState = CdAppState{false, fmt.Sprint(userResult.Error), "", "", ""}
		var appData = RespData{Data: []User{}, RowsAffected: 0, NumberOfResult: 1}
		resp := CdResponse{AppState: appState, Data: appData}
		return resp
	}
	logger.LogInfo("UserModule::Session::SessCreate()/result:" + fmt.Sprint(userResult))
	var appState = CdAppState{true, "User registered successfully", "", "", ""}

	sid := ""
	sidInt, err := CreateSess(req)
	if err != nil {
		logger.LogInfo("UserModule::User::Auth()/Error creating sesson:" + err.Error())
		var appState = CdAppState{false, err.Error(), "", "", ""}
		var appData = RespData{Data: []User{}, RowsAffected: 0, NumberOfResult: 1}
		resp := CdResponse{AppState: appState, Data: appData}
		return resp
	}
	sid = strconv.Itoa(sidInt)
	appState.Sess = sid
	var appData = RespData{Data: []User{}, RowsAffected: 0, NumberOfResult: 1}
	resp := CdResponse{AppState: appState, Data: appData}
	return resp
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
