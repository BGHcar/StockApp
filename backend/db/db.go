package db

import (
	"backend/config"
	"fmt"
	"log"
	"backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Conect() (*gorm.DB, error) {
	var err error
	dsn := config.LoadDB().URL
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Can't get the conection with the database %v", err)
	}

	log.Println("Conection with :", DB.Statement.Table, " Established")

	if DB.Migrator().HasTable(config.LoadDB().TableName) {
		log.Println("The table: ", config.LoadDB().TableName, " exist")
	} else {
		log.Println("The table: ", config.LoadDB().TableName, " doesn't exist")
		if err := DB.AutoMigrate(&models.Stock{}); err != nil {
			return nil, fmt.Errorf("Failed to migrate: %v", err)
		}
	}
	return DB, nil
}

func Drop() error {
	var err error
	dsn := config.LoadDB().URL
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("can't get the conection with the database %v", err)
	}

	if err := DB.Migrator().DropTable(config.LoadDB().TableName); err != nil {
		return fmt.Errorf("failed to drop table: %v", err)
	}
	
	fmt.Println("the database has been deleted")

	return nil
}