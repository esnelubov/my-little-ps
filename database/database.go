package database

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"my-little-ps/config"
)

type DB struct {
	gormDB *gorm.DB
}

func New(config config.IConfig) *DB {
	var (
		db  *gorm.DB
		err error
	)

	err = AutoMigrate(config)
	if err != nil {
		log.Fatalf("failed to auto migrate schema: %s", err)
	}

	dsn := config.GetString("dsn")
	schemaName := config.GetString("schemaName")

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   schemaName + ".", // schema name
			SingularTable: false,
		}})
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	return &DB{
		gormDB: db,
	}
}

func AutoMigrate(config config.IConfig) error {
	var (
		err        error
		db         *sql.DB
		migration  *migrate.Migrate
		dsn        string
		schemaName string
		driver     database.Driver
	)

	autoMigrate := config.GetBool("autoMigrate")

	if !autoMigrate {
		return nil
	}

	dsn = config.GetString("dsn")
	schemaName = config.GetString("schemaName")

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %[1]s; SET search_path TO %[1]s;", schemaName))

	driver, err = migratePostgres.WithInstance(db, &migratePostgres.Config{SchemaName: schemaName})
	if err != nil {
		return err
	}

	migration, err = migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
