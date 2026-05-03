package providers

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	. "safa-went/internal/config"
)

type DatabaseProvider struct {
	DB *gorm.DB
}

func (dp *DatabaseProvider) Connect(config *Config) {

	dbConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Silent),
	}

	var conn gorm.Dialector

	switch config.DBDialect {
	case "sqlite":
		conn = sqlite.Open(config.DBName)
	case "mysql":
		conn = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.DBUser,
			config.DBPassword,
			config.DBHost,
			config.DBPort,
			config.DBName,
		))
	case "postgres":
		conn = postgres.Open(fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
			config.DBHost,
			config.DBPort,
			config.DBUser,
			config.DBName,
			config.DBPassword,
		))
	default:
		panic("Unsupported DB dialect")
	}

	db, err := gorm.Open(conn, dbConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect database: %v", err))
	}
	dp.DB = db
}
