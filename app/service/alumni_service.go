package service

import (
	"crud_alumni/app/model"
	"crud_alumni/app/repository"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetAllAlumni godoc
// @Summary Dapatkan semua data alumni
// @Description Mengambil seluruh data alumni dari database
// @Tags Alumni
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /alumni [get]
func GetAllAlumni(c *fiber.Ctx) error {
	data, err := repository.GetAllAlumni()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data"})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

// CreateAlumni godoc
// @Summary Tambah alumni baru
// @Description Menambahkan data alumni baru (hanya admin)
// @Tags Alumni
// @Accept json
// @Produce json
// @Param alumni body model.Alumni true "Data Alumni"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /alumni [post]
func CreateAlumni(c *fiber.Ctx) error {
	var a model.Alumni
	if err := c.BodyParser(&a); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
	}
	id, err := repository.CreateAlumni(a)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal tambah"})
	}
	a.ID = id
	return c.Status(201).JSON(fiber.Map{"success": true, "data": a})
}

// UpdateAlumni godoc
// @Summary Perbarui data alumni
// @Description Admin dapat memperbarui data alumni berdasarkan ID
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "ID Alumni"
// @Param alumni body model.Alumni true "Data Alumni"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /alumni/{id} [put]
func UpdateAlumni(c *fiber.Ctx) error {
	id := c.Params("id")
	var a model.Alumni
	if err := c.BodyParser(&a); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Body tidak valid"})
	}
	if err := repository.UpdateAlumni(id, a); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update"})
	}
	return c.JSON(fiber.Map{"success": true})
}

// DeleteAlumni godoc
// @Summary Hapus data alumni
// @Description Admin dapat menghapus alumni berdasarkan ID
// @Tags Alumni
// @Param id path string true "ID Alumni"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /alumni/{id} [delete]
func DeleteAlumni(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := repository.DeleteAlumni(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal hapus"})
	}
	return c.JSON(fiber.Map{"success": true})
}

// GetAlumniByID godoc
// @Summary Dapatkan detail alumni berdasarkan ID
// @Description Mengambil data detail 1 alumni berdasarkan ID
// @Tags Alumni
// @Param id path string true "ID Alumni"
// @Success 200 {object} model.Alumni
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /alumni/{id} [get]
func GetAlumniByID(c *fiber.Ctx) error {
	id := c.Params("id")
	a, err := repository.GetAlumniByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Alumni tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"success": true, "data": a})
}

// GetAlumniPagination godoc
// @Summary Dapatkan daftar alumni dengan pagination dan pencarian
// @Description Menampilkan daftar alumni berdasarkan halaman, urutan, dan kata kunci pencarian
// @Tags Alumni
// @Param page query int false "Nomor halaman (default 1)"
// @Param limit query int false "Jumlah data per halaman (default 10)"
// @Param sortBy query string false "Kolom pengurutan (nama/nim/angkatan/tahun_lulus/email)"
// @Param order query string false "Arah pengurutan (asc/desc)"
// @Param search query string false "Kata kunci pencarian"
// @Success 200 {object} model.AlumniResponse
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /alumni/pag [get]
func GetAlumniPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "nama")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	offset := (page - 1) * limit
	whitelist := map[string]bool{"nama": true, "nim": true, "angkatan": true, "tahun_lulus": true, "email": true}
	if !whitelist[sortBy] {
		sortBy = "nama"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	alumni, err := repository.GetAlumniWithPagination(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	total, _ := repository.CountAlumni(search)

	return c.JSON(model.AlumniResponse{
		Data: alumni,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	})
}
