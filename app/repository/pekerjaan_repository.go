package repository

import (
	"context"
	"crud_alumni/app/model"
	"crud_alumni/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllPekerjaan – ambil semua pekerjaan
func GetAllPekerjaan() ([]model.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.PekerjaanCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"_id": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []model.Pekerjaan
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetPekerjaanByID – ambil 1 dokumen berdasarkan ObjectID Mongo atau id lama (integer)
func GetPekerjaanByID(idStr string) (*model.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var pekerjaan model.Pekerjaan

	// coba konversi ke ObjectID dulu
	objID, err := primitive.ObjectIDFromHex(idStr)
	filter := bson.M{"_id": objID}
	if err != nil {
		// fallback ke pencarian berdasarkan id lama (Postgres ID)
		filter = bson.M{"id": idStr}
	}

	err = database.PekerjaanCollection.FindOne(ctx, filter).Decode(&pekerjaan)
	if err != nil {
		return nil, err
	}

	return &pekerjaan, nil
}

// GetPekerjaanByAlumniID – ambil semua pekerjaan dengan alumni_id tertentu
func GetPekerjaanByAlumniID(alumniID int) ([]model.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.PekerjaanCollection.Find(ctx, bson.M{"alumni_id": alumniID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []model.Pekerjaan
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CreatePekerjaan – tambah data baru
func CreatePekerjaan(p model.Pekerjaan) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p.IsDellete = "no"
	p.TanggalMulaiKerja = time.Now().Format("2006-01-02")

	result, err := database.PekerjaanCollection.InsertOne(ctx, p)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// UpdatePekerjaan – update data pekerjaan
func UpdatePekerjaan(idStr string, p model.Pekerjaan) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"alumni_id":            p.AlumniID,
			"nama_perusahaan":      p.NamaPerusahaan,
			"posisi_jabatan":       p.PosisiJabatan,
			"bidang_industri":      p.BidangIndustri,
			"lokasi_kerja":         p.LokasiKerja,
			"gaji_range":           p.GajiRange,
			"tanggal_mulai_kerja":  p.TanggalMulaiKerja,
			"tanggal_selesai_kerja": p.TanggalSelesaiKerja,
			"status_pekerjaan":     p.StatusPekerjaan,
			"deskripsi_pekerjaan":  p.Deskripsi,
			"updated_at":           time.Now(),
		},
	}

	_, err = database.PekerjaanCollection.UpdateByID(ctx, objID, update)
	return err
}

// DeletePekerjaan – hard delete
func DeletePekerjaan(idStr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}

	_, err = database.PekerjaanCollection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// Soft delete
func SoftDeletePekerjaan(idStr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}

	_, err = database.PekerjaanCollection.UpdateByID(ctx, objID, bson.M{
		"$set": bson.M{"isdellete": "yes", "updated_at": time.Now()},
	})
	return err
}

// Restore
func RestorePekerjaan(idStr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}

	_, err = database.PekerjaanCollection.UpdateByID(ctx, objID, bson.M{
		"$set": bson.M{"isdellete": "no", "updated_at": time.Now()},
	})
	return err
}

// GetPekerjaanByTahun – hitung pekerjaan berdasarkan tahun mulai kerja
func GetPekerjaanByTahun(tahun int) (model.JumlahPekerjaanPerTahun, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	startDate := time.Date(tahun, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(1, 0, 0)

	filter := bson.M{
		"tanggal_mulai_kerja": bson.M{
			"$gte": startDate.Format("2006-01-02"),
			"$lt":  endDate.Format("2006-01-02"),
		},
	}

	count, err := database.PekerjaanCollection.CountDocuments(ctx, filter)
	if err != nil {
		return model.JumlahPekerjaanPerTahun{Tahun: tahun, Jumlah: 0}, err
	}

	return model.JumlahPekerjaanPerTahun{
		Tahun:  tahun,
		Jumlah: int(count),
	}, nil
}

// TrashAll – ambil semua pekerjaan yang sudah soft delete
func TrashAll() ([]model.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.PekerjaanCollection.Find(ctx, bson.M{"isdellete": "yes"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []model.Pekerjaan
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}
