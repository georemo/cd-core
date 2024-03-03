package user

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

var mc = memcache.New("localhost:11211")

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

func Auth(string) bool {
	return true
}
