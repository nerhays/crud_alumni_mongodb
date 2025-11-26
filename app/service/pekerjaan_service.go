package service

import (
	"crud_alumni/app/model"
	"crud_alumni/app/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetAllPekerjaan godoc
// @Summary Dapatkan semua data pekerjaan
// @Description Mengambil seluruh data pekerjaan dari database
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Success 200 {array} model.Pekerjaan
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekerjaan [get]
func GetAllPekerjaan(c *fiber.Ctx) error {
	data, err := repository.GetAllPekerjaan()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(data)
}

// GetPekerjaanByID godoc
// @Summary Dapatkan pekerjaan berdasarkan ID
// @Description Mengambil data pekerjaan berdasarkan ID (ObjectID MongoDB)
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "ID pekerjaan"
// @Success 200 {object} model.Pekerjaan
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekerjaan/{id} [get]
func GetPekerjaanByID(c *fiber.Ctx) error {
	id := c.Params("id")

	data, err := repository.GetPekerjaanByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Data tidak ditemukan",
		})
	}

	return c.JSON(data)
}

// GetPekerjaanByAlumniID godoc
// @Summary Dapatkan pekerjaan berdasarkan ID alumni
// @Description Mengambil semua data pekerjaan berdasarkan alumni_id
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param alumni_id path int true "ID Alumni"
// @Success 200 {array} model.Pekerjaan
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekerjaan/alumni/{alumni_id} [get]
func GetPekerjaanByAlumniID(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "alumni_id tidak valid",
		})
	}

	data, err := repository.GetPekerjaanByAlumniID(alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(data)
}

// CreatePekerjaan godoc
// @Summary Tambahkan data pekerjaan
// @Description Menambahkan data pekerjaan baru (admin saja)
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param pekerjaan body model.Pekerjaan true "Data pekerjaan baru"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekerjaan [post]
func CreatePekerjaan(c *fiber.Ctx) error {
	var pekerjaan model.Pekerjaan
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gagal parse body",
		})
	}

	id, err := repository.CreatePekerjaan(pekerjaan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Data pekerjaan berhasil ditambahkan",
		"id":      id.Hex(),
	})
}

// UpdatePekerjaan godoc
// @Summary Perbarui data pekerjaan
// @Description Memperbarui data pekerjaan berdasarkan ID
// @Tags Pekerjaan
// @Accept json
// @Produce json
// @Param id path string true "ID pekerjaan"
// @Param pekerjaan body model.Pekerjaan true "Data pekerjaan yang diperbarui"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekerjaan/{id} [put]
func UpdatePekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")

	var pekerjaan model.Pekerjaan
	if err := c.BodyParser(&pekerjaan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gagal parse body",
		})
	}

	err := repository.UpdatePekerjaan(id, pekerjaan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data pekerjaan berhasil diperbarui",
	})
}

// DeletePekerjaan godoc
// @Summary Hapus pekerjaan secara permanen
// @Description Menghapus data pekerjaan secara hard delete
// @Tags Pekerjaan
// @Param id path string true "ID pekerjaan"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekerjaan/{id} [delete]
func DeletePekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")

	err := repository.DeletePekerjaan(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data pekerjaan berhasil dihapus permanen",
	})
}

// SoftDeletePekerjaan godoc
// @Summary Soft delete pekerjaan
// @Description Menghapus data pekerjaan tanpa benar-benar menghapus dari database
// @Tags Pekerjaan
// @Param id path string true "ID pekerjaan"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekerjaan/{id}/soft-delete [put]
func SoftDeletePekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")

	err := repository.SoftDeletePekerjaan(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data pekerjaan berhasil dihapus (soft delete)",
	})
}

// RestorePekerjaan godoc
// @Summary Pulihkan data pekerjaan yang dihapus
// @Description Mengembalikan data pekerjaan dari status soft delete
// @Tags Pekerjaan
// @Param id path string true "ID pekerjaan"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekerjaan/{id}/restore [put]
func RestorePekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")

	err := repository.RestorePekerjaan(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data pekerjaan berhasil dipulihkan",
	})
}

// GetTrashAll godoc
// @Summary Lihat semua data pekerjaan yang dihapus (soft delete)
// @Description Menampilkan daftar data pekerjaan yang masih tersimpan di trash
// @Tags Pekerjaan
// @Success 200 {array} model.Pekerjaan
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekerjaan/trash [get]
func GetTrashAll(c *fiber.Ctx) error {
	data, err := repository.TrashAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(data)
}

// GetPekerjaanByTahun godoc
// @Summary Statistik pekerjaan berdasarkan tahun
// @Description Mengambil jumlah pekerjaan berdasarkan tahun
// @Tags Pekerjaan
// @Param tahun path int true "Tahun pekerjaan"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekerjaan/tahun/{tahun} [get]
func GetPekerjaanByTahun(c *fiber.Ctx) error {
	tahun, err := strconv.Atoi(c.Params("tahun"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format tahun tidak valid",
		})
	}

	result, err := repository.GetPekerjaanByTahun(tahun)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
