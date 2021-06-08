package data

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	HOST_DB = "localhost"
	PORT_DB = 5432
	USER_DB = "saltcontents"
	PASS_DB = "saltcontents"
	NAME_DB = "saltcontentsdb"
)

type DBConn struct {
	DB *gorm.DB
}

func (db *DBConn) ConnectToDb() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		HOST_DB, USER_DB, PASS_DB, NAME_DB, PORT_DB)
	dbtmp, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.DB = dbtmp
	return err
}

func (db *DBConn) MigradeData() {
	db.DB.AutoMigrate(&Media{})
	db.DB.AutoMigrate(&SharedMedia{})
	db.DB.AutoMigrate(&Post{})
	db.DB.AutoMigrate(&Story{})
	db.DB.AutoMigrate(&Reaction{})
	db.DB.AutoMigrate(&Comment{})
}
