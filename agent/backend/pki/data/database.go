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
	HOST_DB = internal.GetEnvOrDefault("PKIDB_HOST", "localhost")
	PORT_DB = internal.GetEnvOrDefaultInt("PKIDB_PORT", 5432)
	USER_DB = internal.GetEnvOrDefault("PKIDB_USER", "saltpki")
	PASS_DB = internal.GetEnvOrDefault("PKIDB_PASS", "saltpki")
	NAME_DB = internal.GetEnvOrDefault("PKIDB_NAME", "saltpkidb")
)

type DBConn struct {
	DB *gorm.DB
	l  *logrus.Logger
}

func NewDBConn(l *logrus.Logger) *DBConn {
	return &DBConn{l: l}
}

func (db *DBConn) ConnectToDB() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		HOST_DB, USER_DB, PASS_DB, NAME_DB, PORT_DB)
	dbtmp, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	for err != nil {
		db.l.Info("Reattempting connection to PKI db")
		time.Sleep(time.Millisecond * 5000)
		dbtmp, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if dbtmp != nil {
		db.l.Info("Connected to PKI db!")
	}

	db.DB = dbtmp
	return err
}

func (db *DBConn) MigrateData() {
	db.DB.AutoMigrate(&ArchivedCert{})
}
