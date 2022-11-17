package datastore

import (
	"enterprise.sidooh/pkg/entities"
	"enterprise.sidooh/utils"
	"errors"
	permify "github.com/Permify/permify-gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"strings"
	"time"
)

var (
	DB      *gorm.DB
	Permify *permify.Permify
)

func Init() {
	dsn := viper.GetString("DB_DSN")
	env := strings.ToUpper(viper.GetString("APP_ENV"))
	if env != "TEST" && dsn == "" {
		panic("database connection not set")
	}

	var gormDb = new(gorm.DB)
	var err = errors.New("")

	if env == "TEST" {
		gormDb, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else {
		file := utils.GetLogFile("db.log")

		newLogger := logger.New(
			log.New(file, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				LogLevel: logger.Info, // Log level
			},
		)

		config := &gorm.Config{
			Logger: newLogger,
		}

		gormDb, err = gorm.Open(mysql.Open(dsn), config)
		db, _ := gormDb.DB()
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(100)
		db.SetConnMaxLifetime(time.Hour)
	}

	if err != nil {
		logrus.Error(err)
		panic("failed to connect database")
	}

	if viper.GetBool("MIGRATE_DB") {
		err := gormDb.AutoMigrate(
			&entities.Enterprise{},
			&entities.User{},
			&entities.Account{},
			&entities.Team{},
			//&entities.TeamAccount{},
		)
		if err != nil {
			logrus.Error(err)
			panic("failed to auto-migrate")
		}
		logrus.Println("Auto-migrated db")
	}

	// setup permify
	Permify, err = permify.New(permify.Options{
		Migrate: viper.GetBool("MIGRATE_DB"),
		DB:      gormDb,
	})
	if err != nil {
		logrus.Error(err)
		panic("failed to migrate permissions")
	}

	DB = gormDb
}
