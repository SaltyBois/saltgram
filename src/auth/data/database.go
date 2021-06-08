package data

import (
	"fmt"
	"log"
	"saltgram/internal"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	HOST_DB = internal.GetEnvOrDefault("AUTHDB_HOST", "localhost")
	PORT_DB = internal.GetEnvOrDefaultInt("AUTHDB_PORT", 5432)
	USER_DB = internal.GetEnvOrDefault("AUTHDB_USER", "saltauth")
	PASS_DB = internal.GetEnvOrDefault("AUTHDB_PASS", "saltauth")
	NAME_DB = internal.GetEnvOrDefault("AUTHDB_NAME", "saltauthdb")
)

type DBConn struct {
	DB *gorm.DB
	l  *log.Logger
}

func NewDBConn(l *log.Logger) *DBConn {
	return &DBConn{l: l}
}

func (db *DBConn) ConnectToDb() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		HOST_DB, USER_DB, PASS_DB, NAME_DB, PORT_DB)
	dbtmp, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	for err != nil {
		db.l.Print("Reattempting connection to Auth db")
		time.Sleep(time.Millisecond * 5000)
		dbtmp, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if dbtmp != nil {
		db.l.Printf("Connected to Auth db!\n")
	}

	db.DB = dbtmp
	return err
}

func (db *DBConn) MigradeData() {
	db.DB.AutoMigrate(&Refresh{})
}
