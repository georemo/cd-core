/*
This is a golang code that uses GORM to connect and execute CRUD against mysql for any table
that is not known at the time of coding.  The input should be a model based on a table definition.
ChatGPT

To achieve this, you can use reflection in Go to dynamically generate GORM models based on the
table definition and perform CRUD operations. Here's an example code that demonstrates how to do this:

In this code:

    We define a generic TableRow struct to embed common fields such as ID, CreatedAt, UpdatedAt, and DeletedAt.
    We use reflection to dynamically generate a model for a specific table based on its schema.
	You need to define the table's schema by adding appropriate reflect.StructField elements.
    We perform CRUD operations using the dynamically created table model. You can adapt the code to
	handle your specific table's schema and operations accordingly.

Please replace "user:password@tcp(localhost:3306)/dbname" with your actual MySQL connection string and
"your_table_name" with the name of the table you want to interact with. Also, adjust the code to match
your table's schema by adding appropriate struct fields and tags.

*/

package base

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Define a struct that represents a generic table row
type TableRow struct {
	gorm.Model
}

func Conn() *gorm.DB {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database connection parameters from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construct MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to MySQL database
	// dsn := "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return db
}

func Create(tableName string, records []interface{}, db *gorm.DB) {
	// Define a struct for the table schema
	// type YourTable struct {
	// 	TableRow
	// 	// Define fields here
	// 	// Example: Name string `gorm:"column:name"`
	// }

	// Dynamically create table model-
	// tableName := "your_table_name"
	tableModel := reflect.New(reflect.StructOf([]reflect.StructField{
		reflect.StructField{Name: "TableRow", Type: reflect.TypeOf(TableRow{})},
		// Add struct fields based on the table schema
		// Example: reflect.StructField{Name: "Name", Type: reflect.TypeOf(""), Tag: reflect.StructTag(`gorm:"column:name"`)},
	})).Interface()

	// Perform CRUD operations using the dynamically created table model
	// Example: Read data
	// var records []interface{}
	db.Table(tableName).Find(&records)
	fmt.Println("Records:", records)

	// Example: Create data
	newRecord := reflect.New(reflect.TypeOf(tableModel).Elem()).Interface()
	// Set field values
	// Example: newRecord.(YourTable).Name = "NewRecord"
	db.Table(tableName).Create(newRecord)
	// Delete - delete tableName
	db.Delete(&tableName, 1)
}

func Create2(db *gorm.DB) *gorm.DB {
	type User struct {
		Name     string
		Age      int
		Birthday time.Time
	}
	user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}
	return db.Create(&user) // pass pointer of data to Create
}

func Get(tableName string, records []interface{}, db *gorm.DB) []interface{} {
	// Define a struct for the table schema
	type YourTable struct {
		TableRow
		// Define fields here
		// Example: Name string `gorm:"column:name"`
	}

	db.Table(tableName).Find(&records)
	fmt.Println("Records:", records)
	// Delete - delete tableName
	db.Delete(&tableName, 1)
	return records
}

func Get2(db *gorm.DB, tableName string, jsonInput string) (string, error) {
	// Parse the JSON input
	var searchParams map[string]interface{}
	err := json.Unmarshal([]byte(jsonInput), &searchParams)
	if err != nil {
		return "", err
	}

	// Build the query
	query := db.Table(tableName)
	for key, value := range searchParams {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Execute the query
	result := query.Find(&[]map[string]interface{}{})
	if result.Error != nil {
		return "", result.Error
	}

	// Convert result to JSON string
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(resultJSON), nil
}

func Update(tableName string, id string, db *gorm.DB) {
	// Example: Update data
	db.Table(tableName).Where("id = ?", id).Updates(map[string]interface{}{"field": "value"})
	// Delete - delete tableName
	db.Delete(&tableName, 1)
}

func Delete(tableName string, id string, db *gorm.DB) {
	tableModel := reflect.New(reflect.StructOf([]reflect.StructField{
		reflect.StructField{Name: "TableRow", Type: reflect.TypeOf(TableRow{})},
		// Add struct fields based on the table schema
		// Example: reflect.StructField{Name: "Name", Type: reflect.TypeOf(""), Tag: reflect.StructTag(`gorm:"column:name"`)},
	})).Interface()
	newRecord := reflect.New(reflect.TypeOf(tableModel).Elem()).Interface()
	// Example: Delete data
	db.Table(tableName).Where("id = ?", id).Delete(newRecord)
	// Delete - delete tableName
	db.Delete(&tableName, 1)
}

// func main() {

// 	db := Conn()

// 	// Define a struct for the table schema
// 	type YourTable struct {
// 		TableRow
// 		// Define fields here
// 		// Example: Name string `gorm:"column:name"`
// 	}

// 	// Dynamically create table model
// 	tableName := "your_table_name"
// 	tableModel := reflect.New(reflect.StructOf([]reflect.StructField{
// 		reflect.StructField{Name: "TableRow", Type: reflect.TypeOf(TableRow{})},
// 		// Add struct fields based on the table schema
// 		// Example: reflect.StructField{Name: "Name", Type: reflect.TypeOf(""), Tag: reflect.StructTag(`gorm:"column:name"`)},
// 	})).Interface()

// 	// Perform CRUD operations using the dynamically created table model
// 	// Example: Read data
// 	var records []interface{}
// 	db.Table(tableName).Find(&records)
// 	fmt.Println("Records:", records)

// 	// Example: Create data
// 	newRecord := reflect.New(reflect.TypeOf(tableModel).Elem()).Interface()
// 	// Set field values
// 	// Example: newRecord.(YourTable).Name = "NewRecord"
// 	db.Table(tableName).Create(newRecord)

// 	// Example: Update data
// 	db.Table(tableName).Where("id = ?", id).Updates(map[string]interface{}{"field": "value"})

// 	// Example: Delete data
// 	db.Table(tableName).Where("id = ?", id).Delete(newRecord)

// 	// Delete - delete tableName
// 	db.Delete(&tableName, 1)
// }
