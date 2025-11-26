package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID       primitive.ObjectID `bson:"user_id,omitempty" json:"user_id"`
    FileName     string             `bson:"file_name" json:"file_name"`
    OriginalName string             `bson:"original_name" json:"original_name"`
    FilePath     string             `bson:"file_path" json:"file_path"`
    FileSize     int64              `bson:"file_size" json:"file_size"`
    FileType     string             `bson:"file_type" json:"file_type"`
    Category     string             `bson:"category" json:"category"` // "foto" atau "sertifikat"
    UploadedAt   time.Time          `bson:"uploaded_at" json:"uploaded_at"`
}
