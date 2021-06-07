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
	HOST_DB = internal.GetEnvOrDefault("USERSDB_HOST", "localhost")
	PORT_DB = internal.GetEnvOrDefaultInt("USERSDB_PORT", 5432)
	USER_DB = internal.GetEnvOrDefault("USERSDB_USER", "saltusers")
	PASS_DB = internal.GetEnvOrDefault("USERSDB_PASS", "saltusers")
	NAME_DB = internal.GetEnvOrDefault("USERSDB_NAME", "saltusersdb")
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
		db.l.Print("Reattempting connection to User database")
		time.Sleep(5000 * time.Duration(time.Millisecond))
		dbtmp, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if dbtmp != nil {
		db.l.Print("Connected to User db")
	}
	db.DB = dbtmp
	return err
}

func (db *DBConn) MigradeData() {
	db.DB.AutoMigrate(&User{})
}
