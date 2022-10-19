package datastore

import (
	"enterprise.sidooh/pkg/entities"
	"enterprise.sidooh/utils"
	"errors"
	"fmt"
	permify "github.com/Permify/permify-gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"strings"
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
	}

	if err != nil {
		logrus.Println(err)
		panic("failed to connect database")
	}

	if viper.GetBool("MIGRATE_DB") {
		gormDb.AutoMigrate(&entities.Enterprise{}, &entities.Account{})
		fmt.Println("Auto-migrated db")
	}

	// setup permify
	Permify, _ = permify.New(permify.Options{
		Migrate: viper.GetBool("MIGRATE_DB"),
		DB:      gormDb,
	})

	DB = gormDb
}
