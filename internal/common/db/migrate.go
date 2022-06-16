package db

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const DB_MIGRATION_SOURCES = "file://migrations"

func Migrate(db *mongo.Client) {
	dbname := databaseFromConnectionString(GetConnectionString())
	driver, err := mongodb.WithInstance(db, &mongodb.Config{
		MigrationsCollection: MIGRATIONS_COLLECTION,
		DatabaseName:         dbname,
	})
	if err != nil {
		panic(err)
	}
	m, errSetup := migrate.NewWithDatabaseInstance(
		DB_MIGRATION_SOURCES, dbname, driver,
	)
	if errSetup != nil {
		panic(errSetup)
	}
	errMigrate := m.Up()
	if errMigrate == nil {
		log.Println("db migrations complete!")
		return
	}
	if errMigrate.Error() == "no change" {
		log.Println("db is up-to-date. no new migrations!")
		return
	}
	panic(errMigrate)
}
