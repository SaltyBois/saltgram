package data

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"saltgram/internal"
	"saltgram/users/data"
	"time"

	"github.com/sirupsen/logrus"
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
		db.l.Info("Reattempting connection to AuthDB")
		time.Sleep(time.Millisecond * 5000)
		dbtmp, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if dbtmp != nil {
		db.l.Info("Connected to Auth db!")
	}

	db.DB = dbtmp
	return err
}

func (db *DBConn) MigradeData() {
	db.DB.AutoMigrate(&Refresh{})
}

func (db *DBConn) SeedAdmin() {
	if db.DB.Where("username = ?", "admin").First(&Refresh{}).RowsAffected > 0 {
		return
	}

	refreshClaims := data.RefreshClaims{
		Username: "admin",
		StandardClaims: jwt.StandardClaims{
			// TODO(Jovan): Make programmatic?
			ExpiresAt: time.Now().UTC().AddDate(0, 6, 0).Unix(),
			Issuer:    "SaltGram",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	jws, err := token.SignedString([]byte(os.Getenv("REF_SECRET_KEY")))

	token2 := &Refresh{
		Username: "admin",
		Token:    jws,
	}
	if err != nil {
		db.l.Fatalf("Hashing admin password failed")
	}
	db.DB.Create(token2)
}
