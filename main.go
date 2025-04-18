package main

import (
	"log"
	"os"

	"BE-DroneMEQ-Violation-Management/data-layer/config"
	"BE-DroneMEQ-Violation-Management/data-layer/repository"
	"BE-DroneMEQ-Violation-Management/domain-layer/client"
	"BE-DroneMEQ-Violation-Management/domain-layer/middleware"
	"BE-DroneMEQ-Violation-Management/domain-layer/service"
	"BE-DroneMEQ-Violation-Management/domain-layer/service/external"
	"BE-DroneMEQ-Violation-Management/presentation-layer/controller"
	"BE-DroneMEQ-Violation-Management/presentation-layer/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.SetUpDatabaseConnection()
	defer config.CloseDatabaseConnection(db)

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--migrate":
			if err := config.Migrate(db); err != nil {
				log.Fatalf("error migrating: %v", err)
			}
		case "--seed":
			if err := config.Seeder(db); err != nil {
				log.Fatalf("error seeding: %v", err)
			}
		case "--fresh":
			if err := config.Fresh(db); err != nil {
				log.Fatalf("error fresh: %v", err)
			}
			return
		}
	}

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	userManagementURL := os.Getenv("USER_MANAGEMENT_URL")
	if userManagementURL == ""{
		log.Fatalf("USER_MANAGEMENT_URL is not set")
	}
	userManagementClient := client.NewUserManagementClient(userManagementURL)
	userManagementService := external.NewUserManagementService(userManagementClient)

	violationRepo := repository.NewViolationRepository(db)
	violationActionRepo := repository.NewViolationActionRepository(db)
	violationNoteRepo := repository.NewViolationNoteRepository(db)
	violationMediaRepo := repository.NewViolationMediaRepository(db)

	violationService := service.NewViolationService(violationRepo)
	violationActionService := service.NewViolationActionService(violationActionRepo)
	violationNoteService := service.NewViolationNoteService(violationNoteRepo)
	violationMediaService := service.NewViolationMediaService(violationMediaRepo, violationRepo)

	violationController := controller.NewViolationController(violationService)
	violationActionController := controller.NewViolationActionController(violationActionService)
	violationNoteController := controller.NewViolationNoteController(violationNoteService)
	violationMediaController := controller.NewViolationMediaController(violationMediaService)

	routes.Violation(server, violationController, violationMediaController, violationNoteController, violationActionController, userManagementService)

	if err := server.Run(":3030"); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
