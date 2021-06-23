package data

import (
	"fmt"
	"saltgram/internal"
	"time"

	"github.com/sirupsen/logrus"
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
		db.l.Info("Reattempting connection to User database")
		time.Sleep(5000 * time.Duration(time.Millisecond))
		dbtmp, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if dbtmp != nil {
		db.l.Info("Connected to User db")
	}
	db.DB = dbtmp
	return err
}

func (db *DBConn) MigradeData() {
	db.DB.AutoMigrate(&User{})
	db.DB.AutoMigrate(&Profile{})
	db.DB.AutoMigrate(&FollowRequest{})
}

func (db *DBConn) SeedAdmin() {
	if db.DB.Where("username = ?", "admin").First(&User{}).RowsAffected > 0 {
		return
	}

	user := &User{
		Username:       "admin",
		HashedPassword: "Saltadmin123!?",
		Role:           "admin",
		Activated:      true,
	}
	err := user.GenerateSaltAndHashedPassword()
	if err != nil {
		db.l.Fatalf("Hashing admin password failed")
	}
	db.DB.Create(user)
}
