package data

import (
	"fmt"
	"time"

	"saltgram/internal"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	HOST_DB = internal.GetEnvOrDefault("CONTENTSDB_HOST", "localhost")
	PORT_DB = internal.GetEnvOrDefaultInt("CONTENTSDB_PORT", 5432)
	USER_DB = internal.GetEnvOrDefault("CONTENTSDB_USER", "saltcontents")
	PASS_DB = internal.GetEnvOrDefault("CONTENTSDB_PASS", "saltcontents")
	NAME_DB = internal.GetEnvOrDefault("CONTENTSDB_NAME", "saltcontentsdb")
)

type DBConn struct {
	DB *gorm.DB
	l  *logrus.Logger
}

func NewDBConn(l *logrus.Logger) *DBConn {
	return &DBConn{l: l}
}

func (db *DBConn) ConnectToDb() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		HOST_DB, USER_DB, PASS_DB, NAME_DB, PORT_DB)
	dbtmp, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	for err != nil {
		db.l.Info("Reattempting connection to Content database")
		time.Sleep(5000 * time.Duration(time.Millisecond))
		dbtmp, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if dbtmp != nil {
		db.l.Info("Connected to Content db")
	}
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
	db.DB.AutoMigrate(&ProfilePicture{})
	db.DB.AutoMigrate(&Highlight{})
}
