package database

import (
	"fmt"
	"log"
	"os"

	"skills-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DataAccessLayer struct {
	Db *gorm.DB
}

var DAL DataAccessLayer

func ConnectDb() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
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

	MigrateDb(db)

	DAL = DataAccessLayer{
		Db: db,
	}

	return nil
}

func MigrateDb(db *gorm.DB) {
	log.Println("running migrations")
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Skill{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Person{})
	db.AutoMigrate(&models.PersonSkill{})
	db.AutoMigrate(&models.Expertise{})
	// Add default expertise levels
	CreateDefaultExpertise(db)
}

func CreateDefaultExpertise(db *gorm.DB) {
	defaultExpertises := []models.Expertise{
		{
			Model:       gorm.Model{ID: 1},
			Name:        "Beginner",
			Description: "A beginner is a person who is starting to learn or do something.",
			Order:       1,
		},
		{
			Model:       gorm.Model{ID: 2},
			Name:        "Intermediate",
			Description: "An intermediate is a person who has a level of knowledge or skill between a beginner and an expert.",
			Order:       2,
		},
		{
			Model:       gorm.Model{ID: 3},
			Name:        "Advanced",
			Description: "An advanced is a person who is very skilled or highly trained in a particular field.",
			Order:       3,
		},
		{
			Model:       gorm.Model{ID: 4},
			Name:        "Expert",
			Description: "An expert is a person who is very knowledgeable about or skilful in a particular area.",
			Order:       4,
		},
		{
			Model:       gorm.Model{ID: 5},
			Name:        "N/A",
			Description: "N/A is used when the level of expertise is not applicable.",
			Order:       5,
		},
	}

	// Add default expertise levels
	for _, expertise := range defaultExpertises {

		// Check if the expertise already exists
		var existingExpertise models.Expertise
		err := db.Where("name = ?", expertise.Name).First(&existingExpertise).Error
		if err != nil {
			// Log the error
			log.Println(err)
		}

		if existingExpertise.ID > 0 {
			continue
		} else {
			err = db.Create(&expertise).Error

			if err != nil {
				// Log the error
				log.Println(err)
			}
		}
	}
}
