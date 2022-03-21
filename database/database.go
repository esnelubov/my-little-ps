package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"my-little-ps/config"
)

type DB struct {
	gormDB *gorm.DB
}

func New(config config.IConfig) *DB {
	dsn := CreateDatabaseAndReturnDSN(config)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	return &DB{
		gormDB: db,
	}
}

func CreateDatabaseAndReturnDSN(config config.IConfig) string {
	dsn := config.GetString("dsn")
	dbName := config.GetString("dbName")
	dbParams := config.GetString("dbParams")

	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName))

	return fmt.Sprintf("%s/%s?%s", dsn, dbName, dbParams)
}
