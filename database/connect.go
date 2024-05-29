package database

import (
	"app/config"
	"app/model"
	"fmt"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB connect to db
func ConnectDB() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)

	if err != nil {
		panic("failed to parse database port")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Config("DB_HOST"),
		port,
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	//fmt.Println("Connection Opened to Database")
	err = DB.AutoMigrate(
		&model.User{},
		&model.Home{},
		&model.HomeSetting{},
		&model.Room{},
		&model.Device{},
		&model.DeviceLog{},
		&model.RoomShare{},
		&model.DeviceType{},
		&model.DeviceStatus{})
	//if err != nil {
	//	return
	//}
	//fmt.Println("Database Migrated")
}
