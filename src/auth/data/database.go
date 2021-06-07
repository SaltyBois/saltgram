package data

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	HOST_DB = "auth_db"
	PORT_DB = 9090
	USER_DB = "saltauth"
	PASS_DB = "saltauth"
	NAME_DB = "saltauthdb"
)

type DBConn struct {
	DB *gorm.DB
}

func (db *DBConn) ConnectToDb() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		HOST_DB, USER_DB, PASS_DB, NAME_DB, PORT_DB)
	dbtmp, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    for err != nil {
        time.Sleep(time.Millisecond * 2000)
        dbtmp, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    }
	db.DB = dbtmp
	return err
}

func (db *DBConn) MigradeData() {
	db.DB.AutoMigrate(&Refresh{})
}
