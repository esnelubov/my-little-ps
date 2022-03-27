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
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"log"
	"my-little-ps/common/config"
)

type DB struct {
	gormDB            *gorm.DB
	ErrRecordNotFound error
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
		gormDB:            db,
		ErrRecordNotFound: gorm.ErrRecordNotFound,
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

func (db *DB) createWhereQuery(criteria map[string]interface{}, tryLocking bool) (query *gorm.DB) {
	var (
		ok     bool
		offset interface{}
		limit  interface{}
	)

	if tryLocking {
		query = db.gormDB.Clauses(clause.Locking{
			Strength: "UPDATE",
		})
	} else {
		query = db.gormDB
	}

	offset, ok = criteria["offset"]
	if ok {
		query = query.Offset(offset.(int))
		delete(criteria, "offset")
	}

	limit, ok = criteria["limit"]
	if ok {
		query = query.Limit(limit.(int))
		delete(criteria, "limit")
	}

	for c, v := range criteria {
		query = query.Where(c, v)
	}

	return query
}

func (db *DB) Last(record interface{}, criteria map[string]interface{}, tryLocking bool) error {
	query := db.createWhereQuery(criteria, tryLocking)

	return query.Last(record).Error
}

func (db *DB) Has(record interface{}, criteria map[string]interface{}) (bool, error) {
	query := db.createWhereQuery(criteria, false)

	err := query.Take(record).Error

	if err == nil {
		return true, nil
	}

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	return false, err
}

func (db *DB) Find(records interface{}, criteria map[string]interface{}, tryLocking bool) error {
	query := db.createWhereQuery(criteria, tryLocking)

	return query.Find(records).Error
}

func (db *DB) Raw(result interface{}, sql string, params ...interface{}) error {
	return db.gormDB.Raw(sql, params...).Scan(result).Error
}

func (db *DB) Create(value interface{}) error {
	return db.gormDB.Create(value).Error
}

func (db *DB) Save(value interface{}) error {
	return db.gormDB.Save(value).Error
}

func (db *DB) Transaction(fc func(tx *DB) error) error {
	return db.gormDB.Transaction(
		func(tx *gorm.DB) error {
			return fc(&DB{gormDB: tx,
				ErrRecordNotFound: db.ErrRecordNotFound})
		})
}

func (db *DB) TableName(table string) string {
	return db.gormDB.NamingStrategy.TableName(table)
}
