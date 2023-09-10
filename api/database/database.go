package database

import (
	"fmt"
	"log"
	"os"

	"github.com/chadzink/skills-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DataAccessLayer struct {
	Db *gorm.DB
}

var DAL DataAccessLayer

func ConnectDb() error {
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		return err
		// os.Exit(2)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	migrateDb(db)

	DAL = DataAccessLayer{
		Db: db,
	}

	return nil
}

func migrateDb(db *gorm.DB) {
	log.Println("running migrations")
	db.AutoMigrate(&models.Skill{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Person{})
	db.AutoMigrate(&models.PersonSkill{})
}

func ConnectTestDb() error {
	// Open an SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return err
	}

	migrateDb(db)

	DAL = DataAccessLayer{
		Db: db,
	}

	return nil
}
