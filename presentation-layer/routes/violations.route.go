package routes

import (
	"github.com/gin-gonic/gin"

	"BE-DroneMEQ-Violation-Management/domain-layer/middleware"
	"BE-DroneMEQ-Violation-Management/domain-layer/service/external"
	"BE-DroneMEQ-Violation-Management/presentation-layer/controller"
)

func Violation(route *gin.Engine, violationController controller.ViolationController, violationMediaController controller.ViolationMediaController, violationNoteController controller.ViolationNoteController, violationActionController controller.ViolationActionController, userManagementService external.UserManagementService) {
	routes := route.Group("/api/violations")
	routes.Use(middleware.AuthMiddleware(userManagementService))
	{
		// Violation
		routes.POST("", violationController.CreateVs)
		routes.GET("", violationController.GetAllViolation)
		routes.GET("/:Pkid", violationController.ViolationById)
		routes.GET("/filter", violationController.ViolationByFilter)
		routes.PUT("/:Pkid", violationController.UpdateVs)
		routes.PUT("/:Pkid/soft-delete", violationController.SoftDeleteVs)
		routes.PUT("/:Pkid/restore", violationController.RestoreVs)
		routes.DELETE("/:Pkid", violationController.DeleteVs)

		// Violation Media
		routes.POST("/:Pkid/media", violationMediaController.CreateVsMedia)
		// Violation Note
		routes.POST("/:Pkid/notes", violationNoteController.CreateVsNote)                // pastikan bisa by Office Pkid
		routes.GET("/:Pkid/notes", violationNoteController.GetVsNote)                    //filter note type
		routes.GET("/:Pkid/notes/filternote", violationNoteController.GetVsNoteByFilter) //filter note type
		routes.PUT("/:Pkid/notes/:notePkid", violationNoteController.UpdateVsNote)
		routes.PUT("/:Pkid/notes/:notePkid/soft-delete", violationNoteController.SoftDeleteVsNote)
		routes.PUT("/:Pkid/notes/:notePkid/restore", violationNoteController.RestoreVsNote)
		routes.DELETE("/:Pkid/notes/:notePkid", violationNoteController.DeleteVsNote)
		// Violation Action
		routes.POST("/:Pkid/actions", violationActionController.CreateVsAction)
	}
}
