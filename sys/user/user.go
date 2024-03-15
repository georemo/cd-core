package user

import (
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/tcp-x/cd-core/sys/base"
)

var db = base.Conn()
var mc = memcache.New("localhost:11211")

type IUser interface{}

// User model
type User struct {
	UserId        uint      `gorm:"primaryKey"`
	UserGuid      string    `json:"user_guid"`
	UserName      string    `json:"user_name"`
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

// Auth
func (u User) Auth(req string) (string, error) {
	// var records []User
	// usr := base.Get("user", records, db)

	// get user and anon data
	// 1. convert req to struct
	reqStruct, err := base.JSONToICdRequest(req)
	if err != nil {
		fmt.Println("Error:", err)
		return "", nil
	}

	// Accessing fields of MyStruct
	fmt.Println("Module:", reqStruct.M)
	fmt.Println("Dat:", reqStruct.Dat)

	// this.plData = this.b.getPlData(req);
	//     const q: IQuery = {
	//         // get requested user and 'anon' data/ anon data is used in case of failure
	//         where: [
	//             { userName: this.plData.userName },
	//             { userName: "anon" }
	//         ]
	//     };

	// connect to db and check validity of password
	// Auth input should have username and password

	// test if /tcp-x/user/session is accessible
	sid := SessID()
	fmt.Println("cd-user/Auth(): SessionID:", sid)

	// resp := "{name:User, version:0.0.7 publisher: \"EMP Services Ltd\"}"
	resp := `{"state":"success", "data":[]}`
	return resp, nil
}
