package user

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/tcp-x/cd-core/sys/base"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var db = base.Conn()
var mc = memcache.New("localhost:11211")
var logger base.Logger
var jSessData datatypes.JSON
var respMsg = ""

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
	Sess    Session
	Cache   string
	SConfig string
}

type RespData struct {
	Data           *gorm.DB
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

type ServiceInput struct {
	ServiceModel any
	ModelName    string
	DocName      string
	Cmd          Cmd
	DSource      int
}

type Cmd struct {
	Action string
	Query  json.RawMessage
}

var userResult *gorm.DB

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

var u, anon User

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
func Auth(req base.CdRequest) base.CdResponse {
	logger.LogInfo("Module version:v0.0.60")
	logger.LogInfo("Starting UserModule::User::Auth()...")
	// var users []User
	if user, ok := req.Dat.F_vals.Data.(User); ok {
		// If the data is of type User, print the user details
		logger.LogInfo("UserModule::User::Auth()/req.Dat.F_vals.Data.UserName:" + user.UserName)
		logger.LogInfo("UserModule::User::Auth()/req.Dat.F_vals.Data.Password:" + user.Password)
		authenticated, err := AuthenticateUser(user.UserName, user.Password)
		logger.LogInfo("UserModule::User::Auth()/authenticated:" + fmt.Sprint(authenticated))
		if err != nil {
			logger.LogInfo("UserModule::User::Auth()/Error authenticating user:" + err.Error())
			var appState = base.CdAppState{false, err.Error(), nil, "", ""}
			var appData = base.RespData{Data: nil, RowsAffected: 0, NumberOfResult: 1}
			resp := base.CdResponse{AppState: appState, Data: appData}
			return resp
		}

		// sid := ""
		if authenticated {
			respMsg = "User authenticated successfully"
			logger.LogInfo("UserModule::User::Auth()/respMsg:" + respMsg)

			if err != nil {
				logger.LogInfo("UserModule::User::Auth()/Error authenticating user:" + err.Error())
				var appState = base.CdAppState{false, err.Error(), nil, "", ""}
				var appData = base.RespData{Data: nil, RowsAffected: 0, NumberOfResult: 1}
				resp := base.CdResponse{AppState: appState, Data: appData}
				return resp
			}

		} else {
			respMsg = "User authentication failed"
			logger.LogInfo("UserModule::User::Auth()/respMsg:" + respMsg)
			logger.LogWarning("UserModule::User::Auth()/Warning:" + respMsg)
			var appState = base.CdAppState{false, respMsg, nil, "", ""}
			var appData = base.RespData{Data: nil, RowsAffected: 0, NumberOfResult: 0}
			resp := base.CdResponse{AppState: appState, Data: appData}
			return resp
		}

	} else {
		// If the data is not of type User, print an error message
		fmt.Println("Data is not of type User")
	}

	sessResp := CreateSess(req)
	// jSessData, err := json.Marshal(sessResp.Data)
	// if err != nil {
	// 	logger.LogInfo("UserModule::User::Auth()/Error creating sesson:" + err.Error())
	// 	var appState = base.CdAppState{false, err.Error(), nil, "", ""}
	// 	var appData = base.RespData{Data: []User{anon}, RowsAffected: 0, NumberOfResult: 1}
	// 	resp := base.CdResponse{AppState: appState, Data: appData}
	// 	return resp
	// }

	// process authenticated response
	// resp := AuthResponse(req)
	var appState = base.CdAppState{true, respMsg, sessResp.AppState.Sess, "", ""}
	// appState.Sess = sid
	var appData = base.RespData{Data: []User{u}, RowsAffected: 0, NumberOfResult: 1}
	resp := base.CdResponse{AppState: appState, Data: appData}
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
		db.Where("user_name = ?", username).Or("user_name = ?", "anon").Find(&users)
		// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';
	*/

	// result := db.Where("username = ?", username).Or("username = ?", "anon").Find(&users)

	// type Result struct {
	// 	UserId   string
	// 	UserName string
	// 	Password string
	// }
	// var result Result
	userResult := db.Table("user").
		Select("user_id", "user_name", "password").
		Where("user_name = ?", username).
		Or("user_name = ?", "anon").
		Scan(&users)
	if userResult.Error != nil {
		return false, userResult.Error
	}

	rows, err := userResult.Rows()
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		db.ScanRows(rows, &users)
		fmt.Println(users)
		if users.UserName == username {
			u = users
		}
		if users.UserName == "anon" {
			anon = users
		}
	}

	logger.LogInfo("UserModule::User::AuthenticateUser()/result:" + fmt.Sprint(userResult))
	logger.LogInfo("UserModule::User::AuthenticateUser()/user.Password:" + fmt.Sprint(users.Password))
	return CheckPasswordHash(password, u.Password), nil
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

// func CreateUser(req base.CdRequest) base.CdResponse {
// 	logger.LogInfo("Starting UserModule::User::CreateUser()...")
// 	var user User
// 	user.UserName = req.Dat.F_vals.Data.UserName
// 	user.UserGuid = fmt.Sprint(uuid.New())
// 	userResult := db.Create(&user)
// 	if userResult.Error != nil {
// 		logger.LogInfo("UserModule::User::CreateUser()/Error creating user:" + fmt.Sprint(userResult.Error))
// 		var appState = CdAppState{false, fmt.Sprint(userResult.Error), "", "", ""}
// 		var appData = RespData{Data: []User{}, RowsAffected: 0, NumberOfResult: 1}
// 		resp := CdResponse{AppState: appState, Data: appData}
// 		return resp
// 	}
// 	logger.LogInfo("UserModule::Session::SessCreate()/result:" + fmt.Sprint(userResult))
// 	var appState = CdAppState{true, "User registered successfully", "", "", ""}

// 	sid := ""
// 	sidInt, err := CreateSess(req)
// 	if err != nil {
// 		logger.LogInfo("UserModule::User::Auth()/Error creating sesson:" + err.Error())
// 		var appState = CdAppState{false, err.Error(), "", "", ""}
// 		var appData = RespData{Data: []User{}, RowsAffected: 0, NumberOfResult: 1}
// 		resp := CdResponse{AppState: appState, Data: appData}
// 		return resp
// 	}
// 	sid = strconv.Itoa(sidInt)
// 	appState.Sess = sid
// 	var appData = RespData{Data: []User{}, RowsAffected: 0, NumberOfResult: 1}
// 	resp := CdResponse{AppState: appState, Data: appData}
// 	return resp
// }

// func AuthResponse(req CdRequest) CdResponse {
// }

/*
sessResult: sessResult$,

	modulesUserData: {
						consumer: [],
						menuData: [],
						userData: {}
					}
*/
// func processResponse(req CdRequest) CdResponse {
// 	// this process can be handled using seperate subroutines
// }

// // to be move to user/session module
// func GetSessConfig(req CdRequest) CdResponse {

// }

// // to be move to modules module
// func GetModulesUserData(req CdRequest) CdResponse {

// }

// // to be move to menu module
// func GetMenuData(req CdRequest) CdResponse {

// }

// type Module struct{}
// type Menu struct{}
// type MenuTree struct{}
// type UserRole struct {}
// type AclData struct {}
// type CUser struct {
// 	currentUser  User
// 	consumerGuid string
// }

// // to be move to menu module
// func GetAclMenu(req CdRequest) MenuTree {
// 	moduleMenuData := getModuleMenu()
// 	rootMenuId := getRootMenuId(moduleMenuData)
// 	menuData := buildNestedMenu(rootMenuId, moduleMenuData)
// }

// /*
// get allowed modules by the consumer and user
// */
// func GetAclModules(req CdRequest, cUser CUser) []Module {
// 	this.consumerGuid = params.consumerGuid;
// 	return forkJoin({
// 		userRoles: this.srvAcl.aclUser(req, res, params).pipe(map((m) => { return m })),
// 		consumerModules: this.srvAcl.aclModule$(req, res).pipe(map((m) => { return m })),
// 		moduleParents: this.srvAcl.aclModuleMembers$(req, res, params).pipe(map((m) => { return m }))
// 	})
// 	.pipe(
// 		map((acl: any) => {
// 			// this.b.logTimeStamp('ModuleService::getModulesUserData$/02')
// 			// console.log('ModuleService::getAclModule$()/acl:', acl)
// 			// Based on acl result, return appropirate modules
// 			const publicModules = acl.consumerModules.filter(m => m.moduleIsPublic);
// 			if (acl.userRoles.isConsumerRoot.length > 0) { // if userIsConsumerRoot then return all consumerModules
// 				// this.b.logTimeStamp('ModuleService::getModulesUserData$/03')
// 				return acl.consumerModules;
// 			}
// 			else if (acl.userRoles.isConsumerUser.length > 0) { // if user is registered as consumer user then filter consumer modules
// 				// this.b.logTimeStamp('ModuleService::getModulesUserData$/04')
// 				// console.log('ModuleService::getModulesUserData$/acl.userRoles.isConsumerUser:', acl.userRoles.isConsumerUser);
// 				// console.log('ModuleService::getModulesUserData$/acl.moduleParents:', acl.moduleParents);
// 				// console.log('ModuleService::getModulesUserData$/acl.consumerModules:', acl.consumerModules);
// 				const userModules = this.b.intersect(acl.consumerModules, acl.moduleParents, 'moduleGuid');
// 				// console.log('ModuleService::getModulesUserData$/userModules:', userModules);
// 				return userModules.concat(publicModules); // return user modules and public modules
// 			}
// 			else {  // if is neither of the above, return zero modules
// 				// console.log('ModuleService::getAclModule$()/publicModules:', publicModules)
// 				return publicModules; // return only public modules
// 			}
// 		})
// 	);
// }

// func GetConsumerModules(consumerGuid string) []Module {

// }

// /*
//  * get users based on AclUserViewModel and
//  * filtered by current consumer relationship and user role
//  * return user role
//  */
// func aclUser(req, res, params) []UserRole{
// 	const b = new BaseService();
// 	this.consumerGuid = params.consumerGuid;
// 	const q: IQuery = { where: {} };
// 	serviceInput := ServiceInput{
// 		serviceModel AclUserViewModel,
// 		modelName: "AclUserViewModel",
// 		docName: 'AclService::aclUser',
// 		cmd: {
// 			action: 'find',
// 			query: q
// 		},
// 		dSource: 1
// 	}
// 	const user$ = from(b.read(req, res, serviceInput))
// 		.pipe(
// 			share() // to avoid repeated db round trips
// 		)
// 	const isRoot = u => u.userId === 1001;

// 	const isConsumerRoot = u => u.consumerResourceTypeId === 4
// 		&& u.consumerGuid === this.consumerGuid
// 		&& u.objGuid === params.currentUser.userGuid;

// 	const isConsumerTechie = u => u.consumerResourceTypeId === 5
// 		&& u.consumerGuid === this.consumerGuid
// 		&& u.objGuid === params.currentUser.userGuid;

// 	const isConsumerUser = u => u.consumerResourceTypeId === 6
// 		&& u.consumerGuid === this.consumerGuid
// 		&& u.objGuid === params.currentUser.userGuid;

// 	const isRoot$ = user$
// 		.pipe(
// 			map((u) => {
// 				const ret = u.filter(isRoot)
// 				// this.b.logTimeStamp(`AclService::aclUser$/u[isRoot$]:${JSON.stringify(u)}`)
// 				// this.b.logTimeStamp(`AclService::aclUser$/ret[isRoot$]:${JSON.stringify(ret)}`)
// 				return ret;
// 			})
// 			, distinct()
// 		);

// 	const isConsumerRoot$ = user$
// 		.pipe(
// 			map((u) => {
// 				const ret = u.filter(isConsumerRoot)
// 				// this.b.logTimeStamp(`AclService::aclUser$/u[isConsumerRoot$]:${JSON.stringify(u)}`)
// 				// this.b.logTimeStamp(`AclService::aclUser$/ret[isConsumerRoot$]:${JSON.stringify(ret)}`)
// 				return ret;
// 			})
// 			, distinct()
// 		);

// 	const isConsumerTechie$ = user$
// 		.pipe(
// 			map((u) => {
// 				const ret = u.filter(isConsumerTechie)
// 				// this.b.logTimeStamp(`AclService::aclUser$/u[isConsumerTechie$]:${JSON.stringify(u)}`)
// 				// this.b.logTimeStamp(`AclService::aclUser$/ret[isConsumerTechie$]:${JSON.stringify(ret)}`)
// 				return ret;
// 			})
// 			, distinct()
// 		);

// 	const isConsumerUser$ = user$
// 		.pipe(
// 			map((u) => {
// 				const ret = u.filter(isConsumerUser)
// 				// this.b.logTimeStamp(`AclService::aclUser$/u[isConsumerUser$]:${JSON.stringify(u)}`)
// 				// this.b.logTimeStamp(`AclService::aclUser$/ret[isConsumerUser$]:${JSON.stringify(ret)}`)
// 				return ret;
// 			})
// 			, distinct()
// 		);

// 	return forkJoin(
// 		{
// 			isRoot: isRoot$.pipe(map((u) => { return u })),
// 			isConsumerRoot: isConsumerRoot$.pipe(map((u) => { return u })),
// 			isConsumerUser: isConsumerUser$.pipe(map((u) => { return u }))
// 		}
// 	)
// }

// func aclModule(req CdRequest){

// }

// func getModuleMenu(allowedModules []Module) MenuTree {

// }

// func buildNestedMenu(RootMenuId int, moduleMenuData Menu) MenuTree {

// }

// func getRootMenuId(moduleMenuData Menu) {

// }

// // to be move to module/consumer module
// func getConsumerData(req CdRequest) CdResponse {

// }

// // func fVals(fvals string) (FVals, error) {
// // 	fvalStruct, err := JSONToFVals(fvals)
// // 	if err != nil {
// // 		fmt.Println("Error:", err)
// // 		return fvalStruct, err
// // 	}
// // 	return fvalStruct, nil
// // }

// // Function to convert JSON string to fVals
// func JSONToFVals(jsonString string) (FVals, error) {
// 	var fVals FVals
// 	err := json.Unmarshal([]byte(jsonString), &fVals)
// 	return fVals, err
// }
