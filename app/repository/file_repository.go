package repository

import (
	"context"
	"crud_alumni/app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository struct {
    Collection *mongo.Collection
}
type FileRepo interface {
	Create(file *model.File) error
	GetAll() ([]model.File, error)
	GetByUserID(userID string) ([]model.File, error)
	GetByID(id string) (*model.File, error)
	DeleteByID(id string) error
}

func NewFileRepository(db *mongo.Database) *FileRepository {
    return &FileRepository{
        Collection: db.Collection("files"),
    }
}

func (r *FileRepository) Create(file *model.File) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    file.UploadedAt = time.Now()
    _, err := r.Collection.InsertOne(ctx, file)
    return err
}

func (r *FileRepository) FindByUser(userID primitive.ObjectID) ([]model.File, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := r.Collection.Find(ctx, bson.M{"user_id": userID})
    if err != nil {
        return nil, err
    }

    var files []model.File
    if err := cursor.All(ctx, &files); err != nil {
        return nil, err
    }
    return files, nil
}
func (r *FileRepository) GetAll() ([]model.File, error) {
	var files []model.File
	cursor, err := r.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &files)
	return files, nil
}

func (r *FileRepository) GetByUserID(userID string) ([]model.File, error) {
	oid, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.M{"user_id": oid}
	var files []model.File
	cursor, err := r.Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	cursor.All(context.TODO(), &files)
	return files, nil
}

func (r *FileRepository) GetByID(id string) (*model.File, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var file model.File
	err = r.Collection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) DeleteByID(id string) error {
	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := r.Collection.DeleteOne(context.TODO(), bson.M{"_id": oid})
	return err
}