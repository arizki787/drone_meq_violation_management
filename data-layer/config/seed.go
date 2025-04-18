package config

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"BE-DroneMEQ-Violation-Management/data-layer/model"

	"gorm.io/gorm"
)

// func Seeder(db *gorm.DB) error {
// 	if err := seeds.ListUserSeeder(db); err != nil {
// 		return err
// 	}

// 	if err := seeds.ListViolationSeeder(db); err != nil {
// 		return err
// 	}

// 	if err := seeds.ListViolationActionSeeder(db); err != nil {
// 		return err
// 	}

// 	if err := seeds.ListViolationMediaSeeder(db); err != nil {
// 		return err
// 	}

// 	if err := seeds.ListViolationNoteSeeder(db); err != nil {
// 		return err
// 	}

// 	return nil
// }

func RunSeeder(db *gorm.DB) {
	seed := false

	for _, arg := range os.Args[1:] {
		if arg == "--seed" {
			seed = true
			break
		}
	}

	if seed {
		if err := Seeder(db); err != nil {
			log.Fatalf("error during seeder: %v", err)
		}
		log.Println("Seeder completed successfully")
	} else {
		log.Println("No command provided. Use: --seed")
	}
}

func Seeder(db *gorm.DB) error {
	if err := SeedAll(db); err != nil {
		log.Printf("Seeder error: %v", err)
		return err
	}
	log.Println("Seeding completed successfully")
	return nil
}

func SeedAll(db *gorm.DB) error {
	if err := seedListViolationSeeder(db, "data-layer/seeder/violation.json"); err != nil {
		return err
	}

	if err := seedListViolationMediaSeeder(db, "data-layer/seeder/violation_media.json"); err != nil {
		return err
	}

	if err := seedListViolationNoteSeeder(db, "data-layer/seeder/violation_note.json"); err != nil {
		return err
	}

	if err := seedListViolationViolationMediaSeeder(db, "data-layer/seeder/violation_violation_media.json"); err != nil {
		return err
	}

	return nil
}

func seedListViolationSeeder(db *gorm.DB, filePath string) error {
	return seedData[model.Violation](db, filePath)
}
func seedListViolationMediaSeeder(db *gorm.DB, filePath string) error {
	return seedData[model.ViolationMedia](db, filePath)
}
func seedListViolationNoteSeeder(db *gorm.DB, filePath string) error {
	return seedData[model.ViolationNote](db, filePath)
}
func seedListViolationViolationMediaSeeder(db *gorm.DB, filePath string) error {
	return seedData[model.ViolationViolationMedia](db, filePath)
}

func seedData[T any](db *gorm.DB, filePath string) error {
	log.Printf("Attempting to open file: %s\n", filePath)
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var records []T
	if err := json.Unmarshal(jsonData, &records); err != nil {
		return err
	}

	for _, record := range records {
		var existing T
		err := db.Where("pkid = ?", getPkid(record)).First(&existing).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if getPkid(existing) == 0 {
			if err := db.Create(&record).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func getPkid(record interface{}) int64 {
	switch v := record.(type) {
	case model.Violation:
		return v.Pkid
	case model.ViolationMedia:
		return v.Pkid
	case model.ViolationNote:
		return v.Pkid
	default:
		return 0
	}
}
