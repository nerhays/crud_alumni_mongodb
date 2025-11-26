package repository

import (
	"context"
	"crud_alumni/app/model"
	"crud_alumni/database"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ambil semua alumni
func GetAllAlumni() ([]model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("📡 Coba ambil semua alumni...")
	cursor, err := database.AlumniCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("❌ Error MongoDB Find:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.Alumni
	if err = cursor.All(ctx, &list); err != nil {
		fmt.Println("❌ Error Decode:", err)
		return nil, err
	}
	fmt.Println("✅ Jumlah alumni:", len(list))
	return list, nil
}


// Tambah alumni
func CreateAlumni(a model.Alumni) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.ID = primitive.NewObjectID()
	a.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
    a.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")


	_, err := database.AlumniCollection.InsertOne(ctx, a)
	return a.ID, err
}

// Update alumni
func UpdateAlumni(id string, a model.Alumni) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"nama":        a.Nama,
			"jurusan":     a.Jurusan,
			"angkatan":    a.Angkatan,
			"tahun_lulus": a.TahunLulus,
			"email":       a.Email,
			"no_telepon":  a.NoTelepon,
			"alamat":      a.Alamat,
			"updated_at":  time.Now(),
		},
	}
	_, err = database.AlumniCollection.UpdateByID(ctx, objID, update)
	return err
}

// Hapus alumni
func DeleteAlumni(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = database.AlumniCollection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// Get by ID
func GetAlumniByID(id string) (model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var a model.Alumni
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return a, err
	}
	err = database.AlumniCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&a)
	return a, err
}

// Pagination + Sorting + Searching
func GetAlumniWithPagination(search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"nama": bson.M{"$regex": search, "$options": "i"}},
				{"nim": bson.M{"$regex": search, "$options": "i"}},
				{"jurusan": bson.M{"$regex": search, "$options": "i"}},
				{"email": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

	sortOrder := 1
	if order == "desc" {
		sortOrder = -1
	}

	opts := options.Find().
		SetSort(bson.D{{Key: sortBy, Value: sortOrder}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := database.AlumniCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.Alumni
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}

	return list, nil
}

// Count total data
func CountAlumni(search string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"nama": bson.M{"$regex": search, "$options": "i"}},
				{"nim": bson.M{"$regex": search, "$options": "i"}},
				{"jurusan": bson.M{"$regex": search, "$options": "i"}},
				{"email": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

	count, err := database.AlumniCollection.CountDocuments(ctx, filter)
	return int(count), err
}

