// package route

// import (
// 	"crud_alumni/app/model"
// 	"crud_alumni/app/service"
// 	"crud_alumni/middleware"

// 	"github.com/gofiber/fiber/v2"
// )

// func SetupRoutes(app *fiber.App) {
//     api := app.Group("/api")

//     // Public
//     api.Post("/login", func(c *fiber.Ctx) error {
//         var req model.LoginRequest
//         if err := c.BodyParser(&req); err != nil {
//             return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
//         }
//         resp, err := service.Login(req)
//         if err != nil {
//             return c.Status(401).JSON(fiber.Map{"error": err.Error()})
//         }
//         return c.JSON(resp)
//     })

//     // Protected
//     protected := api.Group("", middleware.AuthRequired())

//     alumni := protected.Group("/alumni")
//     alumni.Get("/pag", service.GetAlumniPagination)
//     alumni.Get("/", service.GetAllAlumni)  // Admin + User
//     alumni.Get("/:id", service.GetAlumniByID) // Admin + User
//     alumni.Post("/", middleware.AdminOnly(), service.CreateAlumni)
//     alumni.Put("/:id", middleware.AdminOnly(), service.UpdateAlumni)
//     alumni.Delete("/:id", middleware.AdminOnly(), service.DeleteAlumni)

//     pekerjaan := protected.Group("/pekerjaan")
//     pekerjaan.Get("/", service.GetAllPekerjaan)
//     pekerjaan.Get("/trash", service.Trash)
//     pekerjaan.Get("/pag", service.GetPekerjaanPagination) // Admin + User
//     pekerjaan.Put("/:id/soft-delete", service.SoftDeletePekerjaan)
//     pekerjaan.Put("/:id/restore", service.RestorePekerjaan)
//     pekerjaan.Get("/:id", service.GetPekerjaanByID) // Admin + User
//     pekerjaan.Get("/tahun/:tahun", middleware.AdminOnly(), service.GetPekerjaanByTahun)
//     pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), service.GetPekerjaanByAlumniID)
//     pekerjaan.Post("/", middleware.AdminOnly(), service.CreatePekerjaan)
//     pekerjaan.Put("/:id", middleware.AdminOnly(), service.UpdatePekerjaan)
//     pekerjaan.Delete("/:id", middleware.AdminOnly(), service.DeletePekerjaan)
//     pekerjaan.Delete("/hard/:id", service.HardDeletePekerjaan)

// }

package route

import (
	"crud_alumni/app/repository"
	"crud_alumni/app/service"
	"crud_alumni/database"
	"crud_alumni/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/login", service.LoginHandler)

	// === ROUTES DENGAN AUTH ===
	protected := api.Group("", middleware.AuthRequired())

	// === ALUMNI ===
	alumni := protected.Group("/alumni")
	alumni.Get("/", service.GetAllAlumni)
	alumni.Get("/pag", service.GetAlumniPagination)
	alumni.Get("/:id", service.GetAlumniByID)
	alumni.Post("/", middleware.AdminOnly(), service.CreateAlumni)
	alumni.Put("/:id", middleware.AdminOnly(), service.UpdateAlumni)
	alumni.Delete("/:id", middleware.AdminOnly(), service.DeleteAlumni)

	// === PEKERJAAN ===
	pekerjaan := protected.Group("/pekerjaan")
    pekerjaan.Get("/", service.GetAllPekerjaan)
    pekerjaan.Get("/trash", service.GetTrashAll)
    pekerjaan.Get("/pag", service.GetPekerjaanByTahun) // Admin + User
    pekerjaan.Put("/:id/soft-delete", service.SoftDeletePekerjaan)
    pekerjaan.Put("/:id/restore", service.RestorePekerjaan)
    pekerjaan.Get("/:id", service.GetPekerjaanByID) // Admin + User
    pekerjaan.Get("/tahun/:tahun", middleware.AdminOnly(), service.GetPekerjaanByTahun)
    pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), service.GetPekerjaanByAlumniID)
    pekerjaan.Post("/", middleware.AdminOnly(), service.CreatePekerjaan)
    pekerjaan.Put("/:id", middleware.AdminOnly(), service.UpdatePekerjaan)
    pekerjaan.Delete("/:id", middleware.AdminOnly(), service.DeletePekerjaan)
    pekerjaan.Delete("/hard/:id", service.DeletePekerjaan)

	file := protected.Group("/file")
	fileRepo := repository.NewFileRepository(database.DB)
	fileService := service.NewFileService(fileRepo)
	file.Post("/foto", func(c *fiber.Ctx) error {
		return fileService.UploadFile(c, "foto")
	})
	file.Post("/sertifikat", func(c *fiber.Ctx) error {
        return fileService.UploadFile(c, "sertifikat")
    })
	file.Get("/", fileService.GetAllFiles)
	file.Get("/:id", fileService.GetFileByID)
	file.Delete("/:id", fileService.DeleteFile)

}
