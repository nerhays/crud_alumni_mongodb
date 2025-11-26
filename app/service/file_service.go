package service

import (
	"crud_alumni/app/model"
	"crud_alumni/app/repository"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileService struct {
	Repo repository.FileRepo
}

func NewFileService(repo repository.FileRepo) *FileService {
	return &FileService{Repo: repo}
}


// UploadFile godoc
// @Summary Upload file (foto/sertifikat)
// @Description Upload file foto (jpeg/png) atau sertifikat (pdf)
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Param category path string true "Jenis file (foto/sertifikat)"
// @Success 201 {object} model.File
// @Failure 400 {object} map[string]interface{}
// @Router /file/{category} [post]
func (s *FileService) UploadFile(c *fiber.Ctx, category string) error {
	// === Ambil user & role dari token JWT ===
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	// === Ambil query param target_id (hanya digunakan oleh admin) ===
	targetUserID := c.Query("target_id")

	// === Validasi role: user tidak boleh upload untuk orang lain ===
	if role == "user" && targetUserID != "" && targetUserID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized: user tidak boleh upload file untuk orang lain"})
	}

	// === Tentukan pemilik file ===
	var ownerID primitive.ObjectID
	if targetUserID != "" && role == "admin" {
		// Admin bisa upload untuk orang lain
		ownerID, _ = primitive.ObjectIDFromHex(targetUserID)
	} else {
		// Jika user biasa, hanya untuk dirinya sendiri
		ownerID, _ = primitive.ObjectIDFromHex(userID)
	}

	// === Ambil file dari form-data ===
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "No file uploaded"})
	}

	// === Validasi tipe file & ukuran ===
	contentType := fileHeader.Header.Get("Content-Type")
	var uploadPath string

	if category == "foto" {
		if contentType != "image/jpeg" && contentType != "image/png" {
			return c.Status(400).JSON(fiber.Map{"error": "Only jpeg/jpg/png allowed"})
		}
		if fileHeader.Size > 1*1024*1024 {
			return c.Status(400).JSON(fiber.Map{"error": "Max 1MB allowed"})
		}
		uploadPath = "./uploads/foto"
	} else if category == "sertifikat" {
		if contentType != "application/pdf" {
			return c.Status(400).JSON(fiber.Map{"error": "Only pdf allowed"})
		}
		if fileHeader.Size > 2*1024*1024 {
			return c.Status(400).JSON(fiber.Map{"error": "Max 2MB allowed"})
		}
		uploadPath = "./uploads/sertifikat"
	}

	// === Simpan file ke folder ===
	ext := filepath.Ext(fileHeader.Filename)
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(uploadPath, newFileName)

	os.MkdirAll(uploadPath, os.ModePerm)
	if err := c.SaveFile(fileHeader, filePath); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// === Simpan metadata ke database ===
	file := &model.File{
		UserID:       ownerID,
		FileName:     newFileName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
		FileType:     contentType,
		Category:     category,
	}

	if err := s.Repo.Create(file); err != nil {
		os.Remove(filePath)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save metadata"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"data":    file,
	})
}

// GetAllFiles godoc
// @Summary Dapatkan semua file
// @Description Admin dapat melihat semua file, user hanya file miliknya sendiri
// @Tags File
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} model.File
// @Failure 500 {object} map[string]interface{}
// @Router /file [get]
// === GET SEMUA FILE (admin bisa semua, user hanya miliknya sendiri)
func (s *FileService) GetAllFiles(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)

	var files []model.File
	var err error

	if role == "admin" {
		files, err = s.Repo.GetAll() // ambil semua file
	} else {
		files, err = s.Repo.GetByUserID(userID) // ambil hanya miliknya sendiri
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data file"})
	}

	return c.JSON(fiber.Map{
		"count": len(files),
		"data":  files,
	})
}

// GetFileByID godoc
// @Summary Dapatkan file berdasarkan ID
// @Description Admin dapat melihat semua file, user hanya file miliknya sendiri
// @Tags File
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 200 {object} model.File
// @Failure 404 {object} map[string]interface{}
// @Router /file/{id} [get]
// === GET FILE BY ID (admin bisa semua, user hanya miliknya sendiri)
func (s *FileService) GetFileByID(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)
	fileID := c.Params("id")

	file, err := s.Repo.GetByID(fileID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "File tidak ditemukan"})
	}

	// Jika user biasa, pastikan file miliknya
	if role != "admin" && file.UserID.Hex() != userID {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}

	return c.JSON(file)
}

// DeleteFile godoc
// @Summary Hapus file
// @Description Admin dapat menghapus semua file, user hanya miliknya sendiri
// @Tags File
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /file/{id} [delete]
// === DELETE FILE (admin bisa semua, user hanya miliknya sendiri)
func (s *FileService) DeleteFile(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(string)
	fileID := c.Params("id")

	file, err := s.Repo.GetByID(fileID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "File tidak ditemukan"})
	}

	// User hanya boleh hapus file miliknya
	if role != "admin" && file.UserID.Hex() != userID {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Hapus file fisik
	os.Remove(file.FilePath)

	// Hapus metadata dari database
	if err := s.Repo.DeleteByID(fileID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus data file"})
	}

	return c.JSON(fiber.Map{"message": "File berhasil dihapus"})
}
