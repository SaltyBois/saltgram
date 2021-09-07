package data

import (
	"agent/internal"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	HOST_DB = internal.GetEnvOrDefault("AGENTDB_HOST", "localhost")
	PORT_DB = internal.GetEnvOrDefaultInt("AGENTDB_PORT", 5432)
	USER_DB = internal.GetEnvOrDefault("AGENTDB_USER", "saltagent")
	PASS_DB = internal.GetEnvOrDefault("AGENTDB_PASS", "saltagent")
	NAME_DB = internal.GetEnvOrDefault("AGENTDB_NAME", "saltagentdb")
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
		db.l.Info("Reattempting connection to Agent db")
		time.Sleep(time.Millisecond * 5000)
		dbtmp, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if dbtmp != nil {
		db.l.Info("Connected to Agent db!")
	}

	db.DB = dbtmp
	return err
}

func (db *DBConn) MigradeData() {
	db.DB.AutoMigrate(&User{})
	db.DB.AutoMigrate(&Campaign{})
	db.DB.AutoMigrate(&Product{})
	db.DB.AutoMigrate(&Refresh{})
}
