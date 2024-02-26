package base

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	ID    uint   `gorm:"primary_key"`
	Code  string `gorm:"uniqueIndex"`
	Price uint
}

func main() {
	// Connect to MySQL database
	dsn := "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Automatically create the table if it doesn't exist
	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	// Create some products
	products := []Product{
		{Code: "Laptop", Price: 1000},
		{Code: "Phone", Price: 500},
		{Code: "Tablet", Price: 300},
	}

	// Create products concurrently
	var wg sync.WaitGroup
	for _, p := range products {
		wg.Add(1)
		go func(p Product) {
			defer wg.Done()
			if err := db.Create(&p).Error; err != nil {
				log.Printf("Error creating product: %v", err)
			}
		}(p)
	}
	wg.Wait()

	// Read products concurrently
	var readWG sync.WaitGroup
	readWG.Add(len(products))
	for _, p := range products {
		go func(p Product) {
			defer readWG.Done()
			var product Product
			if err := db.First(&product, "code = ?", p.Code).Error; err != nil {
				log.Printf("Error reading product: %v", err)
				return
			}
			fmt.Printf("Read product: %+v\n", product)
		}(p)
	}
	readWG.Wait()

	// Update products concurrently
	var updateWG sync.WaitGroup
	updateWG.Add(len(products))
	for _, p := range products {
		go func(p Product) {
			defer updateWG.Done()
			if err := db.Model(&Product{}).Where("code = ?", p.Code).Update("price", p.Price+100).Error; err != nil {
				log.Printf("Error updating product: %v", err)
			}
		}(p)
	}
	updateWG.Wait()

	// Delete products concurrently
	var deleteWG sync.WaitGroup
	deleteWG.Add(len(products))
	for _, p := range products {
		go func(p Product) {
			defer deleteWG.Done()
			if err := db.Where("code = ?", p.Code).Delete(&Product{}).Error; err != nil {
				log.Printf("Error deleting product: %v", err)
			}
		}(p)
	}
	deleteWG.Wait()

	// Close database connection
	db.Close()
}
